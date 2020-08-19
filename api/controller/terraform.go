package controller

import (
	"fmt"
	"io/ioutil"

	"github.com/convox/stdapi"
	"github.com/pkg/errors"
)

func (cn *Controller) TerraformLock(c *stdapi.Context) error {
	rid := c.Var("rid")

	if err := cn.model.RackLock(rid); err != nil {
		return stdapi.Errorf(423, "could not lock rack")
	}

	return c.RenderOK()
}

func (cn *Controller) TerraformStateLoad(c *stdapi.Context) error {
	rid := c.Var("rid")
	oid := c.Var("oid")

	r, err := cn.model.RackGet(rid)
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not fetch rack"))
	}
	if r.Organization != oid {
		return errors.WithStack(fmt.Errorf("invalid organization"))
	}

	state, err := cn.model.RackStateLoad(rid)
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not load state"))
	}

	c.Write(state)

	return nil
}

func (cn *Controller) TerraformStateStore(c *stdapi.Context) error {
	rid := c.Var("rid")
	oid := c.Var("oid")

	r, err := cn.model.RackGet(rid)
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not fetch rack"))
	}
	if r.Organization != oid {
		return errors.WithStack(fmt.Errorf("invalid organization"))
	}

	data, err := ioutil.ReadAll(c.Body())
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not read state"))
	}

	if err := cn.model.RackStateStore(rid, data); err != nil {
		return err
	}

	return c.RenderOK()
}

func (cn *Controller) TerraformUnlock(c *stdapi.Context) error {
	rid := c.Var("rid")

	if err := cn.model.RackUnlock(rid); err != nil {
		return stdapi.Errorf(403, "could not unlock rack")
	}

	return nil
}
