// Package sdkv2schema exposes the protocol schemas of resources that are still
// implemented with terraform-plugin-sdk/v2, so that plugin-framework code can act
// on them.
//
// It exists because list resources cannot be implemented in SDKv2 at all:
// helper/schema answers every ListResource and ValidateListResourceConfig RPC with
// "list resource type is not supported by this provider", and offers no way to
// register one. A list resource for an SDKv2-managed resource therefore has to live
// in the framework provider, which then needs the managed resource's schema and
// identity schema in protocol form to satisfy list.ListResourceWithRawV6Schemas.
//
// Those schemas are read back from the SDKv2 provider itself rather than
// hand-written, so the two can never drift: whatever terraform sees for
// ionoscloud_datacenter is exactly what the list resource is told to fill in.
package sdkv2schema

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// schemas holds everything one SDKv2 provider advertises.
type schemas struct {
	resources  map[string]*tfprotov6.Schema
	identities map[string]*tfprotov6.ResourceIdentitySchema
}

// cache memoizes the schemas per SDKv2 provider instance, since building them means
// walking every resource the provider serves.
var cache sync.Map // map[*schema.Provider]*schemas

// Schemas returns the protocol v6 resource schema and resource identity schema that
// sdkv2Provider advertises for typeName, e.g. "ionoscloud_datacenter".
//
// Both are required by list.ListResourceWithRawV6Schemas. A missing identity schema
// means the SDKv2 resource does not declare a schema.ResourceIdentity yet; terraform
// identifies every result a list resource streams back by its identity, so the
// resource has to gain one before it can be listed.
func Schemas(ctx context.Context, sdkv2Provider *schema.Provider, typeName string) (*tfprotov6.Schema, *tfprotov6.ResourceIdentitySchema, error) {
	if sdkv2Provider == nil {
		return nil, nil, fmt.Errorf("no SDKv2 provider was passed to the framework provider, so %q cannot be listed", typeName)
	}

	loaded, err := load(ctx, sdkv2Provider)
	if err != nil {
		return nil, nil, err
	}

	resourceSchema, ok := loaded.resources[typeName]
	if !ok {
		return nil, nil, fmt.Errorf("no SDKv2 managed resource named %q is registered in the provider", typeName)
	}

	identitySchema, ok := loaded.identities[typeName]
	if !ok {
		return nil, nil, fmt.Errorf(
			"the SDKv2 managed resource %q does not declare a resource identity; "+
				"add an Identity to its schema.Resource so that it can be listed", typeName,
		)
	}

	return resourceSchema, identitySchema, nil
}

// load asks the SDKv2 provider, upgraded to protocol v6 the same way main.go does,
// for the schemas of every resource it serves.
func load(ctx context.Context, sdkv2Provider *schema.Provider) (*schemas, error) {
	if cached, ok := cache.Load(sdkv2Provider); ok {
		return cached.(*schemas), nil //nolint:forcetypeassert // only *schemas is ever stored
	}

	server, err := tf5to6server.UpgradeServer(ctx, sdkv2Provider.GRPCProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade the SDKv2 provider server to protocol v6: %w", err)
	}

	providerSchema, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err == nil {
		err = firstError(providerSchema.Diagnostics)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read the SDKv2 provider schema: %w", err)
	}

	identitySchemas, err := server.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err == nil {
		err = firstError(identitySchemas.Diagnostics)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read the SDKv2 resource identity schemas: %w", err)
	}

	loaded := &schemas{
		resources:  providerSchema.ResourceSchemas,
		identities: identitySchemas.IdentitySchemas,
	}
	cache.Store(sdkv2Provider, loaded)

	return loaded, nil
}

// firstError returns the first error diagnostic in diags as an error, or nil.
func firstError(diags []*tfprotov6.Diagnostic) error {
	for _, d := range diags {
		if d != nil && d.Severity == tfprotov6.DiagnosticSeverityError {
			return fmt.Errorf("%s: %s", d.Summary, d.Detail)
		}
	}

	return nil
}
