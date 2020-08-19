package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/convox/console/api"
	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/console/pkg/storage"
	"github.com/convox/stdapi"
	"github.com/gobuffalo/packr/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	a, err := api.New(model.New(storage.New("dynamo")), packr.New("graphql", "../../api/graphql"))
	if err != nil {
		return err
	}

	s := stdapi.New("web", "console.convox")

	s.Router.HandleFunc("/graphql", a.Handler())

	a.Route(s)

	if settings.Development {
		if err := routeAssetsDevelopment(s); err != nil {
			return err
		}
	} else {
		if err := routeAssetsProduction(s); err != nil {
			return err
		}
	}

	if err := s.Listen("https", ":3000"); err != nil {
		return err
	}

	return nil
}

func routeAssetsDevelopment(s *stdapi.Server) error {
	u, err := url.Parse("http://localhost:3001")
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(u)

	s.Router.Handle("/{path:.*}", rp)

	return nil
}

func routeAssetsProduction(s *stdapi.Server) error {
	s.Router.Static("", spaFileSystem{packr.New("dist", "../../web/dist")})

	return nil
}

type spaFileSystem struct {
	http.FileSystem
}

func (fs spaFileSystem) Open(name string) (http.File, error) {
	if file, err := fs.FileSystem.Open(name); err == nil {
		return file, nil
	} else {
		return fs.FileSystem.Open("index.html")
	}
}
