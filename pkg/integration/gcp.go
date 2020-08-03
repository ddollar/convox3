package integration

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/convox/console/pkg/cache"
	"github.com/convox/console/pkg/helpers"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type Google struct {
	id    string
	oid   string
	token string
}

// Integration

func (i *Google) Name() string {
	return "Google Cloud"
}

func (i *Google) Slug() string {
	return "gcp"
}

func (i *Google) Status() (string, error) {
	if _, err := i.Title(nil); err != nil {
		return "disconnected", nil
	}

	return "connected", nil
}

func (i *Google) Title(attrs map[string]string) (string, error) {
	creds, err := i.credentials()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return creds["project"], nil
}

// Runtime

func (i *Google) Credentials() (map[string]string, error) {
	cs, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	creds := map[string]string{
		"GOOGLE_CREDENTIALS": cs["credentials"],
		"GOOGLE_PROJECT":     cs["project"],
	}

	return creds, nil
}

func (i *Google) ParameterList() ([]string, error) {
	ps, err := terraformInputs("https://raw.githubusercontent.com/convox/convox/master/terraform/system/gcp/variables.tf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ps = helpers.SliceRemove(ps, "name")
	ps = helpers.SliceRemove(ps, "region")

	return ps, nil
}

func (i *Google) RegionList() ([]string, error) {
	if rs, ok := cache.Get(i.collection(), "regions").([]string); ok {
		return rs, nil
	}

	ctx := context.Background()

	creds, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c, err := i.computeClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	regions := []string{}

	req := c.Regions.List(creds["project"])

	if err := req.Pages(ctx, func(page *compute.RegionList) error {
		for _, region := range page.Items {
			regions = append(regions, region.Description)
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	sort.Strings(regions)

	if err := cache.Set(i.collection(), "regions", regions, 24*time.Hour); err != nil {
		return nil, errors.WithStack(err)
	}

	return regions, nil
}

func (i *Google) credentials() (map[string]string, error) {
	var params map[string]string

	if err := json.Unmarshal([]byte(i.token), &params); err != nil {
		return nil, errors.WithStack(err)
	}

	return params, nil
}

func (i *Google) computeClient() (*compute.Service, error) {
	ctx := context.Background()

	creds, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c, err := google.CredentialsFromJSON(ctx, []byte(creds["credentials"]), "https://www.googleapis.com/auth/compute")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cc, err := compute.NewService(ctx, option.WithTokenSource(c.TokenSource))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cc, nil
}

func (i *Google) collection() string {
	return fmt.Sprintf("integration.google.%x", sha256.Sum256([]byte(i.token)))
}
