module github.com/ionos-cloud/terraform-provider-ionoscloud/v6

go 1.21

require (
	github.com/aws/aws-sdk-go v1.54.12
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/gofrs/uuid/v5 v5.0.0
	github.com/hashicorp/terraform-plugin-framework v1.9.0
	github.com/hashicorp/terraform-plugin-framework-timetypes v0.4.0
	github.com/hashicorp/terraform-plugin-framework-validators v0.12.0
	github.com/hashicorp/terraform-plugin-mux v0.16.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.33.0
	github.com/hashicorp/terraform-plugin-testing v1.8.0
	github.com/iancoleman/strcase v0.3.0
	github.com/ionos-cloud/sdk-go-bundle/products/logging/v2 v2.0.0
	github.com/ionos-cloud/sdk-go-bundle/shared v0.1.0
	github.com/ionos-cloud/sdk-go-cert-manager v1.0.1
	github.com/ionos-cloud/sdk-go-container-registry v1.1.0
	github.com/ionos-cloud/sdk-go-dataplatform v1.0.3
	github.com/ionos-cloud/sdk-go-dbaas-mariadb v1.0.1
	github.com/ionos-cloud/sdk-go-dbaas-mongo v1.3.1
	github.com/ionos-cloud/sdk-go-dbaas-postgres v1.1.2
	github.com/ionos-cloud/sdk-go-dns v1.1.1
	github.com/ionos-cloud/sdk-go-s3 v1.0.0
	github.com/ionos-cloud/sdk-go-vm-autoscaling v1.0.0
	github.com/ionos-cloud/sdk-go/v6 v6.1.11
	golang.org/x/crypto v0.23.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	gopkg.in/validator.v2 v2.0.1 // indirect
)

replace github.com/ionos-cloud/sdk-go-s3 => /home/dinga/Work/repos/sdk-resources/s3-go

replace github.com/ionos-cloud/sdk-go-bundle/shared => /home/dinga/Work/repos/shared

require (
	github.com/ProtonMail/go-crypto v1.1.0-alpha.2 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-test/deep v1.0.6 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.6.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hc-install v0.6.4 // indirect
	github.com/hashicorp/hcl/v2 v2.20.1 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-exec v0.21.0 // indirect
	github.com/hashicorp/terraform-json v0.22.1 // indirect
	github.com/hashicorp/terraform-plugin-go v0.23.0
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.2.3 // indirect
	github.com/hashicorp/terraform-svchost v0.1.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.34.0 // indirect
)
