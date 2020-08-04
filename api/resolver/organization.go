package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/common"
	"github.com/graph-gophers/graphql-go"
)

type Organization struct {
	model.Organization
}

func (o *Organization) Id() graphql.ID {
	return graphql.ID(o.Organization.ID)
}

func (o *Organization) Access(ctx context.Context) (string, error) {
	u, err := cuser(ctx)
	if err != nil {
		return "", err
	}

	switch {
	case common.SliceContains(o.Organization.Administrators, u.id):
		return "administrator", nil
	case common.SliceContains(o.Organization.Operators, u.id):
		return "operator", nil
	case common.SliceContains(o.Organization.Users, u.id):
		return "developer", nil
	default:
		return "", fmt.Errorf("no access")
	}
}

type IntegrationsArgs struct {
	Kind string
}

func (o *Organization) Integrations(ctx context.Context, args IntegrationsArgs) ([]*Integration, error) {
	is, err := cmodel(ctx).OrganizationIntegrations(o.Organization.ID)
	if err != nil {
		return nil, err
	}

	ris := []*Integration{}

	for _, i := range is {
		fmt.Printf("i.Kind: %+v\n", i.Kind)
		if i.Kind == args.Kind {
			ris = append(ris, &Integration{i})
		}
	}

	return ris, nil
}

func (o *Organization) Members(ctx context.Context) ([]*Member, error) {
	us, err := cmodel(ctx).UserGetBatch(o.Organization.Users)
	if err != nil {
		return nil, err
	}

	rms := []*Member{}

	for _, u := range us {
		access := "developer"

		switch {
		case common.SliceContains(o.Organization.Administrators, u.ID):
			access = "administrator"
		case common.SliceContains(o.Organization.Operators, u.ID):
			access = "operator"
		}

		rms = append(rms, &Member{access: access, user: u})
	}

	return rms, nil
}

func (o *Organization) Locked() bool {
	return o.Organization.Locked
}

func (o *Organization) Name() string {
	return o.Organization.Name
}

func (o *Organization) Racks(ctx context.Context) ([]*Rack, error) {
	rs, err := cmodel(ctx).OrganizationRacks(o.Organization.ID)
	if err != nil {
		return nil, err
	}

	rrs := []*Rack{}

	for _, r := range rs {
		rrs = append(rrs, &Rack{r})
	}

	return rrs, nil
}

type RackArgs struct {
	Id graphql.ID
}

func (o *Organization) Rack(ctx context.Context, args RackArgs) (*Rack, error) {
	r, err := cmodel(ctx).RackGet(string(args.Id))
	if err != nil {
		return nil, err
	}

	if r.Organization != o.Organization.ID {
		return nil, fmt.Errorf("invalid organization")
	}

	return &Rack{*r}, nil
}
