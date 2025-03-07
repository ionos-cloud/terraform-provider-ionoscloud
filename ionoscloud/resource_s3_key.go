package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceS3Key() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceS3KeyCreate,
		ReadContext:   resourceS3KeyRead,
		UpdateContext: resourceS3KeyUpdate,
		DeleteContext: resourceS3KeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceS3KeyImport,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:             schema.TypeString,
				Description:      "The ID of the user that owns the key.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The Object Storage Secret key.",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether this key should be active or not.",
				Optional:    true,
				Default:     true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceS3KeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	userId := d.Get("user_id").(string)
	rsp, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysPost(ctx, userId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating Object Storage key: %w", err))
		return diags
	}

	if rsp.Id == nil {
		return diag.FromErr(fmt.Errorf("the API didn't return an Object Storage key ID"))
	}
	keyId := *rsp.Id
	d.SetId(keyId)
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		return diag.FromErr(errState)
	}

	log.Printf("[INFO] Created Object Storage key: %s", d.Id())

	active := d.Get("active").(bool)
	s3Key := ionoscloud.S3Key{
		Properties: &ionoscloud.S3KeyProperties{
			Active: &active,
		},
	}
	log.Printf("[INFO] Setting key active status to %+v", active)
	_, apiResponse, err = client.UserS3KeysApi.UmUsersS3keysPut(ctx, userId, keyId).S3Key(s3Key).Depth(1).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error saving key data %s: %w", keyId, err))
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceS3KeyRead(ctx, d, meta)
}

func resourceS3KeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	userId := d.Get("user_id").(string)

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) || isS3KeyNotFound(err) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while reading Object Storage key %s: %w, %+v", d.Id(), err, s3Key))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved Object Storage key %+v \n", *s3Key.Id)

	if s3Key.HasProperties() && s3Key.Properties.HasActive() {
		log.Printf("[INFO] Successfully retrieved Object Storage key with status: %t", *s3Key.Properties.Active)
	}

	if err := setS3KeyIdAndProperties(&s3Key, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceS3KeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	request := ionoscloud.S3Key{}
	request.Properties = &ionoscloud.S3KeyProperties{}

	log.Printf("[INFO] Attempting to update Object Storage key %s", d.Id())

	newActiveSetting := d.Get("active")
	log.Printf("[INFO] Object Storage key active setting changed to %+v", newActiveSetting)
	active := newActiveSetting.(bool)
	request.Properties.Active = &active

	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysPut(ctx, userId, d.Id()).S3Key(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) || isS3KeyNotFound(err) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating Object Storage key %s: %w", d.Id(), err))
		return diags
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return diag.FromErr(errState)
	}

	return resourceS3KeyRead(ctx, d, meta)
}

func resourceS3KeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	userId := d.Get("user_id").(string)
	apiResponse, err := client.UserS3KeysApi.UmUsersS3keysDelete(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) || isS3KeyNotFound(err) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Object Storage key %s: %w", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for s3Key %s to be deleted...", d.Id())

		s3KeyDeleted, dsErr := s3KeyDeleted(ctx, client, d)

		if dsErr != nil {
			if isS3KeyNotFound(dsErr) {
				log.Printf("[INFO] Successfully deleted Object Storage key: %s", d.Id())
				return nil
			}
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of Object Storage key %s: %w", d.Id(), dsErr))
			return diags
		}

		if s3KeyDeleted {
			log.Printf("[INFO] Successfully deleted Object Storage key: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] delete timed out")
			diags := diag.FromErr(fmt.Errorf("Object Storage key delete timed out! WARNING: your Object Storage key will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return nil
}

// isS3KeyNotFound needed because api returns 422 instead of 404 on key being not found. will be removed once API issue is fixed
func isS3KeyNotFound(err error) bool {
	var genericOpenAPIError ionoscloud.GenericOpenAPIError
	if !errors.As(err, &genericOpenAPIError) {
		return false
	}
	return genericOpenAPIError.StatusCode() == 422 && strings.Contains(genericOpenAPIError.Error(), "[VDC-21-2] The access key cannot be found, please double-check the key id and try again.")
}

func s3KeyDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking Object Storage key deletion status: %w", err)
	}
	return false, nil
}

func s3Ready(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	rsp, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking Object Storage Key status: %w", err)
	}
	active := d.Get("active").(bool)
	return *rsp.Properties.Active == active, nil
}

func resourceS3KeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import id %q. Expecting {userId}/{s3KeyId}", d.Id())
	}

	userId := parts[0]
	keyId := parts[1]

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, keyId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) || isS3KeyNotFound(err) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Object Storage key %q", keyId)
		}
		return nil, fmt.Errorf("unable to retrieve Object Storage key %q, error:%w", keyId, err)
	}

	if err := setS3KeyIdAndProperties(&s3Key, d); err != nil {
		return nil, err
	}

	if err := d.Set("user_id", userId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setS3KeyIdAndProperties(s3Key *ionoscloud.S3Key, data *schema.ResourceData) error {

	if s3Key == nil {
		return fmt.Errorf("s3key not found")
	}

	if s3Key.Id != nil {
		data.SetId(*s3Key.Id)
	}

	if s3Key.Properties.SecretKey != nil {
		if err := data.Set("secret_key", *s3Key.Properties.SecretKey); err != nil {
			return err
		}
	}

	if s3Key.Properties.Active != nil {
		log.Printf("[INFO] SETTING ACTIVE TO %+v", *s3Key.Properties.Active)
		if err := data.Set("active", *s3Key.Properties.Active); err != nil {
			return err
		}
	}
	return nil
}
