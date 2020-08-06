module github.com/convox/console

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v36.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Azure/go-autorest/autorest/adal v0.8.0
	github.com/Microsoft/hcsshim v0.8.7 // indirect
	github.com/aws/aws-sdk-go v1.29.28
	github.com/containerd/containerd v1.3.3 // indirect
	github.com/containerd/continuity v0.0.0-20200228182428-0f16d7a0959c // indirect
	github.com/convox/convox v0.0.0-20200319135330-fa65d5bb20ba
	github.com/convox/rack v0.0.0-20200720124016-e57ae50684af // indirect
	github.com/convox/stdapi v1.0.0
	github.com/convox/stdsdk v0.0.0-20190422120437-3e80a397e377 // indirect
	github.com/digitalocean/godo v1.42.0
	github.com/fsouza/go-dockerclient v1.6.3 // indirect
	github.com/gbrlsnchs/jwt/v3 v3.0.0-rc.2
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/graph-gophers/graphql-go v0.0.0-20200309224638-dae41bde9ef9
	github.com/graph-gophers/graphql-transport-ws v0.0.0-20190611222414-40c048432299
	github.com/hashicorp/go-retryablehttp v0.6.4 // indirect
	github.com/hashicorp/hcl/v2 v2.6.0
	github.com/karrick/godirwalk v1.15.5 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5
	github.com/xanzy/go-gitlab v0.11.3
	golang.org/x/crypto v0.0.0-20200317142112-1b76d66859c6
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20200319210407-521f4a0cd458 // indirect
	google.golang.org/api v0.9.0
	google.golang.org/genproto v0.0.0-20200319113533-08878b785e9c // indirect
	google.golang.org/grpc v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/graph-gophers/graphql-transport-ws => github.com/convox/graphql-transport-ws v0.0.3
