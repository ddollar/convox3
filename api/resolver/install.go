package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type Install struct {
	*model.Install
}

func (i *Install) Id() graphql.ID {
	return graphql.ID(i.ID)
}
