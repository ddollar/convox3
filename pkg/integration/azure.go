package integration

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/subscriptions"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/convox/console/pkg/cache"
	"github.com/convox/console/pkg/helpers"
	"github.com/pkg/errors"
)

type Azure struct {
	id    string
	oid   string
	token string
}

// Integration

func (i *Azure) Name() string {
	return "Microsoft Azure"
}

func (i *Azure) Slug() string {
	return "azure"
}

func (i *Azure) Status() (string, error) {
	if _, err := i.Title(nil); err != nil {
		return "disconnected", nil
	}

	return "connected", nil
}

func (i *Azure) Title(attrs map[string]string) (string, error) {
	if t, ok := cache.Get(i.collection(), "title").(string); ok {
		return t, nil
	}

	ctx := context.Background()

	c, err := i.subscriptionsClient()
	if err != nil {
		return "", errors.WithStack(err)
	}

	creds, err := i.credentials()
	if err != nil {
		return "", errors.WithStack(err)
	}

	s, err := c.Get(ctx, creds["subscription-id"])
	if err != nil {
		return "", errors.WithStack(err)
	}

	title := helpers.DefaultString(s.DisplayName, creds["tenant-id"])

	if err := cache.Set(i.collection(), "title", title, 24*time.Hour); err != nil {
		return "", errors.WithStack(err)
	}

	return title, nil
}

// Runtime

func (i *Azure) Credentials() (map[string]string, error) {
	cs, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	creds := map[string]string{
		"ARM_CLIENT_ID":       cs["client-id"],
		"ARM_CLIENT_SECRET":   cs["client-secret"],
		"ARM_SUBSCRIPTION_ID": cs["subscription-id"],
		"ARM_TENANT_ID":       cs["tenant-id"],
	}

	return creds, nil
}

func (i *Azure) ParameterList() ([]string, error) {
	ps, err := terraformInputs("https://raw.githubusercontent.com/convox/convox/master/terraform/system/azure/variables.tf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ps = helpers.SliceRemove(ps, "name")
	ps = helpers.SliceRemove(ps, "region")

	return ps, nil
}

func (i *Azure) RegionList() ([]string, error) {
	if rs, ok := cache.Get(i.collection(), "regions").([]string); ok {
		return rs, nil
	}

	ctx := context.Background()

	c, err := i.subscriptionsClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	creds, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := c.ListLocations(ctx, creds["subscription-id"])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Value == nil {
		return []string{}, nil
	}

	regions := []string{}

	for _, v := range *res.Value {
		if v.Name != nil {
			regions = append(regions, *v.Name)
		}
	}

	sort.Strings(regions)

	if err := cache.Set(i.collection(), "regions", regions, 24*time.Hour); err != nil {
		return nil, errors.WithStack(err)
	}

	return regions, nil
}

func (i *Azure) authorizer(resource string) (autorest.Authorizer, error) {
	creds, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	oc, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, creds["tenant-id"])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st, err := adal.NewServicePrincipalToken(*oc, creds["client-id"], creds["client-secret"], resource)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	auth := autorest.NewBearerAuthorizer(st)

	return auth, nil
}

func (i *Azure) collection() string {
	return fmt.Sprintf("integration.azure.%x", sha256.Sum256([]byte(i.token)))
}

func (i *Azure) credentials() (map[string]string, error) {
	var params map[string]string

	if err := json.Unmarshal([]byte(i.token), &params); err != nil {
		return nil, errors.WithStack(err)
	}

	return params, nil
}

func (i *Azure) subscriptionsClient() (*subscriptions.Client, error) {
	auth, err := i.authorizer(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := subscriptions.NewClient()

	c.Authorizer = auth

	return &c, nil
}
