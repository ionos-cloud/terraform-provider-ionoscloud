package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceS3KeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	rsp, apiResponse, err := client.UserManagementApi.UmUsersS3keysPost(ctx, d.Get("user_id").(string)).Execute()

	if err != nil {
		d.SetId("")
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error creating S3 key: %s %s", err, payload))
		return diags
	}

	if rsp.Id != nil {
		d.SetId(*rsp.Id)
	}
	log.Printf("[INFO] Created S3 key: %s", d.Id())

	return resourceS3KeyRead(ctx, d, meta)
}

func resourceS3KeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	userId := d.Get("user_id").(string)

	rsp, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error while reading S3 key %s: %s, %+v %s", d.Id(), err, rsp, payload))
		return diags
	}

	log.Printf("[INFO] Successfully retreived S3 key %s: %+v", d.Id(), rsp)

	if rsp.Id != nil {
		d.SetId(*rsp.Id)
	}

	if rsp.Properties.SecretKey != nil {
		if err := d.Set("secret_key", *rsp.Properties.SecretKey); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	if rsp.Properties.Active != nil {
		if err := d.Set("active", *rsp.Properties.Active); err != nil {
			diags := diag.FromErr(err)
			return diags
		}
	}

	return nil
}

func resourceS3KeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.S3Key{}
	request.Properties = &ionoscloud.S3KeyProperties{}

	log.Printf("[INFO] Attempting to update S3 key %s", d.Id())

	if d.HasChange("active") {
		oldActiveSetting, newActiveSetting := d.GetChange("active")
		log.Printf("[INFO] S3 key active setting changed from %+v to %+v", oldActiveSetting, newActiveSetting)
		active := newActiveSetting.(bool)
		request.Properties.Active = &active
	}

	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserManagementApi.UmUsersS3keysPut(context.TODO(), userId, d.Id()).S3Key(request).Execute()

	time.Sleep(5 * time.Second)

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error while updating S3 key: %s %s", err, payload))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for S3 Key %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		s3KeyReady, rsErr := s3Ready(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of S3 Key %s: %s", d.Id(), rsErr))
			return diags
		}

		if s3KeyReady {
			log.Printf("[INFO] S3 Key ready: %s", d.Id())
			break
		}
	}

	return resourceS3KeyRead(ctx, d, meta)
}

func resourceS3KeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserManagementApi.UmUsersS3keysDelete(ctx, userId, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting S3 key: %s %s", err, payload))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for s3Key %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		s3KeyDeleted, dsErr := s3KeyDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of S3 key %s: %s", d.Id(), dsErr))
			return diags
		}

		if s3KeyDeleted {
			log.Printf("[INFO] Successfully deleted S3 key: %s", d.Id())
			break
		}
	}

	return nil
}

func s3KeyDeleted(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		return true, fmt.Errorf("%s %s", err, payload)
	}
	return false, nil
}

func s3Ready(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	rsp, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(ctx, userId, d.Id()).Execute()

	if err != nil {
		payload := ""
		if apiResponse != nil {
			payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
		}
		return true, fmt.Errorf("%s %s", err, payload)
	}
	active := d.Get("active").(bool)
	return *rsp.Properties.Active == active, nil
}
