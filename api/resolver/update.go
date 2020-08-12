package resolver

import (
	"github.com/convox/console/api/model"
	"github.com/graph-gophers/graphql-go"
)

type Update struct {
	model.Update
}

func (u *Update) Created() int32 {
	return int32(u.Update.Created.Unix())
}

func (u *Update) Finished() *int32 {
	if u.Update.Finished.IsZero() {
		return nil
	}

	t := int32(u.Update.Finished.Unix())

	return &t
}

func (u *Update) Id() graphql.ID {
	return graphql.ID(u.Update.ID)
}

func (u *Update) Parameters() []*Parameter {
	ps := []*Parameter{}

	for k, v := range u.Update.Params {
		ps = append(ps, &Parameter{key: k, value: v})
	}

	return ps
}

func (u *Update) Pid() *string {
	if u.Update.Pid == "" {
		return nil
	}

	pid := u.Update.Pid

	return &pid
}

func (u *Update) Started() *int32 {
	if u.Update.Started.IsZero() {
		return nil
	}

	t := int32(u.Update.Started.Unix())

	return &t
}

func (u *Update) Status() string {
	return u.Update.Status
}

func (u *Update) Version() string {
	return u.Update.Version
}
