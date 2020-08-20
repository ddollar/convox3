package resolver

import (
	"context"
	"fmt"
	"regexp"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/queue"
	"github.com/convox/console/pkg/settings"
	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
)

var (
	reRackName = regexp.MustCompile(`^[a-z0-9-]$`)
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

	return &Rack{Rack: rr, model: r.model}, nil
}

type RackInstallArgs struct {
	Oid        graphql.ID
	Iid        graphql.ID
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

	u, err := currentUser(ctx)
	if err != nil {
		return "", err
	}

	i, err := authenticatedIntegration(ctx, r.model, string(args.Oid), string(args.Iid))
	if err != nil {
		return "", err
	}

	ii, err := i.Integration()
	if err != nil {
		return "", errors.WithStack(err)
	}

	status, err := ii.Status()
	if err != nil {
		return "", errors.WithStack(err)
	}
	if status != "connected" {
		return "", fmt.Errorf("runtime is not connected")
	}

	name := reRackName.ReplaceAllString(args.Name, "")

	rs, err := r.model.OrganizationRacks(string(args.Oid))
	if err != nil {
		return "", errors.WithStack(err)
	}

	for _, r := range rs {
		if r.Name == name {
			return "", fmt.Errorf("rack name in use")
		}
	}

	rr := &model.Rack{
		Creator:      u.id,
		Name:         name,
		Organization: string(args.Oid),
		Provider:     i.Provider,
		Region:       args.Region,
		Runtime:      string(args.Iid),
		Stack:        name,
	}

	if err := r.model.RackSave(rr); err != nil {
		return "", err
	}

	backend, err := rr.TerraformBackend()
	if err != nil {
		return "", err
	}

	in := &model.Install{
		Backend:        backend,
		Engine:         args.Engine,
		Name:           args.Name,
		OrganizationID: string(args.Oid),
		Params:         map[string]string{},
		Provider:       i.Provider,
		RackID:         rr.ID,
		Region:         args.Region,
		UserID:         u.id,
	}

	for _, p := range args.Parameters {
		in.Params[p.Key] = p.Value
	}

	if err := r.model.InstallSave(in); err != nil {
		return "", err
	}

	rr.Install = in.ID

	if err := r.model.RackSave(rr); err != nil {
		return "", err
	}

	work := map[string]string{
		"id":   in.ID,
		"type": "install",
	}

	if err := queue.New(settings.WorkerQueue).Enqueue(in.ID, string(args.Oid), work); err != nil {
		return "", errors.WithStack(err)
	}

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

	uninstall, err := rackUninstallable(r.model, rr.ID)
	if err != nil {
		return "", err
	}

	if uninstall {
		if err := rackUninstall(ctx, r.model, string(args.Oid), rr.ID); err != nil {
			return "", err
		}
	} else {
		if err := r.model.RackDelete(rr.ID); err != nil {
			return "", err
		}
	}

	return rr.ID, nil
}

func rackUninstall(ctx context.Context, m model.Interface, oid, rid string) error {
	r, err := authenticatedRack(ctx, m, oid, rid)
	if err != nil {
		return err
	}

	if r.Runtime == "" {
		return fmt.Errorf("rack is not associated with a runtime integration")
	}

	u := &model.Uninstall{
		Engine:         "v3",
		OrganizationID: oid,
		RackID:         rid,
	}

	if s, err := r.System(); err == nil {
		u.Version = s.Version
	}

	if err := m.UninstallSave(u); err != nil {
		return err
	}

	rr, err := m.RackGet(rid)
	if err != nil {
		return err
	}

	rr.Uninstall = u.ID

	if err := m.RackSave(rr); err != nil {
		return err
	}

	work := map[string]string{
		"id":   u.ID,
		"type": "uninstall",
	}

	if err := queue.New(settings.WorkerQueue).Enqueue(u.ID, r.Organization, work); err != nil {
		return errors.WithStack(err)
	}

	return nil
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
		if _, err := authenticatedIntegration(ctx, r.model, string(args.Oid), string(args.Runtime)); err != nil {
			return "", err
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
