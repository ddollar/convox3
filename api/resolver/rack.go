package resolver

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/common"
	"github.com/convox/convox/pkg/options"
	"github.com/convox/convox/pkg/structs"
	"github.com/convox/convox/sdk"
	"github.com/graph-gophers/graphql-go"
)

type Rack struct {
	model.Rack
	model model.Interface
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
	if err := timeoutError(err); err != nil {
		return nil, err
	}
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
	if err := timeoutError(err); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	ras := []*App{}

	for _, a := range as {
		ras = append(ras, &App{App: a, rack: r})
	}

	// if len(ras) > 0 {
	// 	for i := 0; i < 20; i++ {
	// 		ras = append(ras, ras[0])
	// 	}
	// }

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
	if err := timeoutError(err); err != nil {
		return nil, err
	}
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

func (r *Rack) Install(ctx context.Context) (*Install, error) {
	if r.Rack.Install == "" {
		return nil, nil
	}

	i, err := authenticatedInstall(ctx, r.model, r.Organization, r.Rack.Install)
	if err != nil {
		return nil, err
	}

	ii := &Install{i}

	return ii, nil
}

func (r *Rack) Instances(ctx context.Context) ([]*Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	is, err := c.InstanceList()
	if err := timeoutError(err); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	ris := []*Instance{}

	for _, i := range is {
		ris = append(ris, &Instance{Instance: i})
	}

	return ris, nil
}

func (r *Rack) Runtime() *graphql.ID {
	if r.Rack.Runtime == "" {
		return nil
	}

	id := graphql.ID(r.Rack.Runtime)

	return &id
}

func (r *Rack) Processes(ctx context.Context) ([]*Process, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	ps, err := c.SystemProcesses(structs.SystemProcessesOptions{All: options.Bool(true)})
	if err := timeoutError(err); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	sort.Slice(ps, func(i, j int) bool {
		if ps[i].App != ps[j].App {
			return ps[i].App < ps[j].App
		}
		if ps[i].Name != ps[j].Name {
			return ps[i].Name < ps[j].Name
		}
		return ps[i].Id < ps[j].Id
	})

	rps := []*Process{}

	for _, p := range ps {
		rps = append(rps, &Process{p})
	}

	return rps, nil
}

func (r *Rack) Provider(ctx context.Context) string {
	return common.CoalesceString(r.Rack.Provider, "aws")
}

func (r *Rack) Resources(ctx context.Context) ([]*Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	c, err := r.client(ctx)
	if err != nil {
		return nil, err
	}

	rs, err := c.SystemResourceList()
	if err := timeoutError(err); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	rrs := []*Resource{}

	for _, r := range rs {
		rrs = append(rrs, &Resource{r})
	}

	return rrs, nil
}

func (r *Rack) Status(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	is, err := r.statusInstalling()
	if err != nil {
		return "", err
	}
	if is != "" {
		return is, nil
	}

	us, err := r.statusUninstalling()
	if err != nil {
		return "", err
	}
	if us != "" {
		return us, nil
	}

	c, err := r.client(ctx)
	if err != nil {
		return "", err
	}

	s, err := c.SystemGet()
	if err := timeoutError(err); err != nil {
		return "", err
	}
	if err != nil {
		return "unknown", nil
	}

	return s.Status, nil
}

func (r *Rack) Uninstall(ctx context.Context) (*Uninstall, error) {
	if r.Rack.Uninstall == "" {
		return nil, nil
	}

	u, err := authenticatedUninstall(ctx, r.model, r.Organization, r.Rack.Uninstall)
	if err != nil {
		return nil, err
	}

	uu := &Uninstall{u}

	return uu, nil
}

func (r *Rack) Uninstallable(ctx context.Context) (bool, error) {
	return rackUninstallable(r.model, r.Rack.ID)
}

func (r *Rack) Updates(ctx context.Context) ([]*Update, error) {
	return nil, nil
}

func (r *Rack) UpdateDay(ctx context.Context) int32 {
	return int32(r.Rack.UpdateDay)
}

func (r *Rack) UpdateFrequency(ctx context.Context) string {
	return common.CoalesceString(r.Rack.UpdateFrequency, "never")
}

func (r *Rack) UpdateHour(ctx context.Context) int32 {
	return int32(r.Rack.UpdateHour)
}

func (r *Rack) client(ctx context.Context) (*sdk.Client, error) {
	return rackClient(ctx, r.Host, r.Password)
}

func (r *Rack) statusInstalling() (string, error) {
	if r.Rack.Install == "" {
		return "", nil
	}

	i, err := r.model.InstallGet(r.Rack.Install)
	fmt.Printf("r.Rack.Name: %+v\n", r.Rack.Name)
	fmt.Printf("i: %+v\n", i)
	if err != nil {
		return "", err
	}

	switch i.Status {
	case "pending", "running", "starting":
		return "installing", nil
	case "failed":
		return "failed", nil
	default:
		return "", nil
	}
}

func (r *Rack) statusUninstalling() (string, error) {
	if r.Rack.Uninstall == "" {
		return "", nil
	}

	u, err := r.model.UninstallGet(r.Rack.Uninstall)
	fmt.Printf("r.Rack.Name: %+v\n", r.Rack.Name)
	fmt.Printf("u: %+v\n", u)
	if err != nil {
		return "", err
	}

	switch u.Status {
	case "pending", "running", "starting":
		return "uninstalling", nil
	case "failed":
		return "failed", nil
	default:
		return "", nil
	}
}
