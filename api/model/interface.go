package model

type Interface interface {
	OrganizationGet(id string) (*Organization, error)
	OrganizationRacks(oid string) (Racks, error)
	RackGet(id string) (*Rack, error)
	UserGet(id string) (*User, error)
	UserGetBatch(ids []string) (Users, error)
	UserOrganizations(uid string) (Organizations, error)
}
