package model

import "github.com/convox/console/pkg/integration"

type Integration struct {
	ID string `dynamo:"id" json:"id"`

	AccessToken    string            `dynamo:"access-token,encrypted"`
	Attributes     map[string]string `dynamo:"attributes" json:"attributes"`
	Kind           string            `dynamo:"kind" json:"kind"`
	OrganizationId string            `dynamo:"organization-id" json:"organization-id"`
	Provider       string            `dynamo:"provider" json:"provider"`
}

type Integrations []Integration

func (i *Integration) Integration() (integration.Integration, error) {
	return integration.New(i.ID, i.OrganizationId, i.Provider, i.AccessToken)
}
