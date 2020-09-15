package model

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/convox/console/pkg/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID string `dynamo:"id"`

	CliToken        string    `dynamo:"api-key-hash,encrypted" json:"-"`
	CliTokenCreated time.Time `dynamo:"api-key-created-at" json:"-"`
	CliTokenUsed    time.Time `dynamo:"cli-authenticated" json:"-"`
	Deleted         bool      `dynamo:"deleted" json:"-"`
	Email           string    `dynamo:"email"`
	LastActivity    time.Time `dynamo:"last-activity" json:"-"`
	OrganizationIDs []string  `dynamo:"organization-ids"`
	ResetUntil      time.Time `dynamo:"reset-until" json:"-"`
	Superuser       bool      `dynamo:"superuser" json:"-"`

	passwordHash string `dynamo:"password_hash"`
}

type Users []User

func (m *Model) UserAuthenticatePassword(email, password string) (*User, error) {
	var us Users

	if err := m.storage.GetIndex("users", "email-index", map[string]string{"email": email}, &us); err != nil {
		return nil, storage.NotFound("invalid authentication")
	}
	if len(us) < 1 {
		return nil, storage.NotFound("invalid authentication")
	}

	u := us[0]

	if strings.TrimSpace(u.passwordHash) == "" {
		return nil, storage.NotFound("invalid authentication")
	}

	if !u.Authenticate(password) {
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

func (m *Model) UserSave(u *User) error {
	var us Users

	if err := m.storage.GetIndex("users", "email-index", map[string]string{"email": u.Email}, &us); err != nil {
		return err
	}

	fmt.Printf("us: %+v\n", us)

	for _, uu := range us {
		if uu.ID != u.ID {
			return fmt.Errorf("email is already in use")
		}
	}

	if err := m.storage.Put("users", u); err != nil {
		return err
	}

	return nil
}

func (m *Model) UserTokens(uid string) (Tokens, error) {
	var ts Tokens

	if err := m.storage.GetIndex("tokens", "user-id-index", map[string]string{"user-id": uid}, &ts); err != nil {
		return nil, errors.WithStack(err)
	}

	return ts, nil
}

func (m *Model) UserTokensByKind(uid, kind string) (Tokens, error) {
	ts, err := m.UserTokens(uid)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tsk := Tokens{}

	for _, t := range ts {
		if t.Kind == kind {
			tsk = append(tsk, t)
		}
	}

	return tsk, nil
}

func (u *User) Authenticate(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password)); err != nil {
		return false
	}

	return true
}

func (u *User) CliTokenReset() error {
	u.CliToken = strings.ReplaceAll(uuid.New().String(), "-", "")
	u.CliTokenCreated = time.Now().UTC()

	return nil
}

func (u *User) SetPassword(password string) error {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	u.passwordHash = string(data)

	return nil
}

func (us Users) Less(i, j int) bool {
	return us[i].Email < us[j].Email
}
