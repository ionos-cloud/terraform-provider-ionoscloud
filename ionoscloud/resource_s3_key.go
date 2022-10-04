package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:         schema.TypeString,
				Description:  "The ID of the user that owns the key.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The S3 Secret key.",
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
	client := meta.(SdkBundle).CloudApiClient

	userId := d.Get("user_id").(string)
	rsp, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysPost(ctx, userId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating S3 key: %w", err))
		return diags
	}

	if rsp.Id == nil {
		return diag.FromErr(fmt.Errorf("the API didn't return an s3 key ID"))
	}
	keyId := *rsp.Id
	d.SetId(keyId)

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	log.Printf("[INFO] Created S3 key: %s", d.Id())

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

	_, errState = getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceS3KeyRead(ctx, d, meta)
}

func resourceS3KeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	userId := d.Get("user_id").(string)

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while reading S3 key %s: %w, %+v", d.Id(), err, s3Key))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved S3 key %+v \n", *s3Key.Id)

	if s3Key.HasProperties() && s3Key.Properties.HasActive() {
		log.Printf("[INFO] Successfully retrieved S3 key with status: %t", *s3Key.Properties.Active)
	}

	if err := setS3KeyIdAndProperties(&s3Key, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceS3KeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	request := ionoscloud.S3Key{}
	request.Properties = &ionoscloud.S3KeyProperties{}

	log.Printf("[INFO] Attempting to update S3 key %s", d.Id())

	newActiveSetting := d.Get("active")
	log.Printf("[INFO] S3 key active setting changed to %+v", newActiveSetting)
	active := newActiveSetting.(bool)
	request.Properties.Active = &active

	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysPut(ctx, userId, d.Id()).S3Key(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating S3 key %s: %w", d.Id(), err))
		return diags
	}

	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceS3KeyRead(ctx, d, meta)
}

func resourceS3KeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	userId := d.Get("user_id").(string)
	apiResponse, err := client.UserS3KeysApi.UmUsersS3keysDelete(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting S3 key %s: %w", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for s3Key %s to be deleted...", d.Id())

		s3KeyDeleted, dsErr := s3KeyDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of S3 key %s: %s", d.Id(), dsErr))
			return diags
		}

		if s3KeyDeleted {
			log.Printf("[INFO] Successfully deleted S3 key: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] delete timed out")
			diags := diag.FromErr(fmt.Errorf("s3 key delete timed out! WARNING: your s3 key will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return nil
}

func s3KeyDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking S3 key deletion status: %s", err)
	}
	return false, nil
}

func s3Ready(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	rsp, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking S3 Key status: %s", err)
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

	client := meta.(SdkBundle).CloudApiClient

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, keyId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find S3 key %q", keyId)
		}
		return nil, fmt.Errorf("unable to retreive S3 key %q", keyId)
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
