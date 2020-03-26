package resolver

import "github.com/convox/convox/pkg/structs"

type App struct {
	structs.App
}

func (a *App) Name() string {
	return a.App.Name
}

func (a *App) Status() string {
	return a.App.Status
}
