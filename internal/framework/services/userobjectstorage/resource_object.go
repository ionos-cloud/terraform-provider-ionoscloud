package userobjectstorage

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	objectstorage "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/userobjectstorage"
)

var (
	_ resource.ResourceWithImportState      = (*objectResource)(nil)
	_ resource.ResourceWithConfigure        = (*objectResource)(nil)
	_ resource.ResourceWithConfigValidators = (*objectResource)(nil)
)

// NewObjectResource creates a new resource for the object resource.
func NewObjectResource() resource.Resource {
	return &objectResource{}
}

type objectResource struct {
	client *objectstorage.Client
}

// Metadata returns the metadata for the object resource.
func (r *objectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_object"
}

// Schema returns the schema for the object resource.
func (r *objectResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "The key of the object",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"source": schema.StringAttribute{
				Description: "The path to the file to upload",
				Optional:    true,
			},
			"content": schema.StringAttribute{
				Description: "The utf-8 content of the object",
				Optional:    true,
			},
			"cache_control": schema.StringAttribute{
				Description: "Can be used to specify caching behavior along the request/reply chain",
				Optional:    true,
			},
			"content_disposition": schema.StringAttribute{
				Description: "Specifies presentational information for the object",
				Optional:    true,
			},
			"content_encoding": schema.StringAttribute{
				Description: "Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field",
				Optional:    true,
			},
			"content_language": schema.StringAttribute{
				Description: "The natural language or languages of the intended audience for the object",
				Optional:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "A standard MIME type describing the format of the contents",
				Optional:    true,
				Computed:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The date and time at which the object is no longer cacheable",
				Optional:    true,
			},
			"server_side_encryption": schema.StringAttribute{
				Description: "The server-side encryption algorithm used when storing this object in IONOS Object Storage Object Storage (AES256).",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"storage_class": schema.StringAttribute{
				Description: "The storage class of the object. Valid value is 'STANDARD'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("STANDARD"),
				Validators:  []validator.String{stringvalidator.OneOf("STANDARD")},
			},
			"website_redirect": schema.StringAttribute{
				Description: "If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. IONOS Object Storage Object Storage stores the value of this header in the object metadata",
				Optional:    true,
			},
			"server_side_encryption_customer_algorithm": schema.StringAttribute{
				Description: "Specifies the algorithm to use to when encrypting the object (e.g., AES256).",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("AES256")},
			},
			"server_side_encryption_customer_key": schema.StringAttribute{
				Description: "Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data",
				Optional:    true,
			},
			"server_side_encryption_customer_key_md5": schema.StringAttribute{
				Description: "Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321. IONOS Object Storage Object Storage uses this header for a message integrity check  to ensure that the encryption key was transmitted without error",
				Optional:    true,
			},
			"server_side_encryption_context": schema.StringAttribute{
				Description: " Specifies the IONOS Object Storage Object Storage Encryption Context to use for object encryption. The value of this header is a base64-encoded UTF-8 string holding JSON with the encryption context key-value pairs.",
				Optional:    true,
				Sensitive:   true,
			},
			"request_payer": schema.StringAttribute{
				Description: "Confirms that the requester knows that they will be charged for the request. Bucket owners need not specify this parameter in their requests.",
				Optional:    true,
			},
			"object_lock_mode": schema.StringAttribute{
				Description: "Confirms that the requester knows that they will be charged for the request. Bucket owners need not specify this parameter in their requests.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("GOVERNANCE", "COMPLIANCE")},
			},
			"object_lock_retain_until_date": schema.StringAttribute{
				Description: " The date and time when you want this object's Object Lock to expire. Must be formatted as a timestamp parameter.",
				Optional:    true,
			},
			"object_lock_legal_hold": schema.StringAttribute{
				Description: "Specifies whether a legal hold will be applied to this object.",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("ON", "OFF")},
			},
			"etag": schema.StringAttribute{
				Description: "An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "The tag-set for the object",
				Optional:    true,
				ElementType: types.StringType,
			},
			"metadata": schema.MapAttribute{
				Description: "A map of metadata to store with the object in IONOS Object Storage Object Storage",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Map{mapvalidator.ValueStringsAre([]validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z0-9]+$`), "metadata keys must be lowercase alphanumeric characters"),
				}...)},
			},
			"version_id": schema.StringAttribute{
				Description: "The version of the object",
				Computed:    true,
			},
			"mfa": schema.StringAttribute{
				Description: "The concatenation of the authentication device's serial number, a space, and the value that is displayed on your authentication device. Required to permanently delete a versioned object if versioning is configured with MFA Delete enabled.",
				Optional:    true,
			},
			"force_destroy": schema.BoolAttribute{
				Description: "Specifies whether to delete the object even if it has a governance-type Object Lock in place. You must explicitly pass a value of true for this parameter to delete the object.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// Configure configures the object resource.
func (r *objectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBundle, ok := req.ProviderData.(*bundleclient.SdkBundle)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *services.SdkBundle, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = clientBundle.UserObjectStorageClient
}

func (r *objectResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("source"),
			path.MatchRoot("content"),
		),
	}
}

// Create creates the object.
func (r *objectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UploadObject(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("failed to create object", formatXMLError(err).Error())
		return
	}

	if err = r.client.SetObjectComputedAttributes(ctx, data, result); err != nil {
		resp.Diagnostics.AddError("failed to set computed attributes", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the object.
func (r *objectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, found, err := r.client.GetObject(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("failed to read object", formatXMLError(err).Error())
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	data = result
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the object.
func (r *objectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	bucket, key, err := splitImportID(id)
	if err != nil {
		resp.Diagnostics.AddError("invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket"), bucket)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), key)...)
}

// Update updates the object.
func (r *objectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *objectstorage.ObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateObject(ctx, plan, state); err != nil {
		resp.Diagnostics.AddError("failed to update object", formatXMLError(err).Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the object.
func (r *objectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("object storage api client not configured", "The provider client is not configured")
		return
	}

	var data *objectstorage.ObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteObject(ctx, data); err != nil {
		resp.Diagnostics.AddError("failed to delete object", formatXMLError(err).Error())
		return
	}
}

func splitImportID(path string) (string, string, error) {
	// Ensure the path is not empty
	if path == "" {
		return "", "", fmt.Errorf("path cannot be empty")
	}

	// Split the path into two parts: bucket and the remaining key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid path format")
	}

	bucket := parts[0]
	key := parts[1]

	return bucket, key, nil
}
