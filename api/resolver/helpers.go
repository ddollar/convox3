package resolver

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/common"
	"github.com/convox/convox/sdk"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

var jwtHash = jwt.NewHS256([]byte("secret"))

func authenticatedInstall(ctx context.Context, model model.Interface, oid, iid string) (*model.Install, error) {
	o, err := authenticatedOrganization(ctx, model, oid)
	if err != nil {
		return nil, err
	}

	i, err := model.InstallGet(iid)
	if err != nil {
		return nil, err
	}

	if i.OrganizationID != o.ID {
		return nil, fmt.Errorf("invalid install")
	}

	return i, nil
}

func authenticatedIntegration(ctx context.Context, model model.Interface, oid, iid string) (*model.Integration, error) {
	o, err := authenticatedOrganization(ctx, model, oid)
	if err != nil {
		return nil, err
	}

	i, err := model.IntegrationGet(iid)
	if err != nil {
		return nil, err
	}

	if i.OrganizationID != o.ID {
		return nil, fmt.Errorf("invalid integration")
	}

	return i, nil
}

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
		return nil, fmt.Errorf("invalid rack")
	}

	return r, nil
}

func authenticatedUninstall(ctx context.Context, model model.Interface, oid, uid string) (*model.Uninstall, error) {
	o, err := authenticatedOrganization(ctx, model, oid)
	if err != nil {
		return nil, err
	}

	u, err := model.UninstallGet(uid)
	if err != nil {
		return nil, err
	}

	if u.OrganizationID != o.ID {
		return nil, fmt.Errorf("invalid install")
	}

	return u, nil
}

func checkNonzero(errs []error, value interface{}, message string) []error {
	if reflect.ValueOf(value).IsZero() {
		errs = append(errs, errors.New(message))
	}

	return errs
}

func collateErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	es := []string{}

	for _, err := range errs {
		es = append(es, err.Error())
	}

	return errors.New(strings.Join(es, ", "))
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

func rackUninstallable(m model.Interface, rid string) (bool, error) {
	data, err := m.RackStateLoad(rid)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return false, nil
	}

	return true, nil
}

func timeoutError(err error) error {
	if err == context.DeadlineExceeded {
		return TimeoutError{err}
	}

	switch t := err.(type) {
	case *url.Error:
		if t.Err == context.DeadlineExceeded {
			return TimeoutError{err}
		}
	}

	return nil
}
