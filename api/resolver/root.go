package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/pkg/common"
	"github.com/graph-gophers/graphql-go"
)

type Root struct {
}

type OrganizationArgs struct {
	Id graphql.ID
}

func (r *Root) Organization(ctx context.Context, args OrganizationArgs) (*Organization, error) {
	m, err := cmodel(ctx)
	if err != nil {
		return nil, err
	}

	o, err := m.OrganizationGet(string(args.Id))
	if err != nil {
		return nil, err
	}

	uid, err := cuid(ctx)
	if err != nil {
		return nil, err
	}

	if !common.SliceContains(o.Users, uid) {
		return nil, fmt.Errorf("invalid user")
	}

	return &Organization{*o}, nil
}

func (r *Root) Organizations(ctx context.Context) ([]*Organization, error) {
	m, err := cmodel(ctx)
	if err != nil {
		return nil, err
	}

	uid, err := cuid(ctx)
	if err != nil {
		return nil, err
	}

	os, err := m.UserOrganizations(uid)
	if err != nil {
		return nil, err
	}

	ros := []*Organization{}

	for _, o := range os {
		ros = append(ros, &Organization{o})
	}

	return ros, nil
}
