package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type Token struct {
	model.Token
}

func (t *Token) Id() graphql.ID {
	return graphql.ID(t.Token.ID)
}

func (t *Token) Name() string {
	return t.Token.Name
}

func (t *Token) Used() int32 {
	return int32(t.Token.Used.Unix())
}

type TokenRequest struct {
	data string
	id   string
}

func (t *TokenRequest) Data() string {
	return t.data
}

func (t *TokenRequest) Id() string {
	return t.id
}
