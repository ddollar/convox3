package model

import (
	"github.com/convox/console/pkg/integration"
	"github.com/pkg/errors"
)

type Integration struct {
	ID string `dynamo:"id" json:"id"`

	AccessToken    string            `dynamo:"access-token,encrypted"`
	Attributes     map[string]string `dynamo:"attributes" json:"attributes"`
	Kind           string            `dynamo:"kind" json:"kind"`
	OrganizationId string            `dynamo:"organization-id" json:"organization-id"`
	Provider       string            `dynamo:"provider" json:"provider"`
}

type Integrations []Integration

func (m *Model) IntegrationGet(iid string) (*Integration, error) {
	i := &Integration{}

	if err := m.storage.Get("integrations", iid, i); err != nil {
		return nil, errors.WithStack(err)
	}

	return i, nil
}

func (i *Integration) Integration() (integration.Integration, error) {
	return integration.New(i.ID, i.OrganizationId, i.Provider, i.AccessToken)
}
