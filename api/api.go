package api

import (
	"context"
	"encoding/json"

	"github.com/convox/console/api/model"
	"github.com/convox/console/api/resolver"
	"github.com/convox/stdapi"
	"github.com/gobuffalo/packr/v2"
	"github.com/graph-gophers/graphql-go"
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

func (a *Api) QueryGet(c *stdapi.Context) error {
	q := Query{
		Query: c.Query("query"),
	}

	if vars := c.Query("variables"); vars != "" {
		if err := json.Unmarshal([]byte(vars), &q.Variables); err != nil {
			return err
		}
	}

	return a.query(c, q)
}

func (a *Api) QueryPost(c *stdapi.Context) error {
	var q Query

	if err := c.BodyJSON(&q); err != nil {
		return err
	}

	return a.query(c, q)
}

func (a *Api) Route(s *stdapi.Server) error {
	s.Route("GET", "/graphql", a.QueryGet)
	s.Route("POST", "/graphql", a.QueryPost)

	return nil
}

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

func (a *Api) query(c *stdapi.Context, q Query) error {
	ctx := c.Context()
	ctx = context.WithValue(ctx, "model", a.model)
	ctx = context.WithValue(ctx, "uid", "f8abd4df-f8b4-4cb9-9514-6395e7907f2b")

	res := a.schema.Exec(ctx, q.Query, "", q.Variables)

	if len(res.Errors) > 0 {
		c.Response().WriteHeader(403)
		// fmt.Printf("res.Errors: %+v\n", res.Errors)
		// return nil, fmt.Errorf("server error")
	}

	return c.RenderJSON(res)
}
