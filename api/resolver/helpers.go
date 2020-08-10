package resolver

import (
	"context"
	"fmt"
	"net/url"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/common"
	"github.com/convox/convox/sdk"
)

func authenticatedOrganization(ctx context.Context, model model.Interface, oid string) (*model.Organization, error) {
	o, err := model.OrganizationGet(oid)
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

func authenticatedRack(ctx context.Context, model model.Interface, oid, rid string) (*model.Rack, error) {
	o, err := authenticatedOrganization(ctx, model, oid)
	if err != nil {
		return nil, err
	}

	r, err := model.RackGet(rid)
	if err != nil {
		return nil, err
	}

	if r.Organization != o.ID {
		return nil, fmt.Errorf("invalid organization")
	}

	return r, nil
}

func rackClient(ctx context.Context, host, password string) (*sdk.Client, error) {
	if host == "" {
		return nil, fmt.Errorf("no host")
	}

	u := url.URL{
		Host:   host,
		Scheme: "https",
		User:   url.UserPassword("convox", password),
	}

	s, err := sdk.New(u.String())
	if err != nil {
		return nil, err
	}

	s.Client = s.Client.WithContext(ctx)

	return s, nil
}
