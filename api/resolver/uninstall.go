package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type Uninstall struct {
	*model.Uninstall
}

func (u *Uninstall) Id() graphql.ID {
	return graphql.ID(u.ID)
}
