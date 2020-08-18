package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/helpers"
	"github.com/convox/console/pkg/queue"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/stdapi"
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

type RackInstallArgs struct {
	Oid        graphql.ID
	Runtime    graphql.ID
	Name       string
	Engine     string
	Region     string
	Parameters []*ParameterArg
}

func (r *Root) RackInstall(ctx context.Context, args RackInstallArgs) (string, error) {
	errs := []error{}

	errs = checkNonzero(errs, args.Name, "name required")
	errs = checkNonzero(errs, args.Engine, "engine required")
	errs = checkNonzero(errs, args.Region, "region required")

	if len(errs) > 0 {
		return "", collateErrors(errs)
	}

	fmt.Printf("args: %+v\n", args)

	for _, p := range args.Parameters {
		fmt.Printf("p: %+v\n", p)
	}

	i, err := cn.Models.IntegrationGet(args.Runtime)
	if err != nil {
		return errors.WithStack(err)
	}

	if i.OrganizationId != oid {
		return errors.WithStack(fmt.Errorf("invalid organization"))
	}

	ii, err := i.Integration()
	if err != nil {
		return errors.WithStack(err)
	}

	status, err := ii.Status()
	if err != nil {
		return errors.WithStack(err)
	}
	if status != "connected" {
		return stdapi.Errorf(403, "Integration is not connected")
	}

	name := c.Form("name")
	region := c.Form("region")

	name = reRackName.ReplaceAllString(name, "")

	rs, err := cn.Models.OrganizationRacks(oid)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, r := range rs {
		if r.Name == name {
			return stdapi.Errorf(403, "Duplicate Rack name")
		}
	}

	r, err := models.NewRack(oid, name)
	if err != nil {
		return errors.WithStack(err)
	}

	r.CreatorID = uid
	r.IntegrationID = i.Id
	r.Region = region
	r.Stack = name

	r.SetProvider(ii.Slug())

	if err := cn.Models.RackSave(r); err != nil {
		return errors.WithStack(err)
	}

	in, err := models.NewInstall(oid, uid, r.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	names := c.Request().PostForm["parameter-name"]
	values := c.Request().PostForm["parameter-value"]

	if len(names) == len(values) {
		for i := range names {
			in.Params[names[i]] = values[i]
		}
	}

	provider := ii.Slug()

	backend, err := r.TerraformBackend()
	if err != nil {
		return errors.WithStack(err)
	}

	in.Backend = backend
	in.Engine = helpers.CoalesceString(c.Form("engine"), "v3")
	in.Name = name
	in.Provider = provider
	in.Region = region
	in.Version = helpers.CoalesceString(in.Params["release"], LatestRelease)

	if err := cn.Models.InstallSave(in); err != nil {
		return errors.WithStack(err)
	}

	r.InstallID = in.ID

	if err := cn.Models.RackSave(r); err != nil {
		return errors.WithStack(err)
	}

	work := map[string]string{
		"id":   in.ID,
		"type": "install",
	}

	if err := queue.New(settings.WorkerQueue).Enqueue(in.ID, oid, work); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(302, fmt.Sprintf("/organizations/%s/racks", oid))

	return "", nil
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
