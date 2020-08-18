package resolver

import "context"

type Engine struct {
	description string
	name        string
}

func (e *Engine) Description(ctx context.Context) string {
	return e.description
}

func (e *Engine) Name(ctx context.Context) string {
	return e.name
}
