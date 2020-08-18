package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Install struct {
	ID string `dynamo:"id"`

	Backend  string            `dynamo:"backend,encrypted"`
	Engine   string            `dynamo:"engine"`
	Name     string            `dynamo:"name"`
	Params   map[string]string `dynamo:"params"`
	Pid      string            `dynamo:"pid"`
	Progress int               `dynamo:"progress"`
	Provider string            `dynamo:"provider"`
	Region   string            `dynamo:"region"`
	Status   string            `dynamo:"status"`
	Version  string            `dynamo:"version"`

	OrganizationID string `dynamo:"organization-id"`
	RackID         string `dynamo:"rack-id"`
	UserID         string `dynamo:"user-id"`

	Created  time.Time `dynamo:"created"`
	Started  time.Time `dynamo:"started"`
	Finished time.Time `dynamo:"finished"`
}

type Jobs []Job

func (m *Model) InstallGet(id string) (*Install, error) {
	i := &Install{}

	if err := m.storage.Get("installs", id, i); err != nil {
		return nil, errors.WithStack(err)
	}

	return i, nil
}

func (m *Model) InstallSave(i *Install) error {
	if err := m.storage.Put("installs", i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (i *Install) Defaults() {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}

	if i.Params == nil {
		i.Params = map[string]string{}
	}

	if i.Status == "" {
		i.Status = "pending"
	}

	if i.Created.IsZero() {
		i.Created = time.Now().UTC()
	}
}
