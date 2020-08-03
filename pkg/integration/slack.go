package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/convox/console/pkg/crypt"
	"github.com/convox/console/pkg/settings"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	oauth2slack "golang.org/x/oauth2/slack"
)

type Slack struct {
	client_id     string
	client_secret string
	host          string
	oid           string
	token         string
}

func (i *Slack) Authorize() (string, error) {
	c, err := i.oauthClient()
	if err != nil {
		return "", errors.WithStack(err)
	}

	c.RedirectURL = fmt.Sprintf("https://%s/integrations/authorize/slack?organization=%s", settings.Host, i.oid)

	return c.AuthCodeURL(crypt.OneWay(i.oid), oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (i *Slack) Exchange(code string, reauthorize bool) (string, map[string]string, error) {
	payload := url.Values{
		"client_id":     {settings.SlackClientId},
		"client_secret": {settings.SlackClientSecret},
		"code":          {code},
		"redirect_uri":  {fmt.Sprintf("https://%s/integrations/authorize/slack?organization=%s", settings.Host, i.oid)},
		"state":         {crypt.OneWay(i.oid)},
	}

	res, err := http.PostForm("https://slack.com/api/oauth.access", payload)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	var auth struct {
		AccessToken string `json:"access_token"`
		OK          bool   `json:"ok"`
		Webhook     struct {
			Channel       string `json:"channel"`
			Configuration string `json:"configuration_url"`
			URL           string `json:"url"`
		} `json:"incoming_webhook"`
	}

	if err := json.Unmarshal(data, &auth); err != nil {
		return "", nil, errors.WithStack(err)
	}

	if !auth.OK {
		return "", nil, errors.WithStack(fmt.Errorf("invalid exchange"))
	}

	attrs := map[string]string{
		"channel":       auth.Webhook.Channel,
		"configuration": auth.Webhook.Configuration,
		"webhook":       auth.Webhook.URL,
	}

	return auth.AccessToken, attrs, nil
}

func (i *Slack) EventSend(rack string, attrs map[string]string, e NotificationEvent) error {
	msg, err := formatEvent(rack, e)
	if err != nil {
		return errors.WithStack(err)
	}

	if msg == "" {
		return nil
	}

	hook, ok := attrs["webhook"]
	if !ok {
		return errors.WithStack(fmt.Errorf("no webhook url"))
	}

	color := ""

	switch e.Status {
	case "success":
		color = "good"
	case "error":
		color = "danger"
	}

	data := fmt.Sprintf(`{
		"attachments": [
		  {
				"color": %q,
				"text": %q
			}
		]
	}`, color, msg)

	res, err := http.PostForm(hook, url.Values{"payload": {string(data)}})
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	return nil
}

func (i *Slack) Name() string {
	return "Slack"
}

func (i *Slack) Revoke() error {
	return nil
}

func (i *Slack) Slug() string {
	return "slack"
}

func (i *Slack) Status() (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	if _, err := c.AuthTest(); err != nil {
		return "disconnected", nil
	}

	return "connected", nil
}

func (i *Slack) Title(attrs map[string]string) (string, error) {
	if ch := attrs["channel"]; ch != "" {
		return ch, nil
	}

	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	res, err := c.AuthTest()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return res.Team, nil
}

func (i *Slack) client() (*slack.Client, error) {
	return slack.New(i.token), nil
}

func (i *Slack) oauthClient() (*oauth2.Config, error) {
	c := &oauth2.Config{
		ClientID:     i.client_id,
		ClientSecret: i.client_secret,
		Endpoint:     oauth2slack.Endpoint,
		Scopes:       []string{"incoming-webhook"},
	}

	// if i.host != "" {
	//   c.Endpoint.AuthURL = fmt.Sprintf("https://%s/login/oauth/authorize", i.host)
	//   c.Endpoint.TokenURL = fmt.Sprintf("https://%s/login/oauth/access_token", i.host)
	// }

	return c, nil
}

func formatEvent(rack string, event NotificationEvent) (string, error) {
	var message string

	data := event.Data

	switch event.Action {
	case "app:create":
		message = fmt.Sprintf("Created app *%s*", data["name"])
	case "app:delete":
		message = fmt.Sprintf("Deleted app *%s*", data["name"])
	case "build:create":
		if event.Status == "error" {
			message = fmt.Sprintf("Build `%s` failed for app *%s*", data["id"], data["app"])
		} else {
			message = ""
		}
	case "rack:capacity":
		message = ""
	case "rack:converge":
		message = ""
	case "release:create":
		message = fmt.Sprintf("Created release `%s` for app *%s*", data["id"], data["app"])
	case "release:finish":
		message = fmt.Sprintf("Finished rolling deploy of release `%s` for app *%s*", data["id"], data["app"])
	case "release:promote":
		if event.Status == "success" {
			message = fmt.Sprintf("Promoted release `%s` for app *%s*", data["id"], data["app"])
		} else if event.Status == "error" {
			message = fmt.Sprintf("Promoting release `%s` for app *%s* failed", data["id"], data["app"])
		}
	case "release:scale":
		message = fmt.Sprintf("Scaled release `%s` for app *%s*", data["id"], data["app"])
	case "resource:create", "service:create":
		message = fmt.Sprintf("Created %s resource *%s*", data["type"], data["name"])
	case "resource:delete", "service:delete":
		message = fmt.Sprintf("Deleted %s resource *%s*", data["type"], data["name"])
	case "system:update", "rack:update":
		var upgrades []string
		if data["version"] != "" {
			upgrades = append(upgrades, fmt.Sprintf("version *%s*", data["version"]))
		}
		if data["count"] != "" {
			upgrades = append(upgrades, fmt.Sprintf("count *%s*", data["count"]))
		}
		if data["type"] != "" {
			upgrades = append(upgrades, fmt.Sprintf("instance type *%s*", data["type"]))
		}
		if len(upgrades) > 0 {
			message = fmt.Sprintf("Updating rack to: %s", strings.Join(upgrades, ", "))
		}
	case "workflow:complete":
		switch event.Status {
		case "cancel":
			switch event.Data["kind"] {
			case "deployment":
				message = fmt.Sprintf("Deployment workflow cancelled")
			case "review":
				message = fmt.Sprintf("Review workflow cancelled")
			}
		case "error":
			switch event.Data["kind"] {
			case "deployment":
				message = fmt.Sprintf("Deployment workflow failed")
			case "review":
				message = fmt.Sprintf("Review workflow failed")
			}
		default:
			switch event.Data["kind"] {
			case "deployment":
				message = fmt.Sprintf("Deployment workflow complete")
			case "review":
				message = fmt.Sprintf("Review workflow complete")
			}
		}
	}

	if message != "" {
		message = fmt.Sprintf("[*%s*] %s", rack, message)
	}

	return message, nil
}
