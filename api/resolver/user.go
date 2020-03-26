package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type User struct {
	model.User
}

func (u *User) Id() graphql.ID {
	return graphql.ID(u.User.ID)
}

func (u *User) Email() string {
	return u.User.Email
}
