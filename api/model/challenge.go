package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Challenge struct {
	ID string `dynamo:"id"`

	Data []byte `dynamo:"data"`

	ttl int64 `dynamo:"ttl"`
}

func (m *Model) ChallengeGet(id string) (*Challenge, error) {
	var c Challenge

	if err := m.storage.Get("challenges", id, &c); err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}

func (m *Model) ChallengeSave(c *Challenge) error {
	return m.storage.Put("challenges", c)
}

func (c *Challenge) Defaults() {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	if c.ttl == 0 {
		c.ttl = time.Now().UTC().Add(5 * time.Minute).Unix()
	}
}

func (c *Challenge) Validate() []error {
	errs := []error{}

	errs = checkNonzero(errs, c.ID, "id required")
	errs = checkNonzero(errs, c.Data, "data required")

	return errs
}
