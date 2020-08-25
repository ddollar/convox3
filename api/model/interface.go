package model

import "io"

type Interface interface {
	ChallengeGet(id string) (*Challenge, error)
	ChallengeSave(c *Challenge) error
	InstallGet(id string) (*Install, error)
	InstallLogs(id string) (io.ReadCloser, error)
	InstallSave(i *Install) error
	IntegrationGet(iid string) (*Integration, error)
	JobFail(id string, err error) error
	JobGet(id string) (*Job, error)
	JobListByStatus(status string) (Jobs, error)
	JobSave(j *Job) error
	JobSourceStatus(id, status string) error
	JobSucceed(id string) error
	OrganizationGet(id string) (*Organization, error)
	OrganizationIntegrations(oid string) (Integrations, error)
	OrganizationRacks(oid string) (Racks, error)
	OrganizationSave(o *Organization) error
	RackDelete(id string) error
	RackGet(id string) (*Rack, error)
	RackLock(id string) error
	RackSave(r *Rack) error
	RackStateLoad(id string) ([]byte, error)
	RackStateStore(id string, data []byte) error
	RackUnlock(id string) error
	SessionDelete(id string) error
	SessionGet(id string) (*Session, error)
	SessionSave(s *Session) error
	TaskWriter(t Task) (*TaskWriter, error)
	UninstallGet(id string) (*Uninstall, error)
	UninstallLogs(id string) (io.ReadCloser, error)
	UninstallSave(u *Uninstall) error
	UpdateGet(id string) (*Update, error)
	UpdateSave(u *Update) error
	UserAuthenticatePassword(email, password string) (*User, error)
	UserGet(id string) (*User, error)
	UserGetBatch(ids []string) (Users, error)
	UserOrganizations(uid string) (Organizations, error)
	UserTokens(uid string) (Tokens, error)
	UserTokensByKind(uid, kind string) (Tokens, error)
	TokenDelete(id string) error
	TokenGet(id string) (*Token, error)
	TokenSave(token *Token) error
	WorkflowGet(id string) (*Workflow, error)
}
