package manager

import (
	"io"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/storage"
)

var (
	m = model.New(storage.New("dynamo"))
)

type Manager interface {
	Install(name, version, region string, params map[string]string, output io.Writer) error
	Uninstall(output io.Writer) error
	Update(name, version string, params map[string]string, output io.Writer) error
}

func New(engine, rid string) (Manager, error) {
	switch engine {
	case "v2":
		return NewV2(rid)
	default:
		return NewV3(rid)
	}
}
