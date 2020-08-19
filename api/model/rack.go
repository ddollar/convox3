package model

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/convox/console/pkg/crypt"
	"github.com/convox/console/pkg/integration"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/console/pkg/storage"
	"github.com/convox/convox/pkg/options"
	"github.com/convox/convox/pkg/structs"
	"github.com/convox/convox/sdk"
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

	Created         time.Time         `dynamo:"created"`
	Host            string            `dynamo:"host"`
	Locked          bool              `dynamo:"locked"`
	Name            string            `dynamo:"name"`
	Parameters      map[string]string `dynamo:"parameters"`
	Password        string            `dynamo:"password,encrypted"`
	Provider        string            `dynamo:"provider"`
	Region          string            `dynamo:"region"`
	Stack           string            `dynamo:"stack"`
	UpdateDay       int               `dynamo:"update-day"`
	UpdateFrequency string            `dynamo:"update-frequency"`
	UpdateHour      int               `dynamo:"update-hour"`
	UpdateNext      time.Time         `dynamo:"update-next"`
}

type Racks []Rack

func (m *Model) RackDelete(id string) error {
	if err := m.storage.Delete("racks", id); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) RackGet(id string) (*Rack, error) {
	r := &Rack{}

	if err := m.storage.Get("racks", id, r); err != nil {
		return nil, errors.WithStack(err)
	}

	return r, nil
}

func (m *Model) RackIntegration(id string) (*Integration, error) {
	r, err := m.RackGet(id)
	if err != nil {
		return nil, err
	}

	if r.Runtime == "" {
		return nil, nil
	}

	return m.IntegrationGet(r.Runtime)
}

func (m *Model) RackLock(id string) error {
	r, err := m.RackGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if r.Locked {
		return errors.WithStack(fmt.Errorf("rack is already locked"))
	}

	r.Locked = true

	if err := m.RackSave(r); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) RackRuntime(id string) (integration.Runtime, error) {
	i, err := m.RackIntegration(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if i == nil {
		return nil, nil
	}

	ii, err := i.Integration()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ir, ok := ii.(integration.Runtime)
	if !ok {
		return nil, errors.WithStack(fmt.Errorf("not a runtime integration"))
	}

	return ir, nil
}

func (m *Model) RackStateLoad(id string) ([]byte, error) {
	r, err := m.RackGet(id)
	if err != nil {
		return nil, err
	}

	rr, err := m.rack.ObjectFetch(settings.App, r.stateKey())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}
	defer rr.Close()

	data, err := ioutil.ReadAll(rr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (m *Model) RackStateStore(id string, state []byte) error {
	r, err := m.RackGet(id)
	if err != nil {
		return err
	}

	if _, err := m.rack.ObjectStore(settings.App, r.stateKey(), bytes.NewReader(state), structs.ObjectStoreOptions{}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (m *Model) RackUnlock(id string) error {
	r, err := m.RackGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if !r.Locked {
		return nil
	}

	r.Locked = false

	if err := m.RackSave(r); err != nil {
		return errors.WithStack(err)
	}

	return nil
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

func (r *Rack) System() (*structs.System, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	s, err := c.SystemGet()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *Rack) TerraformBackend() (string, error) {
	pw, err := crypt.Encrypt(settings.RackKey, []byte(r.ID))
	if err != nil {
		return "", errors.WithStack(err)
	}

	u := url.URL{
		Scheme: "https",
		Host:   settings.Host,
		Path:   fmt.Sprintf("/organizations/%s/racks/%s/terraform", r.Organization, r.ID),
		User:   url.UserPassword("terraform", pw),
	}

	return u.String(), nil
}

func (r *Rack) URL() (string, error) {
	return fmt.Sprintf("https://convox:%s@%s", r.Password, r.Host), nil
}

func (r *Rack) Validate() []error {
	errs := []error{}

	errs = checkNonzero(errs, r.ID, "id required")
	errs = checkNonzero(errs, r.Organization, "organization required")
	errs = checkNonzero(errs, r.Name, "name required")

	return errs
}

func (r *Rack) client(ctx context.Context) (*sdk.Client, error) {
	url, err := r.URL()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s, err := sdk.New(url)
	if err != nil {
		return nil, err
	}

	s.Client = s.Client.WithContext(ctx)

	return s, nil
}

func (r *Rack) stateKey() string {
	return fmt.Sprintf("organizations/%s/racks/%s/state", r.Organization, r.ID)
}

func (rs Racks) Less(i, j int) bool {
	return rs[i].Name < rs[j].Name
}
