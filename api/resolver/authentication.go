package resolver

import (
	"time"

	"github.com/convox/console/api/model"
	"github.com/gbrlsnchs/jwt/v3"
)

type Authentication struct {
	session *model.Session
	user    model.User
}

func (a Authentication) Key() (string, error) {
	data := map[string]string{
		"uid": a.user.ID,
	}

	if a.session != nil {
		data["sid"] = a.session.ID
		data["expires"] = a.session.Expires.Format(time.RFC3339)
	}

	key, err := jwt.Sign(data, jwtHash)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func (a Authentication) Session() *Session {
	if a.session == nil {
		return nil
	}

	return &Session{*a.session}
}

func (a Authentication) User() User {
	return User{a.user}
}
