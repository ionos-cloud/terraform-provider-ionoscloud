package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
	rsp, apiResponse, err := createS3KeyWithRetry(ctx, d, meta)
	if err != nil {
		d.SetId("")
		return utils.ToDiags(d, fmt.Sprintf("error creating Object Storage key: %s", err), nil)
	}

	if rsp.Id == nil {
		return utils.ToDiags(d, "the API didn't return an Object Storage key ID", nil)
	}
	keyId := *rsp.Id
	d.SetId(keyId)
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutCreate})
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
		requestLocation, _ := apiResponse.Location()
		return utils.ToDiags(d, fmt.Sprintf("error saving key data %s: %s", keyId, err), &utils.DiagsOpts{RequestLocation: requestLocation, StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
	}

	return resourceS3KeyRead(ctx, d, meta)
}

// createS3KeyWithRetry uses retry.RetryContext to attempt to create an S3 key,
// specifically retrying on privilege propagation errors.
func createS3KeyWithRetry(ctx context.Context, d *schema.ResourceData, meta any) (ionoscloud.S3Key, *ionoscloud.APIResponse, error) {
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	userId := d.Get("user_id").(string)

	var s3Key ionoscloud.S3Key
	var apiResponse *ionoscloud.APIResponse
	var err error

	retryErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		s3Key, apiResponse, err = client.UserS3KeysApi.UmUsersS3keysPost(ctx, userId).Execute()
		logApiRequestTime(apiResponse)

		if err == nil {
			// Success
			return nil
		}

		if isS3KeyPrivilegeError(err) {
			log.Printf("[INFO] Retrying S3 key creation due to privilege error: %v", err)
			return retry.RetryableError(err)
		}

		// Any other error is not retryable.
		return retry.NonRetryableError(err)
	})

	if retryErr != nil {
		return ionoscloud.S3Key{}, apiResponse, retryErr
	}

	return s3Key, apiResponse, nil
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
		return utils.ToDiags(d, fmt.Sprintf("error while reading Object Storage key: %s, %+v", err, s3Key), nil)
	}

	log.Printf("[INFO] Successfully retrieved Object Storage key %+v \n", *s3Key.Id)

	if s3Key.HasProperties() && s3Key.Properties.HasActive() {
		log.Printf("[INFO] Successfully retrieved Object Storage key with status: %t", *s3Key.Properties.Active)
	}

	if err := setS3KeyIdAndProperties(&s3Key, d); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
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
		return utils.ToDiags(d, fmt.Sprintf("error while updating Object Storage key: %s", err), nil)
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		return utils.ToDiags(d, errState.Error(), &utils.DiagsOpts{Timeout: schema.TimeoutUpdate})
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
		return utils.ToDiags(d, fmt.Sprintf("error while deleting Object Storage key: %s", err), nil)
	}

	for {
		log.Printf("[INFO] Waiting for s3Key %s to be deleted...", d.Id())

		s3KeyDeleted, dsErr := s3KeyDeleted(ctx, client, d)

		if dsErr != nil {
			if isS3KeyNotFound(dsErr) {
				log.Printf("[INFO] Successfully deleted Object Storage key: %s", d.Id())
				return nil
			}
			return utils.ToDiags(d, fmt.Sprintf("error while checking deletion status of Object Storage key: %s", dsErr), nil)
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
			return utils.ToDiags(d, "Object Storage key delete timed out! WARNING: your Object Storage key will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", nil)
		}
	}

	return nil
}

// isS3KeyPrivilegeError checks if the error is the specific 422 Unprocessable Entity
// error that indicates a privilege propagation delay IF IONOS_S3_KEY_CREATION_RETRY is set to true
func isS3KeyPrivilegeError(err error) bool {
	envVal, ok := os.LookupEnv("IONOS_S3_KEY_CREATION_RETRY")
	retryEnabled := false
	if ok {
		if b, err := strconv.ParseBool(strings.ToLower(strings.TrimSpace(envVal))); err == nil {
			retryEnabled = b
		} else {
			log.Printf("[WARN] invalid IONOS_S3_KEY_CREATION_RETRY value %q; defaulting to false", envVal)
		}
	}
	if !retryEnabled {
		return false
	}
	var genericOpenAPIError ionoscloud.GenericOpenAPIError
	if !errors.As(err, &genericOpenAPIError) {
		return false
	}

	return genericOpenAPIError.StatusCode() == http.StatusUnprocessableEntity && strings.Contains(string(genericOpenAPIError.Body()), "The user needs to be part of a group that has ACCESS_S3_OBJECT_STORAGE privilege")
}

// isS3KeyNotFound needed because api returns 422 instead of 404 on key being not found. will be removed once API issue is fixed
func isS3KeyNotFound(err error) bool {
	var genericOpenAPIError ionoscloud.GenericOpenAPIError
	if !errors.As(err, &genericOpenAPIError) {
		return false
	}
	return genericOpenAPIError.StatusCode() == 422 && strings.Contains(genericOpenAPIError.Error(), "The access key cannot be found, please double-check the key id and try again.")
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
		return nil, utils.ToError(d, "invalid import. Expecting {userId}/{s3KeyId}", nil)
	}

	userId := parts[0]
	keyId := parts[1]

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	s3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userId, keyId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) || isS3KeyNotFound(err) {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find Object Storage key %q", keyId), nil)
		}
		return nil, utils.ToError(d, fmt.Sprintf("unable to retrieve Object Storage key %q, error:%s", keyId, err), nil)
	}

	if err := setS3KeyIdAndProperties(&s3Key, d); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	if err := d.Set("user_id", userId); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
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
