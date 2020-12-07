package ionoscloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
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
				Type:        schema.TypeString,
				Description: "The ID of the user that owns the key.",
				Required:    true,
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
	client := meta.(*profitbricks.Client)

	createdS3Key, err := client.CreateS3Key(d.Get("user_id").(string))

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating S3 key: %s", err)
	}

	d.SetId(createdS3Key.ID)
	log.Printf("[INFO] Created S3 key: %s", d.Id())

	return resourceS3KeyRead(d, meta)
}

func resourceS3KeyRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*profitbricks.Client)
	s3Key, err := client.GetS3Key(d.Get("user_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}

		return fmt.Errorf("Error while reading S3 key %s: %s, %+v", d.Id(), err, s3Key)
	}

	log.Printf("[INFO] Successfully retreived S3 key %s: %+v", d.Id(), s3Key)

	d.SetId(s3Key.ID)
	d.Set("secret_key", s3Key.Properties.SecretKey)
	d.Set("active", s3Key.Properties.Active)

	return nil
}

func resourceS3KeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	request := profitbricks.S3Key{}

	request.Properties = &profitbricks.S3KeyProperties{}

	log.Printf("[INFO] Attempting to update S3 key %s", d.Id())

	if d.HasChange("active") {
		oldActiveSetting, newActiveSetting := d.GetChange("active")
		log.Printf("[INFO] S3 key active setting changed from %+v to %+v", oldActiveSetting, newActiveSetting)
		request.Properties.Active = newActiveSetting.(bool)
	}

	updatedS3Key, err := client.UpdateS3Key(d.Get("user_id").(string), d.Id(), request)

	time.Sleep(5 * time.Second)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating S3 key: %s", err)
		}
		return fmt.Errorf("Error while updating S3 key %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for S3 Key %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		s3KeyReady, rsErr := s3Ready(client, d, updatedS3Key)

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
	client := meta.(*profitbricks.Client)

	_, err := client.DeleteS3Key(d.Get("user_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting S3 key: %s", err)
		}

		return fmt.Errorf("Error while deleting S3 key %s: %s", d.Id(), err)
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

func s3KeyDeleted(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	_, err := client.GetS3Key(d.Get("user_id").(string), d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking S3 key deletion status: %s", err)
		}
	}
	return false, nil
}

func s3Ready(client *profitbricks.Client, d *schema.ResourceData, s3Key *profitbricks.S3Key) (bool, error) {
	subjectS3Key, err := client.GetS3Key(d.Get("user_id").(string), d.Id())

	if err != nil {
		return true, fmt.Errorf("Error checking S3 Key status: %s", err)
	}
	return subjectS3Key.Properties.Active == d.Get("active").(bool), nil
}
