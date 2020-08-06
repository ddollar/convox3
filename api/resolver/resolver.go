package resolver

import "github.com/convox/console/api/model"

func New(model model.Interface) *Root {
	return &Root{model: model}
}
