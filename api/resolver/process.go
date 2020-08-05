package resolver

import (
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Process struct {
	structs.Process
}

func (p *Process) Cpu() int32 {
	return int32(p.Process.Cpu)
}

func (p *Process) Id() graphql.ID {
	return graphql.ID(p.Process.Id)
}

func (p *Process) Mem() int32 {
	return int32(p.Process.Memory)
}

func (p *Process) Service() string {
	return p.Process.Name
}
