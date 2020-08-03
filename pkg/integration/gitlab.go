package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/convox/console/pkg/crypt"
	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	oauth2gitlab "golang.org/x/oauth2/gitlab"
)

type Gitlab struct {
	client_id     string
	client_secret string
	host          string
	id            string
	oid           string
	token         string
}

func (i *Gitlab) Authorize() (string, error) {
	c, err := i.oauthClient()
	if err != nil {
		return "", errors.WithStack(err)
	}

	if i.id != "" {
		c.RedirectURL = fmt.Sprintf("https://%s/integrations/reauthorize?organization=%s&integration=%s", settings.Host, i.oid, i.id)
	} else {
		c.RedirectURL = fmt.Sprintf("https://%s/integrations/authorize/gitlab?organization=%s", settings.Host, i.oid)
	}

	return c.AuthCodeURL(crypt.OneWay(i.oid), oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (i *Gitlab) Exchange(code string, reauthorize bool) (string, map[string]string, error) {
	c, err := i.oauthClient()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if reauthorize {
		c.RedirectURL = fmt.Sprintf("https://%s/integrations/reauthorize?organization=%s&integration=%s", settings.Host, i.oid, i.id)
	} else {
		c.RedirectURL = fmt.Sprintf("https://%s/integrations/authorize/gitlab?organization=%s", settings.Host, i.oid)
	}

	t, err := c.Exchange(oauth2.NoContext, code, oauth2.AccessTypeOffline)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	attrs := map[string]string{}

	return t.AccessToken, attrs, nil
}

func (i *Gitlab) Name() string {
	return "Gitlab"
}

func (i *Gitlab) RepositoryClone(id int64, ref string, w io.Writer) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	r, _, err := c.Projects.GetProject(int(id))
	if err != nil {
		return "", errors.WithStack(err)
	}

	url := fmt.Sprintf("https://oauth2:%s@gitlab.com/%s", i.token, r.PathWithNamespace)

	return gitClone(w, url, ref)
}

func (i *Gitlab) RepositoryList() (map[int64]string, error) {
	c, err := i.client()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rsh := map[int64]string{}

	page := 1

	for {
		opts := &gitlab.ListProjectsOptions{
			Membership: gitlab.Bool(true),
			OrderBy:    gitlab.String("name"),
			Sort:       gitlab.String("asc"),
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		}

		rs, meta, err := c.Projects.ListProjects(opts)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, r := range rs {
			rsh[int64(r.ID)] = r.PathWithNamespace
		}

		if meta.NextPage == 0 {
			break
		}

		page = meta.NextPage
	}

	return rsh, nil
}

func (i *Gitlab) RepositoryName(id int64) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	r, _, err := c.Projects.GetProject(int(id))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return r.PathWithNamespace, nil
}

func (i *Gitlab) Revoke() error {
	return nil
}

func (i *Gitlab) Slug() string {
	return "gitlab"
}

func (i *Gitlab) Status() (string, error) {
	c, err := i.client()
	if err != nil {
		return "disconnected", nil
	}

	u, _, err := c.Users.CurrentUser()
	if err != nil {
		return "disconnected", nil
	}

	if u.Username == "" {
		return "disconnected", nil
	}

	return "connected", nil
}

func (i *Gitlab) StatusUpdate(repo int64, ref, state, description, url string) error {
	return nil
}

func (i *Gitlab) Title(attrs map[string]string) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	u, _, err := c.Users.CurrentUser()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return u.Username, nil
}

func (i *Gitlab) WebhookCreateMerge(id int64, branch, url string) (int64, error) {
	c, err := i.client()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	opts := &gitlab.AddProjectHookOptions{
		URL:                 gitlab.String(url),
		PushEvents:          gitlab.Bool(true),
		IssuesEvents:        gitlab.Bool(false),
		MergeRequestsEvents: gitlab.Bool(false),
	}

	h, _, err := c.Projects.AddProjectHook(int(id), opts)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(h.ID), nil
}

func (i *Gitlab) WebhookCreateReview(id int64, url string) (int64, error) {
	c, err := i.client()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	opts := &gitlab.AddProjectHookOptions{
		URL:                 gitlab.String(url),
		PushEvents:          gitlab.Bool(false),
		IssuesEvents:        gitlab.Bool(false),
		MergeRequestsEvents: gitlab.Bool(true),
	}

	h, _, err := c.Projects.AddProjectHook(int(id), opts)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(h.ID), nil
}

func (i *Gitlab) WebhookDelete(id, hook int64) error {
	c, err := i.client()
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := c.Projects.DeleteProjectHook(int(id), int(hook)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *Gitlab) WebhookPayload(r *http.Request, validate bool) (*SourceWebhookPayload, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// fmt.Printf("string(data) = %+v\n", string(data))

	var event struct {
		After            string `json:"after"`
		ObjectKind       string `json:"object_kind"`
		ObjectAttributes struct {
			Action     string `json:"action"`
			IID        int64  `json:"iid"`
			LastCommit struct {
				ID string `json:"id"`
			} `json:"last_commit"`
			State string `json:"state"`
		} `json:"object_attributes"`
		Project struct {
			Name string `json:"name"`
		} `json:"project"`
		Ref string `json:"ref"`
	}

	// fmt.Printf("string(data) = %+v\n", string(data))

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, errors.WithStack(err)
	}

	// fmt.Printf("event = %+v\n", event)

	switch event.ObjectKind {
	case "merge_request":
		swh := &SourceWebhookPayload{
			Name: fmt.Sprintf("%s-%d", event.Project.Name, event.ObjectAttributes.IID),
			Ref:  event.ObjectAttributes.LastCommit.ID,
		}

		switch event.ObjectAttributes.Action {
		case "open", "reopen":
			swh.Event = "review.open"
			return swh, nil
		case "update":
			swh.Event = "review.update"
			return swh, nil
		case "close":
			swh.Event = "review.close"
			return swh, nil
		default:
			return nil, errors.WithStack(fmt.Errorf("unknown action: %s", event.ObjectAttributes.Action))
		}
	case "push":
		swh := &SourceWebhookPayload{
			Event: "merge",
			Name:  event.Ref,
			Ref:   event.After,
		}
		return swh, nil
	}

	return nil, errors.WithStack(fmt.Errorf("unknown event"))
}

func (i *Gitlab) client() (*gitlab.Client, error) {
	return gitlab.NewOAuthClient(nil, i.token), nil
}

func (i *Gitlab) oauthClient() (*oauth2.Config, error) {
	c := &oauth2.Config{
		ClientID:     i.client_id,
		ClientSecret: i.client_secret,
		Endpoint:     oauth2gitlab.Endpoint,
		Scopes:       []string{"api"},
	}

	// if i.host != "" {
	//   c.Endpoint.AuthURL = fmt.Sprintf("https://%s/login/oauth/authorize", i.host)
	//   c.Endpoint.TokenURL = fmt.Sprintf("https://%s/login/oauth/access_token", i.host)
	// }

	return c, nil
}
