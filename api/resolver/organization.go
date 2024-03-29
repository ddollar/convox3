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
	model model.Interface
}

func (o *Organization) Id() graphql.ID {
	return graphql.ID(o.Organization.ID)
}

func (o *Organization) Access(ctx context.Context) (string, error) {
	uid, err := currentUid(ctx)
	if err != nil {
		return "", err
	}

	switch {
	case common.SliceContains(o.Organization.Administrators, uid):
		return "administrator", nil
	case common.SliceContains(o.Organization.Operators, uid):
		return "operator", nil
	case common.SliceContains(o.Organization.Users, uid):
		return "developer", nil
	default:
		return "", fmt.Errorf("no access")
	}
}

type IntegrationsArgs struct {
	Kind     string
	Provider *string
}

func (o *Organization) Integrations(ctx context.Context, args IntegrationsArgs) ([]*Integration, error) {
	is, err := o.model.OrganizationIntegrations(o.Organization.ID)
	if err != nil {
		return nil, err
	}

	ris := []*Integration{}

	for _, i := range is {
		if i.Kind == args.Kind {
			if args.Provider == nil || *args.Provider == i.Provider {
				ris = append(ris, &Integration{i})
			}
		}
	}

	return ris, nil
}

func (o *Organization) Members(ctx context.Context) ([]*Member, error) {
	us, err := o.model.UserGetBatch(o.Organization.Users)
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
	rs, err := o.model.OrganizationRacks(o.Organization.ID)
	if err != nil {
		return nil, err
	}

	rrs := []*Rack{}

	for _, r := range rs {
		rrs = append(rrs, &Rack{Rack: r, model: o.model})
	}

	return rrs, nil
}

type RackArgs struct {
	Id graphql.ID
}

func (o *Organization) Rack(ctx context.Context, args RackArgs) (*Rack, error) {
	r, err := authenticatedRack(ctx, o.model, o.Organization.ID, string(args.Id))
	if err != nil {
		return nil, err
	}

	return &Rack{Rack: *r, model: o.model}, nil
}

type RuntimeArgs struct {
	Id graphql.ID
}

func (o *Organization) Runtime(ctx context.Context, args RuntimeArgs) (*Runtime, error) {
	i, err := authenticatedIntegration(ctx, o.model, o.Organization.ID, string(args.Id))
	if err != nil {
		return nil, err
	}

	r := &Runtime{Integration: *i}

	return r, nil
}
