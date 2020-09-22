package resolver

import (
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Release struct {
	structs.Release
}

func (r *Release) Id() graphql.ID {
	return graphql.ID(r.Release.Id)
}

func (r *Release) Build() *string {
	if r.Release.Build == "" {
		return nil
	}

	s := r.Release.Build

	return &s
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
