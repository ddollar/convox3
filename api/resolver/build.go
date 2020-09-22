package resolver

import (
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Build struct {
	structs.Build
}

func (b *Build) Id() graphql.ID {
	return graphql.ID(b.Build.Id)
}

func (b *Build) Description() string {
	return b.Build.Description
}

func (b *Build) Ended() *int32 {
	if b.Build.Ended.IsZero() {
		return nil
	}

	t := int32(b.Build.Ended.UTC().Unix())

	return &t
}

func (b *Build) Logs() string {
	return b.Build.Logs
}

func (b *Build) Manifest() string {
	return b.Build.Manifest
}

func (b *Build) Release() *string {
	if b.Build.Release == "" {
		return nil
	}

	s := b.Build.Release

	return &s
}

func (b *Build) Started() int32 {
	return int32(b.Build.Started.UTC().Unix())
}

func (b *Build) Status() string {
	return b.Build.Status
}
