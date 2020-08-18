package model

import "github.com/pkg/errors"

type Webhook struct {
	ID string `dynamo:"id"`

	Branch  string `dynamo:"branch"`
	Event   string `dynamo:"event"`
	Project int64  `dynamo:"project-id"`
	Remote  int64  `dynamo:"remote-id"`

	IntegrationID string `dynamo:"integration-id"`
}

type Webhooks []Webhook

func (m *Model) WebhookGet(id string) (*Webhook, error) {
	w := &Webhook{}

	if err := m.storage.Get("webhooks", id, w); err != nil {
		return nil, errors.WithStack(err)
	}

	return w, nil
}
