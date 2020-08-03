package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type Root struct {
}

type LoginArgs struct {
	Email    string
	Password string
}

func (r *Root) Login(ctx context.Context, args LoginArgs) (*Authentication, error) {
	mu, err := cmodel(ctx).UserAuthenticatePassword(args.Email, args.Password)
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
	o, err := corg(ctx, string(args.Id))
	if err != nil {
		return nil, err
	}

	ro := &Organization{*o}

	return ro, nil
}

func (r *Root) Organizations(ctx context.Context) ([]*Organization, error) {
	u, err := cuser(ctx)
	if err != nil {
		return nil, err
	}

	os, err := cmodel(ctx).UserOrganizations(u.id)
	if err != nil {
		return nil, err
	}

	ros := []*Organization{}

	for _, o := range os {
		ros = append(ros, &Organization{o})
	}

	return ros, nil
}

type RacksArgs struct {
	Oid graphql.ID
}
