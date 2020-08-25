package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Token struct {
	ID string `dynamo:"id" json:"id"`

	UserID string `dynamo:"user-id" json:"user-id"`

	Counter int       `dynamo:"counter" json:"counter"`
	Data    []byte    `dynamo:"data", json:"data"`
	Kind    string    `dynamo:"kind" json:"kind"`
	Name    string    `dynamo:"name" json:"name"`
	Used    time.Time `dynamo:"used" json:"used"`
}

type Tokens []Token

func (m *Model) TokenDelete(id string) error {
	return m.storage.Delete("tokens", id)
}

func (m *Model) TokenGet(id string) (*Token, error) {
	var t Token

	if err := m.storage.Get("tokens", id, &t); err != nil {
		return nil, errors.WithStack(err)
	}

	return &t, nil
}

func (m *Model) TokenSave(t *Token) error {
	return m.storage.Put("tokens", t)
}

func (t *Token) Defaults() {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	if t.Name == "" {
		t.Name = t.ID[0:10]
	}
}

func (t *Token) Validate() []error {
	errs := []error{}

	errs = checkNonzero(errs, t.ID, "id required")
	errs = checkNonzero(errs, t.UserID, "user-id required")
	errs = checkNonzero(errs, t.Data, "data required")
	errs = checkNonzero(errs, t.Kind, "kind required")
	errs = checkNonzero(errs, t.Name, "name required")

	return errs
}
