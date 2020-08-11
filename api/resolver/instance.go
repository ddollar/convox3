package resolver

import (
	"github.com/convox/convox/pkg/options"
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Instance struct {
	structs.Instance
}

func (i *Instance) Id() graphql.ID {
	return graphql.ID(i.Instance.Id)
}

func (i *Instance) Cpu() float64 {
	return i.Instance.Cpu
}

func (i *Instance) Mem() float64 {
	return i.Instance.Memory
}

func (i *Instance) Private() string {
	return i.Instance.PrivateIp
}

func (i *Instance) Public() *string {
	if i.Instance.PublicIp != "" {
		return options.String(i.Instance.PublicIp)
	}

	return nil
}

func (i *Instance) Started() int32 {
	return int32(i.Instance.Started.Unix())
}

func (i *Instance) Status() string {
	switch i.Instance.Status {
	case "active", "draining":
		return i.Instance.Status
	default:
		return "unknown"
	}
}
