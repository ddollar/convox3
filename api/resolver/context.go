package resolver

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

type ContextKey int

const (
	ContextModel ContextKey = iota
	ContextToken
)

// func cmodel(ctx context.Context) model.Interface {
// 	m, ok := ctx.Value(ContextModel).(model.Interface)
// 	if !ok {
// 		panic("no model available")
// 	}

// 	return m
// }

// func (r *Root) corg(ctx context.Context, oid string) (*model.Organization, error) {
// 	o, err := Model.OrganizationGet(oid)
// 	if err != nil {
// 		return nil, err
// 	}

// 	u, err := cuser(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !common.SliceContains(o.Users, u.id) {
// 		return nil, fmt.Errorf("invalid user")
// 	}

// 	return o, nil
// }

// func cuid(ctx context.Context) (string, error) {
// 	uid, ok := ctx.Value("uid").(string)
// 	if !ok {
// 		return "", fmt.Errorf("invalid uid")
// 	}

// 	return uid, nil
// }

func cuser(ctx context.Context) (*User, error) {
	token, ok := ctx.Value(graphqlws.ContextAuthorization).(string)
	if !ok {
		return nil, AuthenticationError{fmt.Errorf("no token")}
	}

	u, err := UserFromToken(token)
	if err != nil {
		return nil, AuthenticationError{err}
	}

	return u, nil
}
