package model

type Interface interface {
	OrganizationGet(id string) (*Organization, error)
	OrganizationIntegrations(oid string) (Integrations, error)
	OrganizationRacks(oid string) (Racks, error)
	RackGet(id string) (*Rack, error)
	UserAuthenticatePassword(email, password string) (*User, error)
	UserGet(id string) (*User, error)
	UserGetBatch(ids []string) (Users, error)
	UserOrganizations(uid string) (Organizations, error)
}
