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

func (a *App) Builds(ctx context.Context) ([]*Build, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := a.rack.client(ctx)
	if err != nil {
		return nil, err
	}

	bs, err := c.BuildList(a.App.Name, structs.BuildListOptions{})
	if err != nil {
		return nil, err
	}

	rbs := []*Build{}

	for _, b := range bs {
		rbs = append(rbs, &Build{b})
	}

	return rbs, nil
}

func (a *App) Name() string {
	return a.App.Name
}

func (a *App) Processes(ctx context.Context) ([]*Process, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := a.rack.client(ctx)
	if err != nil {
		return nil, err
	}

	ps, err := c.ProcessList(a.App.Name, structs.ProcessListOptions{})
	if err := timeoutError(err); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	rps := []*Process{}

	for _, p := range ps {
		rps = append(rps, &Process{p})
	}

	return rps, nil
}

func (a *App) Releases(ctx context.Context) ([]*Release, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := a.rack.client(ctx)
	if err != nil {
		return nil, err
	}

	rs, err := c.ReleaseList(a.App.Name, structs.ReleaseListOptions{})
	if err != nil {
		return nil, err
	}

	rrs := []*Release{}

	for _, r := range rs {
		rrs = append(rrs, &Release{r})
	}

	return rrs, nil
}

func (a *App) Services(ctx context.Context) ([]*Service, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := a.rack.client(ctx)
	if err != nil {
		return nil, err
	}

	ss, err := c.ServiceList(a.App.Name)
	if err := timeoutError(err); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	rss := []*Service{}

	for _, s := range ss {
		rss = append(rss, &Service{s})
	}

	return rss, nil
}

func (a *App) Status() string {
	return a.App.Status
}
