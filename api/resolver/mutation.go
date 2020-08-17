package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
)

type InstanceTerminateArgs struct {
	Oid graphql.ID
	Rid graphql.ID
	Iid graphql.ID
}

func (r *Root) InstanceTerminate(ctx context.Context, args InstanceTerminateArgs) (string, error) {
	rr, err := authenticatedRack(ctx, r.model, string(args.Oid), string(args.Rid))
	if err != nil {
		return "", err
	}

	c, err := rackClient(ctx, rr.Host, rr.Password)
	if err != nil {
		return "", err
	}

	if err := c.InstanceTerminate(string(args.Iid)); err != nil {
		return "", err
	}

	return string(args.Iid), nil
}

type LoginArgs struct {
	Email    string
	Password string
}

func (r *Root) Login(ctx context.Context, args LoginArgs) (*Authentication, error) {
	mu, err := r.model.UserAuthenticatePassword(args.Email, args.Password)
	if err != nil {
		return nil, err
	}

	u := User{
		id:    mu.ID,
		email: mu.Email,
	}

	a := &Authentication{
		user: u,
	}

	return a, nil
}

type ProcessStopArgs struct {
	Oid graphql.ID
	Rid graphql.ID
	App string
	Pid graphql.ID
}

func (r *Root) ProcessStop(ctx context.Context, args ProcessStopArgs) (string, error) {
	rr, err := authenticatedRack(ctx, r.model, string(args.Oid), string(args.Rid))
	if err != nil {
		return "", err
	}

	c, err := rackClient(ctx, rr.Host, rr.Password)
	if err != nil {
		return "", err
	}

	if err := c.ProcessStop(args.App, string(args.Pid)); err != nil {
		return "", err
	}

	return string(args.Pid), nil
}

type RackImportArgs struct {
	Oid      graphql.ID
	Name     string
	Hostname string
	Password string
}

func (r *Root) RackImport(ctx context.Context, args RackImportArgs) (*Rack, error) {
	u, err := currentUser(ctx)
	if err != nil {
		return nil, err
	}

	o, err := authenticatedOrganization(ctx, r.model, string(args.Oid))
	if err != nil {
		return nil, err
	}

	rs, err := r.model.OrganizationRacks(o.ID)
	if err != nil {
		return nil, err
	}

	for _, r := range rs {
		if r.Name == args.Name {
			return nil, fmt.Errorf("rack name already exists")
		}
	}

	rr := model.Rack{
		Creator:      u.id,
		Organization: o.ID,
		Name:         args.Name,
		Host:         args.Hostname,
		Password:     args.Password,
	}

	if err := r.model.RackSave(&rr); err != nil {
		return nil, errors.WithStack(err)
	}

	return &Rack{rr}, nil
}

type RackRemoveArgs struct {
	Oid graphql.ID
	Id  graphql.ID
}

func (r *Root) RackRemove(ctx context.Context, args RackRemoveArgs) (string, error) {
	rr, err := authenticatedRack(ctx, r.model, string(args.Oid), string(args.Id))
	if err != nil {
		return "", err
	}

	if err := r.model.RackDelete(rr.ID); err != nil {
		return "", err
	}

	return rr.ID, nil
}

type RackUpdateArgs struct {
	Oid             graphql.ID
	Id              graphql.ID
	Name            string
	Runtime         string
	UpdateDay       int32
	UpdateFrequency string
	UpdateHour      int32
}

func (r *Root) RackUpdate(ctx context.Context, args RackUpdateArgs) (string, error) {
	rr, err := authenticatedRack(ctx, r.model, string(args.Oid), string(args.Id))
	if err != nil {
		return "", err
	}

	rr.Name = args.Name
	rr.UpdateDay = int(args.UpdateDay)
	rr.UpdateFrequency = args.UpdateFrequency
	rr.UpdateHour = int(args.UpdateHour)

	if args.Runtime != "" {
		i, err := r.model.IntegrationGet(args.Runtime)
		if err != nil {
			return "", err
		}

		if i.OrganizationId != string(args.Oid) {
			return "", fmt.Errorf("invalid runtime")
		}

		rr.Runtime = args.Runtime
	} else {
		rr.Runtime = ""
	}

	if err := r.model.RackSave(rr); err != nil {
		return "", err
	}

	return rr.ID, nil
}
