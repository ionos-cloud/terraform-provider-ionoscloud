package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

var _ datasource.DataSourceWithConfigure = (*contractsDataSource)(nil)

type contractsDataSource struct {
	client *ionoscloud.APIClient
}

// NewContractsDataSource creates a new data source for the contracts resource.
func NewContractsDataSource() datasource.DataSource {
	return &contractsDataSource{}
}

// Metadata returns the metadata for the data source.
func (d *contractsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contracts"
}

// Schema returns the schema for the data source.
func (d *contractsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"contracts": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of contracts attached to your Ionos account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"contract_number": schema.Int64Attribute{
							Computed:    true,
							Description: "The contract number.",
						},
						"owner": schema.StringAttribute{
							Computed:    true,
							Description: "The contract owner's user name.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The contract status.",
						},
						"reg_domain": schema.StringAttribute{
							Computed:    true,
							Description: "The registration domain of the contract.",
						},
						"resource_limits": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The resource limits of the contract.",
							Attributes: map[string]schema.Attribute{
								"cores_per_server": schema.Int64Attribute{
									Computed:    true,
									Description: "The maximum number of cores per server.",
								},
								"ram_per_server": schema.Int64Attribute{
									Computed:    true,
									Description: "The maximum RAM per server in MB.",
								},
								"ram_per_contract": schema.Int64Attribute{
									Computed:    true,
									Description: "The maximum RAM per contract in MB.",
								},
								"cores_per_contract": schema.Int64Attribute{
									Computed:    true,
									Description: "The maximum number of cores per contract.",
								},
								"cores_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of cores provisioned.",
								},
								"das_volume_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The DAS volume provisioned.",
								},
								"hdd_limit_per_contract": schema.Int64Attribute{
									Computed:    true,
									Description: "The HDD limit per contract.",
								},
								"hdd_limit_per_volume": schema.Int64Attribute{
									Computed:    true,
									Description: "The HDD limit per volume.",
								},
								"hdd_volume_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The HDD volume provisioned.",
								},
								"k8s_cluster_limit_total": schema.Int64Attribute{
									Computed:    true,
									Description: "The total Kubernetes cluster limit.",
								},
								"k8s_clusters_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of Kubernetes clusters provisioned.",
								},
								"nat_gateway_limit_total": schema.Int64Attribute{
									Computed:    true,
									Description: "The total NAT gateway limit.",
								},
								"nat_gateway_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of NAT gateways provisioned.",
								},
								"nlb_limit_total": schema.Int64Attribute{
									Computed:    true,
									Description: "The total NLB limit.",
								},
								"nlb_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of NLBs provisioned.",
								},
								"ram_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The RAM provisioned.",
								},
								"reservable_ips": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of reservable IPs.",
								},
								"reserved_ips_in_use": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of reserved IPs in use.",
								},
								"reserved_ips_on_contract": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of reserved IPs on the contract.",
								},
								"ssd_limit_per_contract": schema.Int64Attribute{
									Computed:    true,
									Description: "The SSD limit per contract.",
								},
								"ssd_limit_per_volume": schema.Int64Attribute{
									Computed:    true,
									Description: "The SSD limit per volume.",
								},
								"ssd_volume_provisioned": schema.Int64Attribute{
									Computed:    true,
									Description: "The SSD volume provisioned.",
								},
								"security_groups_per_vdc": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of security groups per VDC.",
								},
								"security_groups_per_resource": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of security groups per resource.",
								},
								"rules_per_security_group": schema.Int64Attribute{
									Computed:    true,
									Description: "The number of rules per security group.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *contractsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T", req.ProviderData),
		)
		return
	}

	apiClient, err := client.NewCloudAPIClient("")
	if err != nil {
		resp.Diagnostics.AddError("Cloud API Client Error", err.Error())
		return
	}
	d.client = apiClient
}

// contractAttributesModel defines the attributes for a contract.
var contractAttributesModel = map[string]attr.Type{
	"contract_number": types.Int64Type,
	"owner":           types.StringType,
	"status":          types.StringType,
	"reg_domain":      types.StringType,
	"resource_limits": types.ObjectType{AttrTypes: resourceLimitsModel},
}
var resourceLimitsModel = map[string]attr.Type{
	"cores_per_server":             types.Int64Type,
	"ram_per_server":               types.Int64Type,
	"ram_per_contract":             types.Int64Type,
	"cores_per_contract":           types.Int64Type,
	"cores_provisioned":            types.Int64Type,
	"das_volume_provisioned":       types.Int64Type,
	"hdd_limit_per_contract":       types.Int64Type,
	"hdd_limit_per_volume":         types.Int64Type,
	"hdd_volume_provisioned":       types.Int64Type,
	"k8s_cluster_limit_total":      types.Int64Type,
	"k8s_clusters_provisioned":     types.Int64Type,
	"nat_gateway_limit_total":      types.Int64Type,
	"nat_gateway_provisioned":      types.Int64Type,
	"nlb_limit_total":              types.Int64Type,
	"nlb_provisioned":              types.Int64Type,
	"ram_provisioned":              types.Int64Type,
	"reservable_ips":               types.Int64Type,
	"reserved_ips_in_use":          types.Int64Type,
	"reserved_ips_on_contract":     types.Int64Type,
	"ssd_limit_per_contract":       types.Int64Type,
	"ssd_limit_per_volume":         types.Int64Type,
	"ssd_volume_provisioned":       types.Int64Type,
	"security_groups_per_vdc":      types.Int64Type,
	"security_groups_per_resource": types.Int64Type,
	"rules_per_security_group":     types.Int64Type,
}

func (d *contractsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("API client not configured", "The provider client is not configured")
		return
	}

	contracts, apiResponse, err := d.client.ContractResourcesApi.ContractsGet(ctx).Execute()
	apiResponse.LogInfo()
	if err != nil {
		resp.Diagnostics.AddError("Error reading contracts", fmt.Sprintf("Could not read contracts, unexpected error: %s", err.Error()))
		return
	}
	if contracts.Items == nil {
		resp.Diagnostics.AddError("Error reading contracts", "No contracts found")
		return
	}
	contractList := make([]attr.Value, len(*contracts.Items))
	var diags diag.Diagnostics
	for i, contract := range *contracts.Items {
		resourceLimits := types.ObjectNull(resourceLimitsModel)
		if contract.Properties == nil {
			resp.Diagnostics.AddError("Error reading contracts", "Contract properties are nil")
			return
		}
		if contract.Properties.ResourceLimits != nil {
			resourceLimits, diags = types.ObjectValue(resourceLimitsModel, map[string]attr.Value{
				"cores_per_server":             types.Int64Value(int64(*contract.Properties.ResourceLimits.CoresPerServer)),
				"ram_per_server":               types.Int64Value(int64(*contract.Properties.ResourceLimits.RamPerServer)),
				"ram_per_contract":             types.Int64Value(int64(*contract.Properties.ResourceLimits.RamPerContract)),
				"cores_per_contract":           types.Int64Value(int64(*contract.Properties.ResourceLimits.CoresPerContract)),
				"cores_provisioned":            types.Int64Value(int64(*contract.Properties.ResourceLimits.CoresProvisioned)),
				"das_volume_provisioned":       types.Int64Value(*contract.Properties.ResourceLimits.DasVolumeProvisioned),
				"hdd_limit_per_contract":       types.Int64Value(*contract.Properties.ResourceLimits.HddLimitPerContract),
				"hdd_limit_per_volume":         types.Int64Value(*contract.Properties.ResourceLimits.HddLimitPerVolume),
				"hdd_volume_provisioned":       types.Int64Value(*contract.Properties.ResourceLimits.HddVolumeProvisioned),
				"k8s_cluster_limit_total":      types.Int64Value(int64(*contract.Properties.ResourceLimits.K8sClusterLimitTotal)),
				"k8s_clusters_provisioned":     types.Int64Value(int64(*contract.Properties.ResourceLimits.K8sClustersProvisioned)),
				"nat_gateway_limit_total":      types.Int64Value(int64(*contract.Properties.ResourceLimits.NatGatewayLimitTotal)),
				"nat_gateway_provisioned":      types.Int64Value(int64(*contract.Properties.ResourceLimits.NatGatewayProvisioned)),
				"nlb_limit_total":              types.Int64Value(int64(*contract.Properties.ResourceLimits.NlbLimitTotal)),
				"nlb_provisioned":              types.Int64Value(int64(*contract.Properties.ResourceLimits.NlbProvisioned)),
				"ram_provisioned":              types.Int64Value(int64(*contract.Properties.ResourceLimits.RamProvisioned)),
				"reservable_ips":               types.Int64Value(int64(*contract.Properties.ResourceLimits.ReservableIps)),
				"reserved_ips_in_use":          types.Int64Value(int64(*contract.Properties.ResourceLimits.ReservedIpsInUse)),
				"reserved_ips_on_contract":     types.Int64Value(int64(*contract.Properties.ResourceLimits.ReservedIpsOnContract)),
				"ssd_limit_per_contract":       types.Int64Value(*contract.Properties.ResourceLimits.SsdLimitPerContract),
				"ssd_limit_per_volume":         types.Int64Value(*contract.Properties.ResourceLimits.SsdLimitPerVolume),
				"ssd_volume_provisioned":       types.Int64Value(*contract.Properties.ResourceLimits.SsdVolumeProvisioned),
				"security_groups_per_vdc":      types.Int64Value(int64(*contract.Properties.ResourceLimits.SecurityGroupsPerVdc)),
				"security_groups_per_resource": types.Int64Value(int64(*contract.Properties.ResourceLimits.SecurityGroupsPerResource)),
				"rules_per_security_group":     types.Int64Value(int64(*contract.Properties.ResourceLimits.RulesPerSecurityGroup)),
			})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		contractData, diags := types.ObjectValue(contractAttributesModel, map[string]attr.Value{
			"contract_number": types.Int64Value(*contract.Properties.ContractNumber),
			"owner":           types.StringValue(*contract.Properties.Owner),
			"status":          types.StringValue(*contract.Properties.Status),
			"reg_domain":      types.StringValue(*contract.Properties.RegDomain),
			"resource_limits": resourceLimits,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		contractList[i] = contractData
	}

	contractsValue, diags := types.ListValue(types.ObjectType{AttrTypes: contractAttributesModel}, contractList)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateValue, diags := types.ObjectValue(map[string]attr.Type{
		"contracts": types.ListType{ElemType: types.ObjectType{AttrTypes: contractAttributesModel}},
	}, map[string]attr.Value{
		"contracts": contractsValue,
	})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateValue)...)
}
