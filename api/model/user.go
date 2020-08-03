package model

import (
	"sort"
	"strings"

	"github.com/convox/console/pkg/storage"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID string `dynamo:"id"`

	Email           string   `dynamo:"email"`
	OrganizationIDs []string `dynamo:"organization-ids"`

	passwordHash string `dynamo:"password_hash" json:"-"`
}

type Users []User

func (m *Model) UserAuthenticatePassword(email, password string) (*User, error) {
	var us Users

	if err := m.storage.GetIndex("users", "email-index", map[string]string{"email": email}, &us); err != nil {
		return nil, errors.WithStack(err)
	}
	if len(us) < 1 {
		return nil, storage.NotFound("invalid authentication")
	}

	u := us[0]

	if strings.TrimSpace(u.passwordHash) == "" {
		return nil, storage.NotFound("invalid authentication")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)); err != nil {
		return nil, storage.NotFound("invalid authentication")
	}

	return &u, nil
}

func (m *Model) UserGet(id string) (*User, error) {
	u := &User{}

	if err := m.storage.Get("users", id, u); err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}

func (m *Model) UserGetBatch(ids []string) (Users, error) {
	us := Users{}

	if err := m.storage.GetBatch("users", ids, &us); err != nil {
		return nil, err
	}

	sort.Slice(us, us.Less)

	return us, nil
}

func (m *Model) UserOrganizations(uid string) (Organizations, error) {
	u, err := m.UserGet(uid)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(u.OrganizationIDs) == 0 {
		return Organizations{}, nil
	}

	var os Organizations

	if err := m.storage.GetBatch("organizations", u.OrganizationIDs, &os); err != nil {
		return nil, errors.WithStack(err)
	}

	sort.Slice(os, os.Less)

	return os, nil
}

func (us Users) Less(i, j int) bool {
	return us[i].Email < us[j].Email
}
