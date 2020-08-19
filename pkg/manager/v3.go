package manager

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/convox/console/pkg/helpers"
	"github.com/convox/console/pkg/terraform"
	"github.com/pkg/errors"
)

type V3 struct {
	backend  string
	env      map[string]string
	id       string
	provider string
}

func NewV3(rid string) (*V3, error) {
	r, err := m.RackGet(rid)
	if err != nil {
		return nil, err
	}

	backend, err := r.TerraformBackend()
	if err != nil {
		return nil, err
	}

	ri, err := m.RackRuntime(rid)
	if err != nil {
		return nil, err
	}

	creds, err := ri.Credentials()
	if err != nil {
		return nil, err
	}

	m := &V3{
		backend:  backend,
		env:      creds,
		id:       r.ID,
		provider: r.Provider,
	}

	return m, nil
}

func (v *V3) Install(name, version, region string, params map[string]string, output io.Writer) error {
	t, err := v.terraform(output)
	if err != nil {
		return err
	}

	params["region"] = region

	if err := t.Install(name, version, v.env, params); err != nil {
		return errors.WithStack(err)
	}

	if err := v.sync(v.id, params); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (v *V3) Uninstall(output io.Writer) error {
	r, err := m.RackGet(v.id)
	if err != nil {
		return err
	}

	t, err := v.terraform(output)
	if err != nil {
		return err
	}

	version := helpers.CoalesceString(r.Parameters["release"], "master")

	if err := t.Uninstall(r.Name, version, v.env, r.Parameters); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (v *V3) Update(name, version string, params map[string]string, output io.Writer) error {
	t, err := v.terraform(output)
	if err != nil {
		return err
	}

	if err := t.Update(name, version, v.env, params); err != nil {
		return err
	}

	if err := v.sync(v.id, params); err != nil {
		return err
	}

	return nil
}

func (v *V3) sync(rid string, params map[string]string) error {
	r, err := m.RackGet(rid)
	if err != nil {
		return errors.WithStack(err)
	}

	data, err := m.RackStateLoad(rid)
	if err != nil {
		return err
	}

	os, err := stateOutputs(data)
	if err != nil {
		return err
	}

	u, err := url.Parse(os["api"])
	if err != nil {
		return errors.WithStack(err)
	}

	r.Host = u.Host

	if pw, ok := u.User.Password(); ok {
		r.Password = pw
	}

	if r.Parameters == nil {
		r.Parameters = map[string]string{}
	}

	for k, v := range params {
		r.Parameters[k] = v
	}

	r.Parameters["release"] = os["release"]

	if err := m.RackSave(r); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (v *V3) terraform(output io.Writer) (*terraform.Terraform, error) {
	t, err := terraform.New(v.provider, v.backend)
	if err != nil {
		return nil, err
	}

	t.Writer = output

	return t, nil
}

func stateOutputs(data []byte) (map[string]string, error) {
	var state struct {
		Outputs map[string]struct {
			Value string
		}
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, errors.WithStack(err)
	}

	os := map[string]string{}

	for k, o := range state.Outputs {
		os[k] = o.Value
	}

	return os, nil
}
