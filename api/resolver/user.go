package resolver

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/graph-gophers/graphql-go"
)

type User struct {
	id    string
	email string
}

func (u User) Id() graphql.ID {
	return graphql.ID(u.id)
}

func (u User) Email() string {
	return u.email
}

func (u User) Token() (string, error) {
	data := map[string]string{
		"id":    u.id,
		"email": u.email,
	}

	token, err := jwt.Sign(data, jwtHash)
	if err != nil {
		return "", err
	}

	return string(token), nil
}
