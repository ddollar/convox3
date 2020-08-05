package resolver

import (
	"context"
	"time"

	"github.com/convox/convox/pkg/structs"
)

type App struct {
	structs.App
	rack *Rack
}

func (a *App) Name() string {
	return a.App.Name
}

func (a *App) Processes(ctx context.Context) ([]*Process, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := a.rack.client(ctx)
	if err != nil {
		return nil, err
	}

	ps, err := c.ProcessList(a.App.Name, structs.ProcessListOptions{})
	if err != nil {
		return nil, err
	}

	rps := []*Process{}

	for _, p := range ps {
		rps = append(rps, &Process{p})
	}

	return rps, nil
}

func (a *App) Status() string {
	return a.App.Status
}
