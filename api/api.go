package api

import (
	"net/http"

	"github.com/convox/console/api/controller"
	"github.com/convox/console/api/model"
	"github.com/convox/console/api/resolver"
	"github.com/gobuffalo/packr/v2"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

type Api struct {
	box        *packr.Box
	controller controller.Controller
	model      model.Interface
	schema     *graphql.Schema
}

type Query struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func New(m model.Interface, box *packr.Box) (*Api, error) {
	a := &Api{
		box:        box,
		controller: controller.New(m),
		model:      m,
	}

	if err := a.initializeGraphql(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Api) Handler() http.HandlerFunc {
	return graphqlws.NewHandlerFunc(a.schema, &relay.Handler{Schema: a.schema})
}

func (a *Api) initializeGraphql() error {
	r := resolver.New(a.model)

	schema, err := a.box.FindString("schema.graphql")
	if err != nil {
		return err
	}

	s, err := graphql.ParseSchema(schema, r)
	if err != nil {
		return err
	}

	a.schema = s

	return nil
}
