package resolver

import (
	"context"
	"fmt"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/integration"
	"github.com/graph-gophers/graphql-go"
)

type Runtime struct {
	model.Integration
}

func (r *Runtime) Engines() ([]*Engine, error) {
	es := []*Engine{}

	switch r.Integration.Provider {
	case "aws":
		es = append(es, &Engine{name: "v3", description: "EKS (Version 3)"})
		es = append(es, &Engine{name: "v2", description: "ECS (Version 2)"})
	case "azure":
		es = append(es, &Engine{name: "v3", description: "AKS"})
	case "do":
		es = append(es, &Engine{name: "v3", description: "DOKS"})
	case "gcp":
		es = append(es, &Engine{name: "v3", description: "GKE"})
	}

	return es, nil
}

func (r *Runtime) Id() graphql.ID {
	return graphql.ID(r.Integration.ID)
}

func (r *Runtime) Parameters(ctx context.Context) ([]string, error) {
	ri, err := r.runtime()
	if err != nil {
		return nil, err
	}

	ps, err := ri.ParameterList()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (r *Runtime) Provider() string {
	return r.Integration.Provider
}

func (r *Runtime) Regions() ([]string, error) {
	ri, err := r.runtime()
	if err != nil {
		return nil, err
	}

	rs, err := ri.RegionList()
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (r *Runtime) Title() (string, error) {
	ii, err := r.Integration.Integration()
	if err != nil {
		return "", err
	}

	t, err := ii.Title(r.Integration.Attributes)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (r *Runtime) runtime() (integration.Runtime, error) {
	ii, err := r.Integration.Integration()
	if err != nil {
		return nil, err
	}

	ri, ok := ii.(integration.Runtime)
	if !ok {
		return nil, fmt.Errorf("invalid runtime")
	}

	return ri, nil
}
