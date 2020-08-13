package storage

import (
	"github.com/convox/console/pkg/logger"
	"github.com/convox/console/pkg/settings"
)

var (
	log = logger.New("ns=storage")
)

type Defaulter interface {
	Defaults()
}

type Sortable interface {
	Less(i, j int) bool
}

type Validator interface {
	Validate() []error
}

func New(provider string) Interface {
	switch provider {
	case "dynamo":
		return NewDynamo(settings.TablePrefix, settings.RackKey)
	default:
		return NewDynamo(settings.TablePrefix, settings.RackKey)
	}
}
