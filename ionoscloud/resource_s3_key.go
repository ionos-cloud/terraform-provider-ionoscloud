package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceS3Key() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3KeyCreate,
		Read:   resourceS3KeyRead,
		Update: resourceS3KeyUpdate,
		Delete: resourceS3KeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceS3KeyImport,
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

func resourceS3KeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}

	rsp, _, err := client.UserManagementApi.UmUsersS3keysPost(ctx, d.Get("user_id").(string)).Execute()

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating S3 key: %s", err)
	}

	d.SetId(*rsp.Id)
	log.Printf("[INFO] Created S3 key: %s", d.Id())

	return resourceS3KeyRead(d, meta)
}

func resourceS3KeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	userId := d.Get("user_id").(string)

	rsp, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, d.Id()).Execute()
	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}

		return fmt.Errorf("error while reading S3 key %s: %s, %+v", d.Id(), err, rsp)
	}

	log.Printf("[INFO] Successfully retreived S3 key %s: %+v", d.Id(), rsp)

	d.SetId(*rsp.Id)
	d.Set("secret_key", *rsp.Properties.SecretKey)
	d.Set("active", *rsp.Properties.Active)

	return nil
}

func resourceS3KeyUpdate(d *schema.ResourceData, meta interface{}) error {
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
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("error while updating S3 key: %s", err)
		}
		return fmt.Errorf("error while updating S3 key %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for S3 Key %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		s3KeyReady, rsErr := s3Ready(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of S3 Key %s: %s", d.Id(), rsErr)
		}

		if s3KeyReady && rsErr == nil {
			log.Printf("[INFO] S3 Key ready: %s", d.Id())
			break
		}
	}

	return resourceS3KeyRead(d, meta)
}

func resourceS3KeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserManagementApi.UmUsersS3keysDelete(ctx, userId, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("error while deleting S3 key: %s", err)
		}

		return fmt.Errorf("error while deleting S3 key %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for s3Key %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		s3KeyDeleted, dsErr := s3KeyDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of S3 key %s: %s", d.Id(), dsErr)
		}

		if s3KeyDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted S3 key: %s", d.Id())
			break
		}
	}

	return nil
}

func s3KeyDeleted(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	_, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				return true, nil
			}
			return true, fmt.Errorf("error checking S3 key deletion status: %s", err)
		}
	}
	return false, nil
}

func s3Ready(client *ionoscloud.APIClient, d *schema.ResourceData) (bool, error) {
	userId := d.Get("user_id").(string)
	rsp, _, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("Error checking S3 Key status: %s", err)
	}
	active := d.Get("active").(bool)
	return *rsp.Properties.Active == active, nil
}
