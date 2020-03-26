package model

import "github.com/convox/console/pkg/storage"

type Model struct {
	storage storage.Interface
}

func New(storage storage.Interface) *Model {
	return &Model{storage: storage}
}
