package integration

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
)

type Authorizer interface {
	Authorize() (string, error)
	Exchange(code string, reauthorize bool) (string, map[string]string, error)
	Revoke() error
}

type Integration interface {
	Name() string
	Slug() string
	Status() (string, error)
	Title(attrs map[string]string) (string, error)
}

type Notification interface {
	EventSend(string, map[string]string, NotificationEvent) error
}

type Runtime interface {
	Credentials() (map[string]string, error)
	ParameterList() ([]string, error)
	RegionList() ([]string, error)
}

type Source interface {
	RepositoryClone(int64, string, io.Writer) (string, error)
	RepositoryList() (map[int64]string, error)
	RepositoryName(int64) (string, error)
	StatusUpdate(int64, string, string, string, string) error
	WebhookCreateMerge(int64, string, string) (int64, error)
	WebhookCreateReview(int64, string) (int64, error)
	WebhookDelete(int64, int64) error
	WebhookPayload(r *http.Request, validate bool) (*SourceWebhookPayload, error)
}

type NotificationEvent struct {
	Action    string            `json:"action"`
	Data      map[string]string `json:"data"`
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
}

type SourceWebhookPayload struct {
	Description string
	Event       string
	Name        string
	Ref         string
	Title       string
	URL         string
}

func New(id, oid, provider, token string) (Integration, error) {
	switch provider {
	case "aws":
		return &AWS{id: id, oid: oid, token: token}, nil
	case "azure":
		return &Azure{id: id, oid: oid, token: token}, nil
	case "do":
		return &DigitalOcean{id: id, oid: oid, token: token}, nil
	case "gcp":
		return &Google{id: id, oid: oid, token: token}, nil
	case "gitlab":
		return &Gitlab{client_id: settings.GitlabClientId, client_secret: settings.GitlabClientSecret, id: id, oid: oid, token: token}, nil
	case "github":
		return &Github{client_id: settings.GithubClientId, client_secret: settings.GithubClientSecret, id: id, oid: oid, provider: "github", token: token}, nil
	case "github-enterprise":
		return &Github{client_id: settings.GithubEnterpriseClientId, client_secret: settings.GithubEnterpriseClientSecret, host: settings.GithubEnterpriseHost, id: id, oid: oid, provider: "github-enterprise", token: token}, nil
	case "slack":
		return &Slack{client_id: settings.SlackClientId, client_secret: settings.SlackClientSecret, oid: oid, token: token}, nil
	default:
		return nil, errors.WithStack(fmt.Errorf("unknown provider: %s", provider))
	}
}

var (
	_ Notification = &Slack{}
	_ Runtime      = &AWS{}
	_ Runtime      = &Azure{}
	_ Runtime      = &DigitalOcean{}
	_ Runtime      = &Google{}
	_ Source       = &Github{}
	_ Source       = &Gitlab{}
)
