package ionoscloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceBackupUnit() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackupUnitCreate,
		Read:   resourceBackupUnitRead,
		Update: resourceBackupUnitUpdate,
		Delete: resourceBackupUnitDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBackupUnitImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Alphanumeric name you want assigned to the backup unit.",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password you want assigned to the backup unit.",
				Required:    true,
				Sensitive:   true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The e-mail address you want assigned to the backup unit.",
				Required:    true,
			},
			"login": {
				Type:        schema.TypeString,
				Description: "The login associated with the backup unit. Derived from the contract number",
				Computed:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceBackupUnitCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	backupUnit := profitbricks.BackupUnit{
		Properties: &profitbricks.BackupUnitProperties{
			Name:     d.Get("name").(string),
			Password: d.Get("password").(string),
			Email:    d.Get("email").(string),
		},
	}

	createdBackupUnit, err := client.CreateBackupUnit(backupUnit)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error creating backup unit: %s", err)
	}

	d.SetId(createdBackupUnit.ID)
	log.Printf("[INFO] Created backup unit: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		backupUnitReady, rsErr := backupUnitReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of backup unit %s: %s", d.Id(), rsErr)
		}

		if backupUnitReady && rsErr == nil {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}
	}

	return resourceBackupUnitRead(d, meta)
}

func resourceBackupUnitRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*profitbricks.Client)
	backupUnit, err := client.GetBackupUnit(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}

		return fmt.Errorf("Error while fetching backup unit %s: %s", d.Id(), err)
	}

	contractResources, cErr := client.GetContractResources()

	if cErr != nil {
		return fmt.Errorf("Error while fetching contract resources for backup unit %s: %s", d.Id(), cErr)
	}

	log.Printf("[INFO] Successfully retreived contract resource for backup unit unit %s: %+v", d.Id(), contractResources)

	d.Set("name", backupUnit.Properties.Name)
	d.Set("email", backupUnit.Properties.Email)
	d.Set("login", fmt.Sprintf("%s-%d", backupUnit.Properties.Name, int64(contractResources.Properties.PBContractNumber)))

	log.Printf("[INFO] Successfully retreived backup unit %s: %+v", d.Id(), backupUnit)

	return nil
}

func resourceBackupUnitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)
	request := profitbricks.BackupUnit{}

	request.Properties = &profitbricks.BackupUnitProperties{}

	log.Printf("[INFO] Attempting update backup unit %s", d.Id())

	if d.HasChange("email") {
		oldEmail, newEmail := d.GetChange("email")
		log.Printf("[INFO] backup unit email changed from %+v to %+v", oldEmail, newEmail)
		request.Properties.Email = newEmail.(string)
	}

	if d.HasChange("password") {
		oldPassword, newPassword := d.GetChange("password")
		log.Printf("[INFO] backup unit password changed from %+v to %+v", oldPassword, newPassword)
		request.Properties.Password = newPassword.(string)
	}

	_, err := client.UpdateBackupUnit(d.Id(), request)

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while updating backup unit: %s", err)
		}
		return fmt.Errorf("Error while updating backup unit %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		backupUnitReady, rsErr := backupUnitReady(client, d)

		if rsErr != nil {
			return fmt.Errorf("Error while checking readiness status of backup unit %s: %s", d.Id(), rsErr)
		}

		if backupUnitReady && rsErr == nil {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}
	}

	return resourceBackupUnitRead(d, meta)
}

func resourceBackupUnitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*profitbricks.Client)

	_, err := client.DeleteBackupUnit(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error while deleting backup unit: %s", err)
		}

		return fmt.Errorf("Error while deleting backup unit %s: %s", d.Id(), err)
	}

	for {
		log.Printf("[INFO] Waiting for backupUnit %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		backupUnitDeleted, dsErr := backupUnitDeleted(client, d)

		if dsErr != nil {
			return fmt.Errorf("Error while checking deletion status of backup unit %s: %s", d.Id(), dsErr)
		}

		if backupUnitDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted backup unit: %s", d.Id())
			break
		}
	}

	return nil
}

func backupUnitReady(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	subjectBackupUnit, err := client.GetBackupUnit(d.Id())

	if err != nil {
		return true, fmt.Errorf("Error checking backup unit status: %s", err)
	}
	return subjectBackupUnit.Metadata.State == "AVAILABLE", nil
}

func backupUnitDeleted(client *profitbricks.Client, d *schema.ResourceData) (bool, error) {
	_, err := client.GetBackupUnit(d.Id())

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				return true, nil
			}
			return true, fmt.Errorf("Error checking backup unit deletion status: %s", err)
		}
	}
	return false, nil
}
