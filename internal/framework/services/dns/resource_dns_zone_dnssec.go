package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	dnssdk "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	dnsservice "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

var (
	_ resource.ResourceWithImportState = (*zoneDNSSECResource)(nil)
	_ resource.ResourceWithConfigure   = (*zoneDNSSECResource)(nil)
)

type zoneDNSSECResource struct {
	client *dnsservice.Client
}

type zoneDNSSECResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	ZoneID          types.String   `tfsdk:"zone_id"`
	Algorithm       types.String   `tfsdk:"algorithm"`
	KskBits         types.Int64    `tfsdk:"ksk_bits"`
	ZskBits         types.Int64    `tfsdk:"zsk_bits"`
	NsecMode        types.String   `tfsdk:"nsec_mode"`
	Nsec3Iterations types.Int64    `tfsdk:"nsec3_iterations"`
	Nsec3SaltBits   types.Int64    `tfsdk:"nsec3_salt_bits"`
	Validity        types.Int64    `tfsdk:"validity"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	dnssecKeyModel
}

// NewZoneDNSSECResource creates a new resource managing the DNSSEC signing key of a DNS zone.
func NewZoneDNSSECResource() resource.Resource {
	return &zoneDNSSECResource{}
}

// Metadata returns the metadata for the resource.
func (r *zoneDNSSECResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone_dnssec"
}

// replaceIfChanged reports whether a configured value that differs from the one in state requires replacement.
// The API does not return every creation parameter, so after an import the state holds nulls; those are adopted
// from the configuration instead of forcing the key to be recreated.
func replaceIfChanged(stateIsNull, planIsNull, equal bool) bool {
	return !stateIsNull && !planIsNull && !equal
}

func requiresReplaceString() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = replaceIfChanged(req.StateValue.IsNull(), req.PlanValue.IsNull(), req.PlanValue.Equal(req.StateValue))
		},
		"Changing this value recreates the DNSSEC key.", "Changing this value recreates the DNSSEC key.",
	)
}

func requiresReplaceInt64() planmodifier.Int64 {
	return int64planmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.Int64Request, resp *int64planmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = replaceIfChanged(req.StateValue.IsNull(), req.PlanValue.IsNull(), req.PlanValue.Equal(req.StateValue))
		},
		"Changing this value recreates the DNSSEC key.", "Changing this value recreates the DNSSEC key.",
	)
}

// Schema returns the schema for the resource.
func (r *zoneDNSSECResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	computedString := func(description string) schema.StringAttribute {
		return schema.StringAttribute{
			Computed:      true,
			Description:   description,
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		}
	}
	computedInt64 := func(description string) schema.Int64Attribute {
		return schema.Int64Attribute{
			Computed:      true,
			Description:   description,
			PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
		}
	}
	bitsValidator := int64validator.OneOf(1024, 2048, 4096)

	resp.Schema = schema.Schema{
		Description: "Enables DNSSEC for a DNS zone by managing its signing key. The attributes of the key needed to build the DS record for the registrar are exported.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "The ID of the resource, equal to `zone_id`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"zone_id": schema.StringAttribute{
				Required:      true,
				Description:   "The ID (UUID) of the DNS zone to enable DNSSEC for. Changing it recreates the resource.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"algorithm": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Default:       stringdefault.StaticString(string(dnssdk.ALGORITHM_RSASHA256)),
				Description:   "The signing algorithm. Only `RSASHA256` is supported. Defaults to `RSASHA256`.",
				Validators:    []validator.String{stringvalidator.OneOf(string(dnssdk.ALGORITHM_RSASHA256))},
				PlanModifiers: []planmodifier.String{requiresReplaceString()},
			},
			"ksk_bits": schema.Int64Attribute{
				Required:      true,
				Description:   "The size in bits of the key signing key. One of `1024`, `2048`, `4096`.",
				Validators:    []validator.Int64{bitsValidator},
				PlanModifiers: []planmodifier.Int64{requiresReplaceInt64()},
			},
			"zsk_bits": schema.Int64Attribute{
				Required:      true,
				Description:   "The size in bits of the zone signing key. One of `1024`, `2048`, `4096`.",
				Validators:    []validator.Int64{bitsValidator},
				PlanModifiers: []planmodifier.Int64{requiresReplaceInt64()},
			},
			"nsec_mode": schema.StringAttribute{
				Required:      true,
				Description:   "The authenticated denial of existence mode. One of `NSEC`, `NSEC3`.",
				Validators:    []validator.String{stringvalidator.OneOf(string(dnssdk.NSECMODE_NSEC), string(dnssdk.NSECMODE_NSEC3))},
				PlanModifiers: []planmodifier.String{requiresReplaceString()},
			},
			"nsec3_iterations": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Default:       int64default.StaticInt64(0),
				Description:   "The number of NSEC3 iterations. Defaults to `0`.",
				PlanModifiers: []planmodifier.Int64{requiresReplaceInt64()},
			},
			"nsec3_salt_bits": schema.Int64Attribute{
				Optional:      true,
				Computed:      true,
				Default:       int64default.StaticInt64(64),
				Description:   "The NSEC3 salt length in bits. Defaults to `64`.",
				PlanModifiers: []planmodifier.Int64{requiresReplaceInt64()},
			},
			"validity": schema.Int64Attribute{
				Required:      true,
				Description:   "The signature validity in days.",
				PlanModifiers: []planmodifier.Int64{requiresReplaceInt64()},
			},
			"key_tag":                   computedInt64(descKeyTag),
			"flags":                     computedInt64(descFlags),
			"public_key":                computedString(descPublicKey),
			"composed_key_data":         computedString(descComposedKeyData),
			"digest":                    computedString(descDigest),
			"digest_algorithm_mnemonic": computedString(descDigestMnemonic),
			"digest_type":               computedInt64(descDigestType),
			"algorithm_number":          computedInt64(descAlgorithmNumber),
			"ds_record":                 computedString(descDSRecord),
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

// Configure configures the resource.
func (r *zoneDNSSECResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *bundleclient.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = clientBundle.DNSClient
}

// Create enables DNSSEC for the zone.
func (r *zoneDNSSECResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data zoneDNSSECResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, utils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := data.ZoneID.ValueString()
	errCtx := &diagutil.ErrorContext{ResourceID: zoneID, Timeout: createTimeout.String()}

	// The zone has to be settled before a key can be added to it.
	if err := retry(ctx, createTimeout, func() error { return r.client.IsZoneAvailable(ctx, zoneID) }); err != nil {
		resp.Diagnostics.AddError("error while waiting for the DNS zone to become AVAILABLE", diagutil.WrapError(err, errCtx).Error())
		return
	}

	properties := dnssdk.DnssecKeyParameters{
		KeyParameters: dnssdk.KeyParameters{
			Algorithm: dnssdk.Algorithm(data.Algorithm.ValueString()),
			KskBits:   dnssdk.KskBits(data.KskBits.ValueInt64()),
			ZskBits:   dnssdk.ZskBits(data.ZskBits.ValueInt64()),
		},
		NsecParameters: dnssdk.NsecParameters{
			NsecMode:        dnssdk.NsecMode(data.NsecMode.ValueString()),
			Nsec3Iterations: int32(data.Nsec3Iterations.ValueInt64()),
			Nsec3SaltBits:   int32(data.Nsec3SaltBits.ValueInt64()),
		},
		Validity: int32(data.Validity.ValueInt64()),
	}

	_, apiResponse, err := r.client.CreateDNSSECKey(ctx, zoneID, properties)
	if err != nil {
		errCtx.StatusCode = apiResponse.SafeStatusCode()
		summary := "failed to enable DNSSEC for the DNS zone"
		if errCtx.StatusCode == http.StatusConflict {
			summary = "DNSSEC is already enabled for the DNS zone; import the existing key with `terraform import` or remove it first"
		}
		resp.Diagnostics.AddError(summary, diagutil.WrapError(err, errCtx).Error())
		return
	}

	// The key material is not part of the creation response, so wait for it to show up in the zone.
	var keys dnssdk.DnssecKeyReadList
	err = retry(ctx, createTimeout, func() error {
		if err := r.client.IsZoneAvailable(ctx, zoneID); err != nil {
			return err
		}
		var getResponse *shared.APIResponse
		keys, getResponse, err = r.client.GetDNSSECKeys(ctx, zoneID)
		if err != nil {
			// The key is created asynchronously; until then the API reports the zone as not signed.
			if getResponse.HttpNotFound() || dnsservice.IsZoneNotSigned(getResponse) {
				return err
			}
			return backoff.Permanent(err)
		}
		if _, found := dnsservice.SigningKey(keys); !found {
			return errors.New(zoneNotFoundMessage(zoneID))
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("error while waiting for the DNSSEC key to become available", diagutil.WrapError(err, errCtx).Error())
		return
	}

	data.ID = data.ZoneID
	data.dnssecKeyModel.setFromKeys(keys, dnssdk.Algorithm(data.Algorithm.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the state from the API. Parameters that the API does not return stay as they are in state.
func (r *zoneDNSSECResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data zoneDNSSECResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := data.ZoneID.ValueString()
	keys, apiResponse, err := r.client.GetDNSSECKeys(ctx, zoneID)
	if err != nil {
		if apiResponse.HttpNotFound() || dnsservice.IsZoneNotSigned(apiResponse) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error while fetching the DNSSEC key of the DNS zone", diagutil.WrapError(err, &diagutil.ErrorContext{ResourceID: zoneID, StatusCode: apiResponse.SafeStatusCode()}).Error())
		return
	}

	found, algorithm, nsecMode := data.dnssecKeyModel.setFromKeys(keys, dnssdk.Algorithm(data.Algorithm.ValueString()))
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	data.ID = data.ZoneID
	data.Algorithm = types.StringValue(string(algorithm))
	if nsecMode != nil {
		data.NsecMode = types.StringValue(string(*nsecMode))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update only runs to adopt configured values for parameters that the API does not return (e.g. after an import),
// or when only the timeouts changed. Every real change of the key recreates it.
func (r *zoneDNSSECResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan zoneDNSSECResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete disables DNSSEC for the zone.
func (r *zoneDNSSECResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data zoneDNSSECResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneID := data.ZoneID.ValueString()
	errCtx := &diagutil.ErrorContext{ResourceID: zoneID}

	apiResponse, err := r.client.DeleteDNSSECKey(ctx, zoneID)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return
		}
		errCtx.StatusCode = apiResponse.SafeStatusCode()
		resp.Diagnostics.AddError("failed to disable DNSSEC for the DNS zone", diagutil.WrapError(err, errCtx).Error())
		return
	}

	// The API only acknowledges the request; the key is removed asynchronously and the zone stays signed until the
	// change has propagated, which may take several hours. Do not block on it.
	resp.Diagnostics.AddWarning(
		"DNSSEC removal is asynchronous",
		"The API accepted the request to disable DNSSEC. The key is removed asynchronously, which may take several hours, and enabling DNSSEC again for the zone fails until then.",
	)
}

// ImportState imports the DNSSEC key of a zone using the zone ID.
func (r *zoneDNSSECResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), req.ID)...)
}

// retry runs op with exponential backoff until it succeeds, returns a permanent error, or the timeout elapses.
func retry(ctx context.Context, timeout time.Duration, op func() error) error {
	return backoff.Retry(op, backoff.WithContext(backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(timeout)), ctx))
}
