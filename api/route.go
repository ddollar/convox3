package api

import "github.com/convox/stdapi"

func (a *Api) Route(server *stdapi.Server) {
	cn := a.controller

	server.Subrouter("/organizations/{oid}/racks/{rid}/terraform", func(tf *stdapi.Router) {
		tf.Use(cn.EnsureTerraformAuthentication)

		tf.Route("POST", "/lock", cn.TerraformLock)
		tf.Route("DELETE", "/lock", cn.TerraformUnlock)
		tf.Route("GET", "/state", cn.TerraformStateLoad)
		tf.Route("POST", "/state", cn.TerraformStateStore)
	})
}
