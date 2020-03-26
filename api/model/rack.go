package model

import (
	"time"

	"github.com/pkg/errors"
)

type Rack struct {
	ID string `dynamo:"id"`

	Creator      string `dynamo:"creator"`
	Install      string `dynamo:"install-id"`
	Integration  string `dynamo:"integration-id"`
	Organization string `dynamo:"organization-id"`
	Uninstall    string `dynamo:"uninstall-id"`

	Created          time.Time         `dynamo:"created"`
	Host             string            `dynamo:"host"`
	Locked           bool              `dynamo:"locked"`
	Name             string            `dynamo:"name"`
	Parameters       map[string]string `dynamo:"parameters"`
	Password         string            `dynamo:"password,encrypted"`
	Provider         string            `dynamo:"provider"`
	Region           string            `dynamo:"region"`
	Stack            string            `dynamo:"stack"`
	UnreachableCount int               `dynamo:"unreachable-count"`
	UpdateDay        int               `dynamo:"update-day"`
	UpdateFrequency  string            `dynamo:"update-frequency"`
	UpdateHour       int               `dynamo:"update-hour"`
	UpdateNext       time.Time         `dynamo:"update-next"`
}

type Racks []Rack

func (m *Model) RackGet(id string) (*Rack, error) {
	r := &Rack{}

	if err := m.storage.Get("racks", id, r); err != nil {
		return nil, errors.WithStack(err)
	}

	return r, nil
}

func (rs Racks) Less(i, j int) bool {
	return rs[i].Name < rs[j].Name
}
