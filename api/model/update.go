package model

import "time"

type Update struct {
	ID string `dynamo:"id"`

	OrganizationID string `dynamo:"organization-id"`
	RackID         string `dynamo:"rack-id"`

	Created  time.Time         `dynamo:"created"`
	Finished time.Time         `dynamo:"finished"`
	Params   map[string]string `dynamo:"params,encrypted"`
	Pid      string            `dynamo:"pid"`
	Started  time.Time         `dynamo:"started"`
	Status   string            `dynamo:"status"`
	Version  string            `dynamo:"version"`
}

type Updates []Update
