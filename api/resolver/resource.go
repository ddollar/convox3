package resolver

import (
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Resource struct {
	structs.Resource
}

func (p *Resource) Name() graphql.ID {
	return graphql.ID(p.Resource.Name)
}

func (p *Resource) Status() string {
	return p.Resource.Status
}

func (p *Resource) Type() string {
	return p.Resource.Type
}
