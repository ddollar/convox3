package resolver

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/convox/console/api/model"
	"github.com/convox/convox/sdk"
	"github.com/graph-gophers/graphql-go"
)

type Rack struct {
	model.Rack
}

func (r *Rack) Id() graphql.ID {
	return graphql.ID(r.Rack.ID)
}

func (r *Rack) Name() string {
	return r.Rack.Name
}

type AppArgs struct {
	Name string
}

func (r *Rack) App(ctx context.Context, args AppArgs) (*App, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	a, err := c.AppGet(args.Name)
	if err != nil {
		return nil, err
	}

	ra := &App{App: *a, rack: r}

	return ra, nil
}

func (r *Rack) Apps(ctx context.Context) ([]*App, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	as, err := c.AppList()
	if err != nil {
		return nil, err
	}

	ras := []*App{}

	for _, a := range as {
		ras = append(ras, &App{App: a, rack: r})
	}

	if len(ras) > 0 {
		for i := 0; i < 20; i++ {
			ras = append(ras, ras[0])
		}
	}

	return ras, nil
}

func (r *Rack) Capacity(ctx context.Context) (*Capacity, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	sc, err := c.CapacityGet()
	if err != nil {
		return nil, err
	}

	cc := &Capacity{
		cpu: CapacityMetric{
			total: int32(sc.ClusterCPU),
			used:  int32(sc.ProcessCPU),
		},
		mem: CapacityMetric{
			total: int32(sc.ClusterMemory),
			used:  int32(sc.ProcessMemory),
		},
	}

	return cc, nil
}

func (r *Rack) Instances(ctx context.Context) ([]*Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	is, err := c.InstanceList()
	if err != nil {
		return nil, err
	}

	ris := []*Instance{}

	for _, i := range is {
		fmt.Printf("i: %+v\n", i)
		ris = append(ris, &Instance{Instance: i})
	}

	fmt.Printf("ris: %+v\n", ris)

	return ris, nil
}

func (r *Rack) Runtime() *graphql.ID {
	if r.Rack.Runtime == "" {
		return nil
	}

	id := graphql.ID(r.Rack.Runtime)

	return &id
}

func (r *Rack) Status(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return "", err
	}

	s, err := c.SystemGet()
	if err != nil {
		return "unknown", nil
	}

	return s.Status, nil
}

func (r *Rack) client(ctx context.Context) (*sdk.Client, error) {
	if r.Rack.Host == "" {
		return nil, fmt.Errorf("no host")
	}

	u := url.URL{
		Host:   r.Rack.Host,
		Scheme: "https",
		User:   url.UserPassword("convox", r.Password),
	}

	s, err := sdk.New(u.String())
	if err != nil {
		return nil, err
	}

	s.Client = s.Client.WithContext(ctx)

	return s, nil
}
