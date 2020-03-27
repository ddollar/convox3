package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/convox/console/api"
	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/console/pkg/storage"
	"github.com/convox/stdapi"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	a, err := api.New(model.New(storage.New("dynamo")))
	if err != nil {
		return err
	}

	s := stdapi.New("web", "console.convox")

	s.Router.HandleFunc("/graphql", a.Handler())

	// if err := a.Route(s); err != nil {
	// 	return err
	// }

	if settings.Development {
		if err := routeWebDevelopment(s); err != nil {
			return err
		}
	}

	if err := s.Listen("https", ":3000"); err != nil {
		return err
	}

	return nil
}

func routeWebDevelopment(s *stdapi.Server) error {
	u, err := url.Parse("http://localhost:3001")
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	s.Router.Handle("/{path:.*}", rp)

	return nil
}
