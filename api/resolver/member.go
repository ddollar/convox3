package resolver

import "github.com/convox/console/api/model"

type Member struct {
	access string
	user   model.User
}

func (m *Member) Access() string {
	return m.access
}

func (m *Member) User() *User {
	return &User{id: m.user.ID, email: m.user.Email}
}
