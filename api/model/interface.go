package model

type Interface interface {
	InstallGet(id string) (*Install, error)
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
	RackSave(r *Rack) error
	TaskWriter(t Task) (*TaskWriter, error)
	UninstallGet(id string) (*Uninstall, error)
	UninstallSave(u *Uninstall) error
	UpdateGet(id string) (*Update, error)
	UpdateSave(u *Update) error
	UserAuthenticatePassword(email, password string) (*User, error)
	UserGet(id string) (*User, error)
	UserGetBatch(ids []string) (Users, error)
	UserOrganizations(uid string) (Organizations, error)
	WorkflowGet(id string) (*Workflow, error)
}
