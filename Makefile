.PHONY: all build clean clean-package compress dist mocks package release test

commands = web
dist = $(wildcard web/dist/*)
webpack  = public/assets.js

assets   = $(wildcard assets/*)
binaries = $(addprefix $(GOPATH)/bin/, $(commands))
sources  = $(shell find . -name '*.go')

MODE ?= production

all: build

build: $(binaries) # $(webpack)

clean: clean-package

clean-package:
	find . -name '*-packr.go' -delete

compress: $(binaries)
	upx-ucl -1 $^

dist:
	cd web && npm run build

mocks:
	# make -C models mocks
	# make -C pkg/storage mocks

package:
	go run -mod=vendor vendor/github.com/gobuffalo/packr/v2/packr2/main.go

release:
	test -n "$(VERSION)" # VERSION
	git tag $(VERSION) -m $(VERSION)
	git push origin refs/tags/$(VERSION)

test:
	env TEST=true go test -covermode atomic -coverprofile coverage.txt -mod=vendor ./...

$(binaries): $(GOPATH)/bin/%: $(sources)
	go install -mod=vendor --ldflags="-s -w" ./cmd/$*

$(GOPATH)/bin/web: dist
