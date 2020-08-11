package resolver

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
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

	fmt.Printf("args: %+v\n", args)

	if err := c.ProcessStop(args.App, string(args.Pid)); err != nil {
		return "", err
	}

	return string(args.Pid), nil
}
