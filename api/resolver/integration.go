package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type Integration struct {
	model.Integration
}

func (i *Integration) Id() graphql.ID {
	return graphql.ID(i.Integration.ID)
}

func (i *Integration) Provider() string {
	return i.Integration.Provider
}

func (i *Integration) Title() (string, error) {
	ii, err := i.Integration.Integration()
	if err != nil {
		return "", err
	}

	t, err := ii.Title(i.Integration.Attributes)
	if err != nil {
		return "", err
	}

	return t, nil
}
