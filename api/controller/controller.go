package controller

import "github.com/convox/console/api/model"

type Controller struct {
	model model.Interface
}

func New(m model.Interface) Controller {
	return Controller{model: m}
}
