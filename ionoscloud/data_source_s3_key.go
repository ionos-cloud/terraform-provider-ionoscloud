package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceObjectStorageKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectStorageKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Id of the key.",
				Computed:    true,
				Optional:    true,
			},
			"user_id": {
				Type:             schema.TypeString,
				Description:      "The ID of the user that owns the key.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				Required:         true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The Secret key.",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether this key should be active or not.",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceObjectStorageKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	userIDItf, idOk := d.GetOk("user_id")
	if !idOk {
		return utils.ToDiags(d, "please provide the userID", nil)
	}
	userID := userIDItf.(string)
	client := meta.(bundleclient.SdkBundle).CloudApiClient
	var s3Keys ionoscloud.S3Keys
	var s3Key ionoscloud.S3Key
	var err error
	var apiResponse *ionoscloud.APIResponse
	if IDItf, idOk := d.GetOk("id"); idOk {
		id := IDItf.(string)
		s3Key, apiResponse, err = client.UserS3KeysApi.UmUsersS3keysFindByKeyId(ctx, userID, id).Execute()
		apiResponse.LogInfo()
		if err != nil {
			if apiResponse.HttpNotFound() || isS3KeyNotFound(err) {
				return utils.ToDiags(d, fmt.Sprintf("no storage key found with the specified criteria: userID = %s id = %s", userID, id), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
			}
			return utils.ToDiags(d, fmt.Sprintf("error while reading Object Storage key: %s, %s", err, userID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		s3Keys, apiResponse, err = client.UserS3KeysApi.UmUsersS3keysGet(ctx, userID).Depth(2).Execute()
		apiResponse.LogInfo()
		if apiResponse.HttpNotFound() || isS3KeyNotFound(err) {
			return utils.ToDiags(d, fmt.Sprintf("no storage key found with the specified criteria: userID = %s", userID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("error while reading Object Storage key: %s, %s", err, userID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		if s3Keys.Items == nil || len(*s3Keys.Items) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no storage key found with the specified criteria: userID = %s", userID), nil)
		} else if len(*s3Keys.Items) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one storage key found with the specified criteria: userID = %s", userID), nil)
		}

		s3Key = (*s3Keys.Items)[0]
	}

	if err := setS3KeyIdAndProperties(&s3Key, d); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
