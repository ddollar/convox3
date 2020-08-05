package resolver

import (
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Service struct {
	structs.Service
}

func (p *Service) Count() int32 {
	return int32(p.Service.Count)
}

func (p *Service) Cpu() int32 {
	return int32(p.Service.Cpu)
}

func (p *Service) Domain() string {
	return p.Service.Domain
}

func (p *Service) Mem() int32 {
	return int32(p.Service.Memory)
}

func (p *Service) Name() graphql.ID {
	return graphql.ID(p.Service.Name)
}
