package integration

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/convox/console/pkg/crypt"
	"github.com/convox/console/pkg/helpers"
	"github.com/convox/console/pkg/settings"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

type Github struct {
	client_id     string
	client_secret string
	host          string
	id            string
	oid           string
	provider      string
	token         string
}

func (i *Github) Authorize() (string, error) {
	c, err := i.oauthClient()
	if err != nil {
		return "", errors.WithStack(err)
	}

	if i.id != "" {
		c.RedirectURL = fmt.Sprintf("https://%s/integrations/reauthorize?organization=%s&integration=%s", settings.Host, i.oid, i.id)
	} else {
		c.RedirectURL = fmt.Sprintf("https://%s/integrations/authorize/%s?organization=%s", settings.Host, i.provider, i.oid)
	}

	return c.AuthCodeURL(crypt.OneWay(i.oid), oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (i *Github) Exchange(code string, reauthorize bool) (string, map[string]string, error) {
	c, err := i.oauthClient()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	t, err := c.Exchange(i.oauthContext(), code, oauth2.AccessTypeOffline)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	attrs := map[string]string{}

	return t.AccessToken, attrs, nil
}

func (i *Github) Name() string {
	switch i.provider {
	case "github-enterprise":
		return "Github Enterprise"
	default:
		return "Github"
	}
}

func (i *Github) RepositoryClone(id int64, ref string, w io.Writer) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	r, _, err := c.Repositories.GetByID(context.Background(), id)
	if err != nil {
		return "", errors.WithStack(err)
	}

	cu, err := url.Parse(r.GetCloneURL())
	if err != nil {
		return "", errors.WithStack(err)
	}

	cu.User = url.User(i.token)

	return gitClone(w, cu.String(), ref)
}

func (i *Github) RepositoryList() (map[int64]string, error) {
	// if rsh, ok := cache.Get("repositories.github", i.id).(map[int64]string); ok {
	//   return rsh, nil
	// }

	c, err := i.client()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rsh := map[int64]string{}

	page := 1

	for {
		opts := &github.RepositoryListOptions{
			Type:      "all",
			Sort:      "full_name",
			Direction: "asc",
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		}

		rs, meta, err := c.Repositories.List(context.Background(), "", opts)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, r := range rs {
			rsh[r.GetID()] = r.GetFullName()
		}

		if meta.NextPage == 0 {
			break
		}

		page = meta.NextPage
	}

	// cache.Set("repositories.github", i.id, rsh, 30*time.Second)

	return rsh, nil
}

func (i *Github) RepositoryName(id int64) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	r, _, err := c.Repositories.GetByID(context.Background(), id)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return r.GetFullName(), nil
}

func (i *Github) Revoke() error {
	return nil
}

func (i *Github) Slug() string {
	return i.provider
}

func (i *Github) Status() (string, error) {
	c, err := i.client()
	if err != nil {
		return "disconnected", nil
	}

	u, _, err := c.Users.Get(context.Background(), "")
	if err != nil {
		return "disconnected", nil
	}

	if u.Login == nil {
		return "disconnected", nil
	}

	return "connected", nil
}

func (i *Github) StatusUpdate(repo int64, ref, state, description, url string) error {
	c, err := i.client()
	if err != nil {
		return errors.WithStack(err)
	}

	r, _, err := c.Repositories.GetByID(context.Background(), repo)
	if err != nil {
		return errors.WithStack(err)
	}

	ctx := "convox"

	if settings.Development {
		ctx = "convox/dev"
	}

	status := &github.RepoStatus{
		Context:     github.String(ctx),
		Description: github.String(description),
		State:       github.String(state),
	}

	if url != "" {
		status.TargetURL = github.String(url)
	}

	if _, _, err := c.Repositories.CreateStatus(context.Background(), r.GetOwner().GetLogin(), r.GetName(), ref, status); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *Github) Title(attrs map[string]string) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	u, _, err := c.Users.Get(context.Background(), "")
	if err != nil {
		return "", nil
	}

	return u.GetLogin(), nil
}

func (i *Github) WebhookCreateMerge(repo int64, branch, url string) (int64, error) {
	c, err := i.client()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	r, res, err := c.Repositories.GetByID(context.Background(), repo)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if i.host != "" {
		if res.Header.Get("X-Github-Enterprise-Version") < "2.15" {
			params := map[string]interface{}{
				"config": map[string]string{
					"content_type": "json",
					"secret":       settings.GithubWebhookSecret,
					"url":          url,
				},
				"events": []string{"push"},
				"name":   "web",
			}

			hreq, err := c.NewRequest("POST", fmt.Sprintf("repos/%v/%v/hooks", r.Owner.GetLogin(), r.GetName()), params)
			if err != nil {
				return 0, errors.WithStack(err)
			}

			h := new(github.Hook)

			if _, err := c.Do(context.Background(), hreq, h); err != nil {
				return 0, errors.WithStack(err)
			}

			return h.GetID(), nil
		}
	}

	hook := &github.Hook{
		Config: map[string]interface{}{
			"content_type": "json",
			"secret":       settings.GithubWebhookSecret,
			"url":          url,
		},
		Events: []string{"push"},
	}

	h, _, err := c.Repositories.CreateHook(context.Background(), r.Owner.GetLogin(), r.GetName(), hook)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return h.GetID(), nil
}

func (i *Github) WebhookCreateReview(repo int64, url string) (int64, error) {
	c, err := i.client()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	r, res, err := c.Repositories.GetByID(context.Background(), repo)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if i.host != "" {
		if res.Header.Get("X-Github-Enterprise-Version") < "2.15" {
			params := map[string]interface{}{
				"config": map[string]string{
					"content_type": "json",
					"secret":       settings.GithubWebhookSecret,
					"url":          url,
				},
				"events": []string{"pull_request"},
				"name":   "web",
			}

			hreq, err := c.NewRequest("POST", fmt.Sprintf("repos/%v/%v/hooks", r.Owner.GetLogin(), r.GetName()), params)
			if err != nil {
				return 0, errors.WithStack(err)
			}

			h := new(github.Hook)

			if _, err := c.Do(context.Background(), hreq, h); err != nil {
				return 0, errors.WithStack(err)
			}

			return h.GetID(), nil
		}
	}

	hook := &github.Hook{
		Config: map[string]interface{}{
			"content_type": "json",
			"secret":       settings.GithubWebhookSecret,
			"url":          url,
		},
		Events: []string{"pull_request"},
	}

	h, _, err := c.Repositories.CreateHook(context.Background(), r.Owner.GetLogin(), r.GetName(), hook)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return h.GetID(), nil
}

func (i *Github) WebhookDelete(repo, hook int64) error {
	c, err := i.client()
	if err != nil {
		return errors.WithStack(err)
	}

	r, _, err := c.Repositories.GetByID(context.Background(), repo)
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := c.Repositories.DeleteHook(context.Background(), r.Owner.GetLogin(), r.GetName(), hook); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *Github) WebhookPayload(r *http.Request, validate bool) (*SourceWebhookPayload, error) {
	var payload []byte

	if validate {
		p, err := github.ValidatePayload(r, []byte(settings.GithubWebhookSecret))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		payload = p
	} else {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		payload = data
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch t := event.(type) {
	case *github.PingEvent:
		return &SourceWebhookPayload{Event: "ping"}, nil
	case *github.PullRequestEvent:
		if !t.GetRepo().GetPrivate() {
			c, err := i.client()
			if err != nil {
				return nil, errors.WithStack(err)
			}

			pl, _, err := c.Repositories.GetPermissionLevel(context.Background(), t.GetRepo().GetOwner().GetLogin(), t.GetRepo().GetName(), t.GetSender().GetLogin())
			if err != nil {
				return nil, errors.WithStack(err)
			}

			switch pl.GetPermission() {
			case "admin", "push", "write":
			default:
				fmt.Printf("disallowed permission: %s\n", pl.GetPermission())
				return nil, nil
			}
		}

		swh := &SourceWebhookPayload{
			Name:  fmt.Sprintf("%s-%d", t.GetRepo().GetName(), t.GetPullRequest().GetNumber()),
			Ref:   t.GetPullRequest().GetHead().GetSHA(),
			Title: fmt.Sprintf("[#%d] %s", t.GetNumber(), t.GetPullRequest().GetTitle()),
			URL:   t.GetPullRequest().GetHTMLURL(),
		}

		switch t.GetAction() {
		case "opened", "reopened":
			swh.Event = "review.open"
			return swh, nil
		case "synchronize":
			swh.Event = "review.update"
			return swh, nil
		case "closed":
			swh.Event = "review.close"
			return swh, nil
		default:
			return nil, errors.WithStack(fmt.Errorf("unknown action: %s", t.GetAction()))
		}
	case *github.PushEvent:
		if t.GetDeleted() {
			return nil, nil
		}

		repo := t.GetRepo()
		head := t.GetHeadCommit()
		rparts := strings.Split(t.GetRef(), "/")

		swh := &SourceWebhookPayload{
			Description: fmt.Sprintf("push %s:%s %s %s", repo.GetFullName(), rparts[len(rparts)-1], helpers.Truncate(head.GetID(), 10, false), helpers.Truncate(head.GetMessage(), 32, true)),
			Event:       "merge",
			Name:        t.GetRef(),
			Ref:         t.GetAfter(),
			Title:       t.GetHeadCommit().GetMessage(),
			URL:         t.GetHeadCommit().GetURL(),
		}
		return swh, nil
	case *github.StatusEvent:
		return nil, nil
	}

	return nil, errors.WithStack(fmt.Errorf("unknown event"))
}

func (i *Github) client() (*github.Client, error) {
	t := &oauth2.Token{
		AccessToken: i.token,
	}

	oc, err := i.oauthClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := github.NewClient(oc.Client(i.oauthContext(), t))

	if i.host != "" {
		u, err := url.Parse(fmt.Sprintf("https://%s/api/v3/", i.host))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		c.BaseURL = u
		c.UploadURL = u
	}

	return c, nil
}

func (i *Github) defaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func (i *Github) oauthClient() (*oauth2.Config, error) {
	c := &oauth2.Config{
		ClientID:     i.client_id,
		ClientSecret: i.client_secret,
		Endpoint:     oauth2github.Endpoint,
		Scopes:       []string{"repo"},
	}

	if i.host != "" {
		c.Endpoint.AuthURL = fmt.Sprintf("https://%s/login/oauth/authorize", i.host)
		c.Endpoint.TokenURL = fmt.Sprintf("https://%s/login/oauth/access_token", i.host)
	}

	return c, nil
}

func (i *Github) oauthContext() context.Context {
	ctx := context.Background()

	if i.host != "" && settings.GithubEnterpriseInsecure {
		t := i.defaultTransport()

		t.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}

		c := *http.DefaultClient

		c.Transport = t

		ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &c)
	}

	return ctx
}
