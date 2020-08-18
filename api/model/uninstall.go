package model

import (
	"time"

	"github.com/pkg/errors"
)

type Uninstall struct {
	ID string `dynamo:"id"`

	OrganizationID string `dynamo:"organization-id"`
	RackID         string `dynamo:"rack-id"`

	Created  time.Time `dynamo:"created"`
	Engine   string    `dynamo:"engine"`
	Finished time.Time `dynamo:"finished"`
	Pid      string    `dynamo:"pid"`
	Started  time.Time `dynamo:"started"`
	Status   string    `dynamo:"status"`
	Version  string    `dynamo:"version"`
}

type Uninstalls []Uninstall

func (m *Model) UninstallGet(id string) (*Uninstall, error) {
	u := &Uninstall{}

	if err := m.storage.Get("uninstalls", id, u); err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}

func (m *Model) UninstallSave(u *Uninstall) error {
	if err := m.storage.Put("uninstalls", u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
