package settings

import (
	"os"

	"github.com/convox/console/pkg/common"
)

var (
	App              = os.Getenv("APP")
	AuditLogsBucket  = os.Getenv("AUDIT_LOGS_OBJECT_STORE")
	Authentication   = os.Getenv("AUTHENTICATION")
	Development      = os.Getenv("MODE") == "development"
	DiscourseSsoKey  = os.Getenv("DISCOURSE_SSO_KEY")
	ExternalHost     = common.CoalesceString(os.Getenv("TUNNEL_HOST"), os.Getenv("HOST"))
	Host             = os.Getenv("HOST")
	LdapAddr         = os.Getenv("LDAP_ADDR")
	LdapBind         = os.Getenv("LDAP_BIND")
	LdapVerify       = os.Getenv("LDAP_VERIFY")
	LicenseKey       = os.Getenv("LICENSE_KEY")
	RackKey          = os.Getenv("RACK_KEY")
	SamlMetadata     = os.Getenv("SAML_METADATA")
	SegmentClientKey = os.Getenv("SEGMENT_CLIENT_KEY")
	SegmentServerKey = os.Getenv("SEGMENT_SERVER_KEY")
	SessionKey       = os.Getenv("SESSION_KEY")
	StorageProvider  = os.Getenv("STORAGE_PROVIDER")
	StripePublicKey  = os.Getenv("STRIPE_PUBLISHABLE_KEY")
	StripePrivateKey = os.Getenv("STRIPE_SECRET_KEY")
	TablePrefix      = os.Getenv("TABLE_PREFIX")
	Version          = os.Getenv("VERSION")
	WorkerQueue      = os.Getenv("WORKER_QUEUE")

	GithubClientId      = os.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret  = os.Getenv("GITHUB_CLIENT_SECRET")
	GithubWebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")

	GithubEnterpriseClientId     = os.Getenv("GITHUB_ENTERPRISE_CLIENT_ID")
	GithubEnterpriseClientSecret = os.Getenv("GITHUB_ENTERPRISE_CLIENT_SECRET")
	GithubEnterpriseHost         = os.Getenv("GITHUB_ENTERPRISE_HOST")
	GithubEnterpriseInsecure     = os.Getenv("MODE") == "development"

	GitlabClientId     = os.Getenv("GITLAB_CLIENT_ID")
	GitlabClientSecret = os.Getenv("GITLAB_CLIENT_SECRET")

	SlackClientId     = os.Getenv("SLACK_CLIENT_ID")
	SlackClientSecret = os.Getenv("SLACK_CLIENT_SECRET")
)

func AllowLogout() bool {
	switch Authentication {
	case "saml":
		return false
	default:
		return true
	}
}

func AllowSignup() bool {
	switch Authentication {
	case "ldap":
		return false
	case "saml":
		return false
	default:
		return true
	}
}

func BillingEnabled() bool {
	return StripePublicKey != ""
}

func IntegrationEnabled(name string) bool {
	switch name {
	case "github":
		return GithubClientId != ""
	case "github-enterprise":
		return GithubEnterpriseClientId != "" && GithubEnterpriseHost != ""
	case "gitlab":
		return GitlabClientId != ""
	case "slack":
		return SlackClientId != ""
	default:
		return false
	}
}
