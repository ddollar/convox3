package resolver

import (
	"context"
	"fmt"
	"net/url"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/common"
	"github.com/convox/convox/sdk"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

var jwtHash = jwt.NewHS256([]byte("secret"))

func authenticatedOrganization(ctx context.Context, model model.Interface, oid string) (*model.Organization, error) {
	o, err := model.OrganizationGet(oid)
	if err != nil {
		return nil, err
	}

	u, err := currentUser(ctx)
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

func currentUser(ctx context.Context) (*User, error) {
	token, ok := ctx.Value(graphqlws.ContextAuthorization).(string)
	if !ok {
		return nil, AuthenticationError{fmt.Errorf("no token")}
	}

	var data map[string]string

	if _, err := jwt.Verify([]byte(token), jwtHash, &data); err != nil {
		return nil, err
	}

	u := &User{
		id:    data["id"],
		email: data["email"],
	}

	return u, nil
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
