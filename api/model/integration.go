package model

import (
	"fmt"

	"github.com/convox/console/pkg/integration"
	"github.com/pkg/errors"
)

type Integration struct {
	ID string `dynamo:"id" json:"id"`

	AccessToken    string            `dynamo:"access-token,encrypted"`
	Attributes     map[string]string `dynamo:"attributes" json:"attributes"`
	Kind           string            `dynamo:"kind" json:"kind"`
	OrganizationID string            `dynamo:"organization-id" json:"organization-id"`
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
	return integration.New(i.ID, i.OrganizationID, i.Provider, i.AccessToken)
}

func (i *Integration) Notification() (integration.Notification, error) {
	ii, err := i.Integration()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ir, ok := ii.(integration.Notification)
	if !ok {
		return nil, errors.WithStack(fmt.Errorf("not a notification integration"))
	}

	return ir, nil
}

func (i *Integration) Source() (integration.Source, error) {
	ii, err := i.Integration()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	is, ok := ii.(integration.Source)
	if !ok {
		return nil, errors.WithStack(fmt.Errorf("not a source integration"))
	}

	return is, nil
}

func (is Integrations) ByKind(kind string) Integrations {
	isk := Integrations{}

	for _, i := range is {
		if i.Kind == kind {
			isk = append(isk, i)
		}
	}

	return isk
}
