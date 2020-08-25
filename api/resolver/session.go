package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type Session struct {
	model.Session
}

func (s *Session) Id() graphql.ID {
	return graphql.ID(s.Session.ID)
}

func (s *Session) Expires() int32 {
	return int32(s.Session.Expires.UTC().Unix())
}
