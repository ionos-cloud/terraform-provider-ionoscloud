package s3

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/s3"
)

var (
	_ resource.ResourceWithImportState = (*objectCopyResource)(nil)
	_ resource.ResourceWithConfigure   = (*objectCopyResource)(nil)
)

// NewObjectCopyResource creates a new resource for the object copy resource.
func NewObjectCopyResource() resource.Resource {
	return &objectCopyResource{}
}

type objectCopyResource struct {
	client *s3.Client
}

// Metadata returns the metadata for the object copy resource.
func (r *objectCopyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_object_copy"
}

// Schema returns the schema for the object copy resource.
func (r *objectCopyResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				Description: "The name of the bucket",
				Required:    true,
				Validators:  []validator.String{stringvalidator.LengthBetween(3, 63)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The key of the object copy",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"source": schema.StringAttribute{
				Description: "The key of the source object",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"copy_if_match": schema.StringAttribute{
				Description: "Copies the object if its entity tag (ETag) matches the specified tag",
				Optional:    true,
			},
			"copy_if_modified_since": schema.StringAttribute{
				Description: "Copies the object if it has been modified since the specified time",
				Optional:    true,
			},
			"copy_if_none_match": schema.StringAttribute{
				Description: "Copies the object if its entity tag (ETag) is different than the specified ETag",
				Optional:    true,
			},
			"copy_if_unmodified_since": schema.StringAttribute{
				Description: "Copies the object if it hasn't been modified since the specified time",
				Optional:    true,
			},
			"cache_control": schema.StringAttribute{
				Description: "Can be used to specify caching behavior along the request/reply chain",
				Optional:    true,
			},
			"content_disposition": schema.StringAttribute{
				Description: "Specifies presentational information for the object copy",
				Optional:    true,
			},
			"content_encoding": schema.StringAttribute{
				Description: "Specifies what content encodings have been applied to the object copy and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field",
				Optional:    true,
			},
			"content_language": schema.StringAttribute{
				Description: "The natural language or languages of the intended audience for the object copy",
				Optional:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "A standard MIME type describing the format of the contents",
				Optional:    true,
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The date and time at which the object copy is no longer cacheable",
				Optional:    true,
			},
			"metadata_directive": schema.StringAttribute{
				Description: "Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("COPY", "REPLACE")},
			},
			"tagging_directive": schema.StringAttribute{
				Description: "Specifies whether the object copy tag-set is copied from the source object or replaced with tag-set provided in the request",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("COPY", "REPLACE")},
			},
			"server_side_encryption": schema.StringAttribute{
				Description: "The server-side encryption algorithm used when storing this object copy in IONOS S3 Object Copy Storage (AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"storage_class": schema.StringAttribute{
				Description: "The storage class of the object copy. Valid value is 'STANDARD'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Validators:  []validator.String{stringvalidator.OneOf("STANDARD")},
			},
			"website_redirect": schema.StringAttribute{
				Description: "If the bucket is configured as a website, redirects requests for this object copy to another object copy in the same bucket or to an external URL. IONOS S3 Object Copy Storage stores the value of this header in the object copy metadata",
				Optional:    true,
			},
			"server_side_encryption_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when encrypting the object copy (e.g., AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"server_side_encryption_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data",
				Optional:    true,
			},
			"server_side_encryption_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Copy Storage uses this header for a message integrity check  to ensure that the encryption key was transmitted without error",
				Optional:    true,
			},
			"source_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when decrypting the source object (e.g., AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"source_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to decrypt the source object",
				Optional:    true,
			},
			"source_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS S3 Object Copy Storage uses this header for a message integrity check  to ensure that the encryption key was transmitted without error",
				Optional:    true,
			},
			"object_lock_mode": schema.StringAttribute{
				Description: "Confirms that the requester knows that they will be charged for the request. Bucket owners need not specify this parameter in their requests.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
			},
			"object_lock_retain_until_date": schema.StringAttribute{
				Description: " The date and time when you want this object copy's Object Copy Lock to expire. Must be formatted as a timestamp parameter.",
				Optional:    true,
			},
			"object_lock_legal_hold": schema.StringAttribute{
				Description: "Specifies whether a legal hold will be applied to this object copy.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("ON", "OFF")},
			},
			"etag": schema.StringAttribute{
				Description: "An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.",
				Computed:    true,
			},
			"last_modified": schema.StringAttribute{
				Description: "The date and time at which the object copy was last modified",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "The tag-set for the object copy",
				Optional:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.MapAttribute{
				Description: "A map of metadata to store with the object copy in IONOS S3 Object Copy Storage",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{mapvalidator.ValueStringsAre([]validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9]+$`), "metadata keys must be lowercase alphanumeric characters"),
				}...)},
			},
			"version_id": schema.StringAttribute{
				Description: "The version of the object copy",
				Computed:    true,
			},
			"force_destroy": schema.BoolAttribute{
				Description: "Specifies whether to delete the object copy even if it has a governance-type Object Copy Lock in place. You must explicitly pass a value of true for this parameter to delete the object copy.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// Configure configures the object copy resource.
func (r *objectCopyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *s3.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the object copy.
func (r *objectCopyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *s3.ObjectCopyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CopyObject(ctx, data); err != nil {
		resp.Diagnostics.AddError("failed to create object copy", formatXMLError(err).Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the object copy.
func (r *objectCopyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *s3.ObjectCopyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := r.client.GetObjectCopy(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("failed to read object copy", formatXMLError(err).Error())
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	data = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the object copy.
func (r *objectCopyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	bucket, key, err := splitImportID(id)
	if err != nil {
		resp.Diagnostics.AddError("invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket"), bucket)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
}

// Update updates the object copy.
func (r *objectCopyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *s3.ObjectCopyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateObjectCopy(ctx, plan, state); err != nil {
		resp.Diagnostics.AddError("failed to update object copy", formatXMLError(err).Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the object copy.
func (r *objectCopyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data *s3.ObjectCopyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteObjectCopy(ctx, data); err != nil {
		resp.Diagnostics.AddError("failed to delete object copy", formatXMLError(err).Error())
		return
	}
}
