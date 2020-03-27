package api

import (
	"net/http"

	"github.com/convox/console/api/model"
	"github.com/convox/console/api/resolver"
	"github.com/gobuffalo/packr/v2"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
)

type Api struct {
	model  model.Interface
	schema *graphql.Schema
}

type Query struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func New(m model.Interface) (*Api, error) {
	a := &Api{
		model: m,
	}

	if err := a.initializeGraphql(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Api) Handler() http.HandlerFunc {
	return graphqlws.NewHandlerFunc(a.schema, &Handler{api: a})
}

// func (a *Api) Route(s *stdapi.Server) error {
// 	s.Route("GET", "/graphql", a.QueryGet)
// 	s.Route("POST", "/graphql", a.QueryPost)

// 	return nil
// }

func (a *Api) initializeGraphql() error {
	r := resolver.New()

	box := packr.New("graphql", "./graphql")

	schema, err := box.FindString("schema.graphql")
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
