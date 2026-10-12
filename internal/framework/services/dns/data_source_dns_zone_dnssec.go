package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dnssdk "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	dnsservice "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var _ datasource.DataSourceWithConfigure = (*zoneDNSSECDataSource)(nil)

type zoneDNSSECDataSource struct {
	client *dnsservice.Client
}

type zoneDNSSECDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	ZoneID    types.String `tfsdk:"zone_id"`
	Algorithm types.String `tfsdk:"algorithm"`
	NsecMode  types.String `tfsdk:"nsec_mode"`
	dnssecKeyModel
}

// NewZoneDNSSECDataSource creates a new data source reading the DNSSEC signing key of a DNS zone.
func NewZoneDNSSECDataSource() datasource.DataSource {
	return &zoneDNSSECDataSource{}
}

// Metadata returns the metadata for the data source.
func (d *zoneDNSSECDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone_dnssec"
}

// Configure configures the data source.
func (d *zoneDNSSECDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = clientBundle.DNSClient
}

// Schema returns the schema for the data source.
func (d *zoneDNSSECDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the DNSSEC signing key of a DNS zone, including the values needed for the DS record.",
		Attributes: map[string]schema.Attribute{
			"id":                        schema.StringAttribute{Computed: true, Description: "The ID of the data source, equal to `zone_id`."},
			"zone_id":                   schema.StringAttribute{Required: true, Description: "The ID (UUID) of the DNS zone."},
			"algorithm":                 schema.StringAttribute{Computed: true, Description: "The signing algorithm."},
			"nsec_mode":                 schema.StringAttribute{Computed: true, Description: "The authenticated denial of existence mode."},
			"key_tag":                   schema.Int64Attribute{Computed: true, Description: descKeyTag},
			"flags":                     schema.Int64Attribute{Computed: true, Description: descFlags},
			"public_key":                schema.StringAttribute{Computed: true, Description: descPublicKey},
			"composed_key_data":         schema.StringAttribute{Computed: true, Description: descComposedKeyData},
			"digest":                    schema.StringAttribute{Computed: true, Description: descDigest},
			"digest_algorithm_mnemonic": schema.StringAttribute{Computed: true, Description: descDigestMnemonic},
			"digest_type":               schema.Int64Attribute{Computed: true, Description: descDigestType},
			"algorithm_number":          schema.Int64Attribute{Computed: true, Description: descAlgorithmNumber},
			"ds_record":                 schema.StringAttribute{Computed: true, Description: descDSRecord},
		},
	}
}

// Read reads the DNSSEC key of the zone.
func (d *zoneDNSSECDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data zoneDNSSECDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := data.ZoneID.ValueString()
	keys, apiResponse, err := d.client.GetDNSSECKeys(ctx, zoneID)
	if err != nil {
		if dnsservice.IsZoneNotSigned(apiResponse) {
			resp.Diagnostics.AddError("DNSSEC key not found", dnssecNotEnabledMessage(zoneID))
			return
		}
		resp.Diagnostics.AddError("failed to get the DNSSEC key of the DNS zone", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceID: zoneID, StatusCode: apiResponse.SafeStatusCode()}).Error())
		return
	}

	found, algorithm, nsecMode, diags := data.setFromKeys(keys, dnssdk.ALGORITHM_RSASHA256)
	resp.Diagnostics.Append(diags...)
	if !found {
		resp.Diagnostics.AddError("DNSSEC key not found", dnssecNotEnabledMessage(zoneID))
		return
	}

	data.ID = data.ZoneID
	data.Algorithm = types.StringValue(string(algorithm))
	data.NsecMode = types.StringNull()
	if nsecMode != nil {
		data.NsecMode = types.StringValue(string(*nsecMode))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
