package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/common"
	"github.com/graph-gophers/graphql-go"
)

type Root struct {
	model model.Interface
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

type SignupArgs struct {
	Email    string
	Password string
}

// func (r *Root) Signup(ctx context.Context, args SignupArgs) (*Authentication, error) {
// 	mu, err := cmodel(ctx).UserCreate(args.Email, args.Password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	u := User{
// 		id:    mu.Id,
// 		email: mu.Email,
// 	}

// 	a := &Authentication{
// 		user: u,
// 	}

// 	return a, nil
// }

type OrganizationArgs struct {
	Id graphql.ID
}

func (r *Root) Organization(ctx context.Context, args OrganizationArgs) (*Organization, error) {
	o, err := r.authenticatedOrganization(ctx, string(args.Id))
	if err != nil {
		return nil, err
	}

	ro := &Organization{
		Organization: *o,
		model:        r.model,
	}

	return ro, nil
}

func (r *Root) Organizations(ctx context.Context) ([]*Organization, error) {
	u, err := cuser(ctx)
	if err != nil {
		return nil, err
	}

	os, err := r.model.UserOrganizations(u.id)
	if err != nil {
		return nil, err
	}

	ros := []*Organization{}

	for _, o := range os {
		ros = append(ros, &Organization{Organization: o, model: r.model})
	}

	return ros, nil
}

type RackLog struct {
	line string
}

type RackLogsArgs struct {
	Oid   graphql.ID
	Rid   graphql.ID
	Since *int32
}

func (r *Root) RackLogs(ctx context.Context, args RackLogsArgs) (chan *RackLog, error) {
	fmt.Printf("ctx: %+v\n", ctx)

	xx := ctx.Value(ContextModel)
	fmt.Printf("xx: %+v\n", xx)

	o, err := r.authenticatedOrganization(ctx, string(args.Oid))
	if err != nil {
		return nil, err
	}

	rr, err := r.model.RackGet(string(args.Rid))
	if err != nil {
		return nil, err
	}

	fmt.Printf("rr: %+v\n", rr)

	if rr.Organization != o.ID {
		return nil, fmt.Errorf("invalid organization")
	}

	ch := make(chan *RackLog)

	fmt.Printf("ch: %+v\n", ch)

	go rackLogs(ctx, &Rack{*rr}, ch)

	return ch, nil
}

func (r *Root) authenticatedOrganization(ctx context.Context, oid string) (*model.Organization, error) {
	o, err := r.model.OrganizationGet(oid)
	if err != nil {
		return nil, err
	}

	u, err := cuser(ctx)
	if err != nil {
		return nil, err
	}

	if !common.SliceContains(o.Users, u.id) {
		return nil, fmt.Errorf("invalid authentication")
	}

	return o, nil
}
