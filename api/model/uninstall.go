package model

import (
	"fmt"
	"io"
	"time"

	"github.com/convox/console/pkg/settings"
	"github.com/google/uuid"
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

func (m *Model) UninstallFail(id string, failure error) error {
	u, err := m.UninstallGet(id)
	if err != nil {
		return err
	}

	u.Finished = time.Now().UTC()
	u.Status = "failed"

	if err := m.UninstallSave(u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) UninstallGet(id string) (*Uninstall, error) {
	u := &Uninstall{}

	if err := m.storage.Get("uninstalls", id, u); err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}

func (m *Model) UninstallLogs(id string) (io.ReadCloser, error) {
	u, err := m.UninstallGet(id)
	if err != nil {
		return nil, err
	}

	r, err := m.rack.ObjectFetch(settings.App, u.Key())
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *Model) UninstallSave(u *Uninstall) error {
	if err := m.storage.Put("uninstalls", u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) UninstallSucceed(id string) error {
	u, err := m.UninstallGet(id)
	if err != nil {
		return err
	}

	u.Finished = time.Now().UTC()
	u.Status = "complete"

	if err := m.UninstallSave(u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *Uninstall) Defaults() {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	if u.Created.IsZero() {
		u.Created = time.Now().UTC()
	}

	if u.Status == "" {
		u.Status = "pending"
	}
}

func (u *Uninstall) Key() string {
	return fmt.Sprintf("organizations/%s/racks/%s/uninstalls/%s", u.OrganizationID, u.RackID, u.ID)
}
