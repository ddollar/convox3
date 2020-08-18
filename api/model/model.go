package model

import (
	"github.com/convox/console/pkg/storage"
	"github.com/convox/convox/sdk"
)

type Model struct {
	rack    *sdk.Client
	storage storage.Interface
}

func New(storage storage.Interface) *Model {
	r, err := sdk.NewFromEnv()
	if err != nil {
		panic(err)
	}

	return &Model{rack: r, storage: storage}
}
