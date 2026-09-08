package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/identity"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// A list resource for ionoscloud_target_group, whose managed resource is still
// implemented with terraform-plugin-sdk/v2 and therefore lives on the other half of
// the mux. The protocol schemas the framework needs come from that resource via
// identity.SetRawV6Schemas, and the Identity it requires is declared alongside it in
// ionoscloud/resource_target_group.go.

const targetGroupResourceType = "ionoscloud_target_group"

var (
	_ list.ListResource                 = (*targetGroupListResource)(nil)
	_ list.ListResourceWithConfigure    = (*targetGroupListResource)(nil)
	_ list.ListResourceWithRawV6Schemas = (*targetGroupListResource)(nil)
)

// targetGroupListResource lists ionoscloud_target_group instances.
type targetGroupListResource struct {
	bundle *bundleclient.SdkBundle

	// resourceSchema is the SDKv2 managed resource being listed, the source of the
	// protocol schemas returned by RawV6Schemas.
	resourceSchema *schema.Resource
}

// targetGroupResourceModel mirrors the SDKv2 schema of ionoscloud_target_group. It has
// to cover every attribute and block in that schema, because the framework fills the
// whole resource object from it.
//
// The fields are plain Go pointers rather than the types.String / types.Int64 values
// a framework-native resource would use. The schema behind this model is converted
// from the SDKv2 protocol schema, where an SDKv2 TypeInt arrives as a protocol
// Number and becomes a NumberAttribute - assigning a types.Int64 to it fails the
// type check. Plain Go types reflect into whichever framework type the converted
// schema ended up with, and the pointer carries the null/absent distinction.
//
// health_check and http_health_check are MaxItems: 1 blocks, but a MaxItems: 1 list is
// still NestingList in the protocol (core_schema.go:201-205), so they are slices here
// and decode as single-element lists, not bare objects.
type targetGroupResourceModel struct {
	ID              *string                            `tfsdk:"id"`
	Name            *string                            `tfsdk:"name"`
	Algorithm       *string                            `tfsdk:"algorithm"`
	Protocol        *string                            `tfsdk:"protocol"`
	ProtocolVersion *string                            `tfsdk:"protocol_version"`
	Targets         *[]targetGroupTargetModel          `tfsdk:"targets"`
	HealthCheck     *[]targetGroupHealthCheckModel     `tfsdk:"health_check"`
	HTTPHealthCheck *[]targetGroupHTTPHealthCheckModel `tfsdk:"http_health_check"`
	Timeouts        *targetGroupTimeoutsModel          `tfsdk:"timeouts"`
}

// targetGroupTargetModel mirrors an entry of the targets block.
type targetGroupTargetModel struct {
	IP                 *string `tfsdk:"ip"`
	Port               *int32  `tfsdk:"port"`
	Weight             *int32  `tfsdk:"weight"`
	ProxyProtocol      *string `tfsdk:"proxy_protocol"`
	HealthCheckEnabled *bool   `tfsdk:"health_check_enabled"`
	MaintenanceEnabled *bool   `tfsdk:"maintenance_enabled"`
}

// targetGroupHealthCheckModel mirrors the health_check block.
type targetGroupHealthCheckModel struct {
	CheckTimeout  *int32 `tfsdk:"check_timeout"`
	CheckInterval *int32 `tfsdk:"check_interval"`
	Retries       *int32 `tfsdk:"retries"`
}

// targetGroupHTTPHealthCheckModel mirrors the http_health_check block.
type targetGroupHTTPHealthCheckModel struct {
	Path      *string `tfsdk:"path"`
	Method    *string `tfsdk:"method"`
	MatchType *string `tfsdk:"match_type"`
	Response  *string `tfsdk:"response"`
	Regex     *bool   `tfsdk:"regex"`
	Negate    *bool   `tfsdk:"negate"`
}

// targetGroupTimeoutsModel mirrors the timeouts block that SDKv2 adds to the schema.
// A listed target group has no timeouts, so this is always left null.
type targetGroupTimeoutsModel struct {
	Create  *string `tfsdk:"create"`
	Default *string `tfsdk:"default"`
	Delete  *string `tfsdk:"delete"`
	Update  *string `tfsdk:"update"`
}

// NewTargetGroupListResource creates a new list resource for ionoscloud_target_group.
// targetGroupResource is the SDKv2 managed resource being listed; it is the source of
// the protocol schemas the framework needs.
func NewTargetGroupListResource(targetGroupResource *schema.Resource) list.ListResource {
	return &targetGroupListResource{resourceSchema: targetGroupResource}
}

// RawV6Schemas hands the framework the protocol schemas of the SDKv2 managed
// resource. A framework-native list resource inherits them from the resource
// itself; this is only needed because ionoscloud_target_group lives on the SDKv2 side.
func (r *targetGroupListResource) RawV6Schemas(ctx context.Context, _ list.RawV6SchemaRequest, resp *list.RawV6SchemaResponse) {
	identity.SetRawV6Schemas(ctx, resp, targetGroupResourceType, r.resourceSchema)
}

// Metadata returns the type name of the managed resource being listed. It must match
// the SDKv2 resource exactly, otherwise terraform has no resource to attach the
// results to.
func (r *targetGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = targetGroupResourceType
}

// Configure stores the client bundle shared by the provider.
func (r *targetGroupListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected List Resource Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.bundle = clientBundle
}

// ListResourceConfigSchema returns the schema for the list resource config block.
func (r *targetGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			identity.FiltersKey: identity.FilterAttribute("name", "algorithm"),
		},
	}
}

// List fetches every target group on the contract and streams the results. Target
// groups are not location-scoped - the resource has no location attribute at all - so
// there is nothing to fan out over here.
func (r *targetGroupListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	identity.StreamList(ctx, stream, req,
		func(ctx context.Context) ([]ionoscloud.TargetGroup, error) {
			// The target group collection is not location-scoped, so this uses the same
			// client every other global Cloud API listing uses - and the same one the
			// managed resource's own CRUD paths use.
			client, err := r.bundle.NewCloudAPIClientWithFailover(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create the Cloud API client: %w", err)
			}

			// Depth(1) is what makes the API return the properties of every target group
			// instead of just its links. Filtering stays client-side, in the mapper.
			//
			// The explicit Limit matches ionoscloud/data_source_target_group.go. Without
			// it the SDK client falls back to limit=100 (sdk-go/v6/api_target_groups.go:489-492,
			// which is the client's own default, not the endpoint's). Note that
			// constant.TargetGroupLimit is 200, a much tighter ceiling than the 1000 the
			// datacenter and ip block listings get - a contract with more than 200 target
			// groups is truncated silently, which the docs call out.
			targetGroups, apiResponse, err := client.TargetGroupsApi.TargetgroupsGet(ctx).Depth(1).Limit(constant.TargetGroupLimit).Execute()
			if apiResponse != nil {
				tflog.Debug(ctx, "listed target groups", map[string]any{"status_code": apiResponse.SafeStatusCode()})
			}
			if err != nil {
				return nil, fmt.Errorf("failed to list target groups: %w", err)
			}
			if targetGroups.Items == nil {
				return nil, nil
			}

			return *targetGroups.Items, nil
		},
		r.mapTargetGroup,
	)
}

// mapTargetGroup maps a target group to an identity.MappedItem, or returns nil to skip it.
func (r *targetGroupListResource) mapTargetGroup(_ context.Context, includeResource bool, filters []identity.Filter, targetGroup ionoscloud.TargetGroup) (*identity.MappedItem, diag.Diagnostics) {
	if targetGroup.Id == nil || targetGroup.Properties == nil {
		return nil, nil
	}

	name := shared.ToValueDefault(targetGroup.Properties.Name)
	algorithm := shared.ToValueDefault(targetGroup.Properties.Algorithm)

	// protocol is deliberately NOT filterable: the API allows exactly one value, "HTTP"
	// (model_target_group_properties.go:23, and the resource's own validator at
	// resource_target_group.go:63), so a protocol filter could never narrow anything.
	if !identity.MatchesFilters(map[string]string{
		"name":      name,
		"algorithm": algorithm,
	}, filters) {
		return nil, nil
	}

	mapped := &identity.MappedItem{
		// name is Required on ionoscloud_target_group, so it is always populated and
		// needs no fallback - unlike ionoscloud_ipblock.
		DisplayName: name,
		// Target groups are addressed by UUID alone, so the shared single-`id` identity
		// model fits and no per-resource struct is needed.
		Identity: &identity.Model{ID: types.StringValue(*targetGroup.Id)},
	}

	if !includeResource {
		return mapped, nil
	}

	mapped.Resource = &targetGroupResourceModel{
		ID:              targetGroup.Id,
		Name:            targetGroup.Properties.Name,
		Algorithm:       targetGroup.Properties.Algorithm,
		Protocol:        targetGroup.Properties.Protocol,
		ProtocolVersion: targetGroup.Properties.ProtocolVersion,
		Targets:         mapTargetGroupTargets(targetGroup.Properties.Targets),
		HealthCheck:     mapTargetGroupHealthCheck(targetGroup.Properties.HealthCheck),
		HTTPHealthCheck: mapTargetGroupHTTPHealthCheck(targetGroup.Properties.HttpHealthCheck),
	}

	return mapped, nil
}

// mapTargetGroupTargets maps the targets balanced by a target group.
func mapTargetGroupTargets(targets *[]ionoscloud.TargetGroupTarget) *[]targetGroupTargetModel {
	if targets == nil {
		return nil
	}

	mapped := make([]targetGroupTargetModel, 0, len(*targets))
	for _, target := range *targets {
		mapped = append(mapped, targetGroupTargetModel{
			IP:                 target.Ip,
			Port:               target.Port,
			Weight:             target.Weight,
			ProxyProtocol:      target.ProxyProtocol,
			HealthCheckEnabled: target.HealthCheckEnabled,
			MaintenanceEnabled: target.MaintenanceEnabled,
		})
	}

	return &mapped
}

// mapTargetGroupHealthCheck maps the health check of a target group. The SDKv2 schema
// declares it as a MaxItems: 1 list, so a single object becomes a one-element slice.
func mapTargetGroupHealthCheck(healthCheck *ionoscloud.TargetGroupHealthCheck) *[]targetGroupHealthCheckModel {
	if healthCheck == nil {
		return nil
	}

	return &[]targetGroupHealthCheckModel{{
		CheckTimeout:  healthCheck.CheckTimeout,
		CheckInterval: healthCheck.CheckInterval,
		Retries:       healthCheck.Retries,
	}}
}

// mapTargetGroupHTTPHealthCheck maps the HTTP health check of a target group, likewise
// a MaxItems: 1 list in the SDKv2 schema.
func mapTargetGroupHTTPHealthCheck(healthCheck *ionoscloud.TargetGroupHttpHealthCheck) *[]targetGroupHTTPHealthCheckModel {
	if healthCheck == nil {
		return nil
	}

	return &[]targetGroupHTTPHealthCheckModel{{
		Path:      healthCheck.Path,
		Method:    healthCheck.Method,
		MatchType: healthCheck.MatchType,
		Response:  healthCheck.Response,
		Regex:     healthCheck.Regex,
		Negate:    healthCheck.Negate,
	}}
}
