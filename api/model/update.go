package model

import (
	"time"

	"github.com/pkg/errors"
)

type Update struct {
	ID string `dynamo:"id"`

	OrganizationID string `dynamo:"organization-id"`
	RackID         string `dynamo:"rack-id"`

	Created  time.Time         `dynamo:"created"`
	Finished time.Time         `dynamo:"finished"`
	Params   map[string]string `dynamo:"params,encrypted"`
	Pid      string            `dynamo:"pid"`
	Started  time.Time         `dynamo:"started"`
	Status   string            `dynamo:"status"`
	Version  string            `dynamo:"version"`
}

type Updates []Update

func (m *Model) UpdateGet(id string) (*Update, error) {
	u := &Update{}

	if err := m.storage.Get("updates", id, u); err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}

func (m *Model) UpdateSave(u *Update) error {
	if err := m.storage.Put("updates", u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
