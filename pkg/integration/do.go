package integration

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/convox/console/pkg/helpers"
	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type DigitalOcean struct {
	id    string
	oid   string
	token string
}

// Integration

func (i *DigitalOcean) Name() string {
	return "Digital Ocean"
}

func (i *DigitalOcean) Slug() string {
	return "do"
}

func (i *DigitalOcean) Status() (string, error) {
	c, err := i.client()
	if err != nil {
		return "disconnected", nil
	}

	if _, _, err := c.Account.Get(context.Background()); err != nil {
		return "disconnected", nil
	}

	return "connected", nil
}

func (i *DigitalOcean) Title(attrs map[string]string) (string, error) {
	c, err := i.client()
	if err != nil {
		return "", errors.WithStack(err)
	}

	a, _, err := c.Account.Get(context.Background())
	if err != nil {
		return "", errors.WithStack(err)
	}

	return a.Email, nil
}

// Runtime

func (i *DigitalOcean) Credentials() (map[string]string, error) {
	cs, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	creds := map[string]string{
		"DIGITALOCEAN_ACCESS_ID":  cs["access"],
		"DIGITALOCEAN_SECRET_KEY": cs["secret"],
		"DIGITALOCEAN_TOKEN":      cs["token"],
	}

	return creds, nil
}

func (i *DigitalOcean) ParameterList() ([]string, error) {
	ps, err := terraformInputs("https://raw.githubusercontent.com/convox/convox/master/terraform/system/do/variables.tf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ps = helpers.SliceRemove(ps, "access_id")
	ps = helpers.SliceRemove(ps, "name")
	ps = helpers.SliceRemove(ps, "region")
	ps = helpers.SliceRemove(ps, "secret_key")
	ps = helpers.SliceRemove(ps, "token")

	return ps, nil
}

func (i *DigitalOcean) RegionList() ([]string, error) {
	// c, err := i.client()
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// opts, _, err := c.Kubernetes.GetOptions(context.Background())
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// kuberegions := map[string]bool{}

	// for _, region := range opts.Regions {
	// 	kuberegions[region.Slug] = true
	// }

	// regions, _, err := c.Regions.List(context.Background(), &godo.ListOptions{})
	// if err != nil {
	// 	return nil, err
	// }

	// spacesregions := map[string]bool{}

	// for _, region := range regions {
	// 	fmt.Printf("region.Features: %s %+v\n", region.Slug, region.Features)
	// 	if helpers.SliceContains(region.Features, "storage") {
	// 		spacesregions[region.Slug] = true
	// 	}
	// }

	// fmt.Printf("kuberegions: %+v\n", kuberegions)
	// fmt.Printf("spacesregions: %+v\n", spacesregions)

	// rs := []string{}

	rs := []string{"nyc3", "ams3", "sfo2", "sgp1", "fra1"}

	sort.Strings(rs)

	return rs, nil
}

func (i *DigitalOcean) client() (*godo.Client, error) {
	creds, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	token := &oauth2.Token{
		AccessToken: creds["token"],
	}

	c := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	dc := godo.NewClient(c)

	return dc, nil
}

func (i *DigitalOcean) credentials() (map[string]string, error) {
	var params map[string]string

	if err := json.Unmarshal([]byte(i.token), &params); err != nil {
		return nil, errors.WithStack(err)
	}

	return params, nil
}
