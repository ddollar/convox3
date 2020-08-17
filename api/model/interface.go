package model

type Interface interface {
	IntegrationGet(iid string) (*Integration, error)
	OrganizationGet(id string) (*Organization, error)
	OrganizationIntegrations(oid string) (Integrations, error)
	OrganizationRacks(oid string) (Racks, error)
	OrganizationSave(o *Organization) error
	RackDelete(id string) error
	RackGet(id string) (*Rack, error)
	RackSave(r *Rack) error
	UserAuthenticatePassword(email, password string) (*User, error)
	UserGet(id string) (*User, error)
	UserGetBatch(ids []string) (Users, error)
	UserOrganizations(uid string) (Organizations, error)
}
