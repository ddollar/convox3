package resolver

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/graph-gophers/graphql-go"
)

const (
	
)
var jwtHash = jwt.NewHS256([]byte("secret"))

type User struct {
	id    string
	email string
}

func UserFromToken(token string) (*User, error) {
	var data map[string]string

	if _, err := jwt.Verify([]byte(token), jwtHash, &data); err != nil {
		return nil, err
	}

	u := &User{
		id:    data["id"],
		email: data["email"],
	}

	return u, nil
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