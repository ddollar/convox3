package controller

import (
	"github.com/convox/console/pkg/crypt"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/stdapi"
)

func (cn *Controller) EnsureTerraformAuthentication(fn stdapi.HandlerFunc) stdapi.HandlerFunc {
	return func(c *stdapi.Context) error {
		_, pass, ok := c.Request().BasicAuth()
		if !ok {
			return stdapi.Errorf(401, "authentication failed")
		}

		dec, err := crypt.Decrypt(settings.RackKey, pass)
		if err != nil {
			return stdapi.Errorf(401, "authentication failed")
		}

		if string(dec) != c.Var("rid") {
			return stdapi.Errorf(401, "authentication failed")
		}

		return fn(c)
	}
}