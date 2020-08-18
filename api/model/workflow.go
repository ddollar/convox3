package model

import (
	"fmt"

	"github.com/convox/console/pkg/integration"
	"github.com/pkg/errors"
)

type Workflow struct {
	ID string `dynamo:"id" json:"id"`

	Hook   int64             `dynamo:"hook"`
	Kind   string            `dynamo:"kind"`
	Name   string            `dynamo:"name"`
	Params map[string]string `dynamo:"params"`
	Tasks  Tasks             `dynamo:"tasks"`

	// legacy
	TriggerKind string `dynamo:"trigger-kind"`
	TriggerID   string `dynamo:"trigger-id"`

	IntegrationID  string `dynamo:"integration-id"`
	OrganizationId string `dynamo:"organization-id"`

	repository int64 `dynamo:"repository"`
}

type Workflows []Workflow

func (m *Model) WorkflowGet(id string) (*Workflow, error) {
	w := &Workflow{}

	if err := m.storage.Get("workflows", id, w); err != nil {
		return nil, errors.WithStack(err)
	}

	return w, nil
}

func (m *Model) WorkflowIntegration(id string) (*Integration, error) {
	w, err := m.WorkflowGet(id)
	if err != nil {
		return nil, err
	}

	if w.IntegrationID == "" {
		return nil, nil
	}

	return m.IntegrationGet(w.IntegrationID)
}

func (m *Model) WorkflowRepository(id string) (int64, error) {
	w, err := m.WorkflowGet(id)
	if err != nil {
		return 0, err
	}

	if w.TriggerID == "" {
		return 0, errors.WithStack(fmt.Errorf("no repository"))
	}

	wh, err := m.WebhookGet(w.TriggerID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return wh.Project, nil
}

func (m *Model) WorkflowSource(id string) (integration.Source, error) {
	w, err := m.WorkflowGet(id)
	if err != nil {
		return nil, err
	}

	if w.IntegrationID == "" && w.TriggerID != "" {
		wh, err := m.WebhookGet(w.TriggerID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		i, err := m.IntegrationGet(wh.IntegrationID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return i.Source()
	}

	i, err := m.WorkflowIntegration(w.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if i == nil {
		return nil, nil
	}

	return i.Source()
}
