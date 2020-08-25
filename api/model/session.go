package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Session struct {
	ID string `dynamo:"id" json:"id"`

	UserID string `dynamo:"user-id" json:"user-id"`

	Created time.Time `dynamo:"created" json:"created"`
	Expires time.Time `dynamo:"expires" json:"expires"`

	ttl int64 `dynamo:"ttl" json:"-"`
}

func (m *Model) SessionDelete(id string) error {
	return m.storage.Delete("sessions", id)
}

func (m *Model) SessionGet(id string) (*Session, error) {
	var s Session

	if err := m.storage.Get("sessions", id, &s); err != nil {
		return nil, errors.WithStack(err)
	}

	return &s, nil
}

func (m *Model) SessionSave(s *Session) error {
	return m.storage.Put("sessions", s)
}

func (s *Session) Defaults() {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	if s.Created.IsZero() {
		s.Created = time.Now().UTC()
	}

	if s.Expires.IsZero() {
		s.Expires = time.Now().UTC().Add(60 * time.Minute)
	}

	s.ttl = s.Expires.Unix()
}

func (s *Session) Key() {

}

func (s *Session) Validate() []error {
	errs := []error{}

	errs = checkNonzero(errs, s.ID, "id required")
	errs = checkNonzero(errs, s.UserID, "user-id required")

	return errs
}
