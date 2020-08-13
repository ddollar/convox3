package model

import (
	"time"

	"github.com/convox/console/pkg/storage"
	"github.com/convox/convox/pkg/options"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Rack struct {
	ID string `dynamo:"id"`

	Creator      string `dynamo:"creator"`
	Install      string `dynamo:"install-id"`
	Organization string `dynamo:"organization-id"`
	Runtime      string `dynamo:"integration-id"`
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

func (m *Model) RackUpdates(id string) (Updates, error) {
	opts := storage.QueryOptions{
		Forward: options.Bool(false),
		Index:   options.String("rack-id-created-index"),
		Limit:   options.Int64(10),
	}

	var us Updates

	if err := m.storage.Query("updates", map[string]string{"rack-id": id}, opts, &us); err != nil {
		return nil, errors.WithStack(err)
	}

	return us, nil
}

func (m *Model) RackSave(r *Rack) error {
	if err := m.storage.Put("racks", r); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Rack) Defaults() {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}

	if r.Created.IsZero() {
		r.Created = time.Now().UTC()
	}
}

func (r *Rack) Validate() []error {
	errs := []error{}

	errs = checkNonzero(errs, r.ID, "id required")
	errs = checkNonzero(errs, r.Creator, "creator required")
	errs = checkNonzero(errs, r.Organization, "organization required")
	errs = checkNonzero(errs, r.Host, "host required")
	errs = checkNonzero(errs, r.Name, "name required")
	errs = checkNonzero(errs, r.Password, "password required")

	return errs
}

func (rs Racks) Less(i, j int) bool {
	return rs[i].Name < rs[j].Name
}
