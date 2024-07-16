package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	s3 "github.com/ionos-cloud/sdk-go-s3"
)

var (
	_ resource.ResourceWithImportState = (*s3BucketObjectResource)(nil)
	_ resource.ResourceWithConfigure   = (*s3BucketObjectResource)(nil)
)

// NewBucketObjectResource creates a new resource for the bucket resource.
func NewS3BucketObjectResource() resource.Resource {
	return &s3BucketObjectResource{}
}

type s3BucketObjectResource struct {
	client *s3.APIClient
}

type s3BucketObjectModel struct {
	BucketName types.String `tfsdk:"bucket_name"`
	ObjectKey  types.String `tfsdk:"object_key"`
	Body       types.String `tfsdk:"body"`
	BodyFile   types.String `tfsdk:"body_file"`
}

// Metadata returns the metadata for the bucket policy resource.
func (r *s3BucketObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket_object" // todo: use constant here maybe
}

// Schema returns the schema for the bucket policy resource.
func (r *s3BucketObjectResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Description: "Name of the S3 bucket to which the object will be added.",
				Required:    true,
			},
			"object_key": schema.StringAttribute{
				Description: "Key name of the object",
				Required:    true,
			},
			"body_file": schema.StringAttribute{
				Description: "The name of the file containing the body.",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "The body.",
				Computed:    true,
			},
			// "cache_control": schema.StringAttribute{
			// 	Description: "an be used to specify caching behavior along the request/reply chain.",
			// 	Optional:    true,
			// },
			// "content_disposition": schema.StringAttribute{
			// 	Description: "Specifies presentational information for the object.",
			// 	Optional:    true,
			// },
			// "content_encoding": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "content_language": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "content_length": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "content_md5": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "content_type": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "expires": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_server_side_encryption": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_storage_class": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_website_redirect_location": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_server_side_encryption_customer_algorithm": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_server_side_encryption_customer_key": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_server_side_encryption_customer_key_md5": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_server_side_encryption_context": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_request_payer": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_tagging": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_object_lock_mode": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_object_lock_retain_until_date": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
			// "x_amz_object_lock_legal_hold": schema.StringAttribute{
			// 	Description: "The name of the file containing the body.",
			// 	Optional:    true,
			// },
		},
	}
}

// Configure configures the bucket object resource.
func (r *s3BucketObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*s3.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the bucket policy.
func (r *s3BucketObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured") // todo: const for this error maybe?
		return
	}

	var data s3BucketObjectModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	f, err := os.Open(data.BodyFile.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot open file", err.Error())
		return
	}

	_, err = r.client.ObjectsApi.PutObject(ctx, data.BucketName.ValueString(), data.ObjectKey.ValueString()).Body(f).Execute()

	f.Close()
	if err != nil {
		resp.Diagnostics.AddError("failed to create bucket", err.Error())
		return
	}
}

// Read reads the bucket policy.
func (r *s3BucketObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

	var data s3BucketObjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	obj, _, err := r.client.ObjectsApi.GetObject(ctx, data.BucketName.ValueString(), data.ObjectKey.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to retrieve object bucket", err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(obj, buf)
	obj.Close()
	if err != nil {
		resp.Diagnostics.AddError("cannot read from file", err.Error())
		return
	}

	body := string(buf.Bytes())

	data.Body = types.StringValue(body)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState imports the state of the bucket policy.
func (r *s3BucketObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("s3_bucket_object"), req, resp)
}

// Update updates the bucket policy.
func (r *s3BucketObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data bucketResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Nothing to update for now
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the object.
func (r *s3BucketObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("s3 api client not configured", "The provider client is not configured")
		return
	}

}
