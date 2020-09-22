package resolver

import (
	"context"
	"time"

	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Release struct {
	structs.Release
	app *App
}

func (r *Release) Id() graphql.ID {
	return graphql.ID(r.Release.Id)
}

func (r *Release) Build(ctx context.Context) (*Build, error) {
	if r.Release.Build == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := r.app.rack.client(ctx)
	if err != nil {
		return nil, err
	}

	b, err := c.BuildGet(r.app.App.Name, r.Release.Build)
	if err != nil {
		return nil, err
	}

	return &Build{*b}, nil
}

func (r *Release) Created() int32 {
	return int32(r.Release.Created.UTC().Unix())
}

func (r *Release) Description() string {
	return r.Release.Description
}

func (r *Release) Env() string {
	return r.Release.Env
}

func (r *Release) Manifest() string {
	return r.Release.Manifest
}
