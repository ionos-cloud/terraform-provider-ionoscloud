package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBackupUnit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupUnitCreate,
		ReadContext:   resourceBackupUnitRead,
		UpdateContext: resourceBackupUnitUpdate,
		DeleteContext: resourceBackupUnitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBackupUnitImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "Alphanumeric name you want assigned to the backup unit.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:         schema.TypeString,
				Description:  "The password you want assigned to the backup unit.",
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"email": {
				Type:         schema.TypeString,
				Description:  "The e-mail address you want assigned to the backup unit.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
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

func resourceBackupUnitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	backupUnitName := d.Get("name").(string)
	backupUnitPassword := d.Get("password").(string)
	backupUnitEmail := d.Get("email").(string)

	backupUnit := ionoscloud.BackupUnit{
		Properties: &ionoscloud.BackupUnitProperties{
			Name:     &backupUnitName,
			Password: &backupUnitPassword,
			Email:    &backupUnitEmail,
		},
	}

	createdBackupUnit, apiResponse, err := client.BackupUnitApi.BackupunitsPost(ctx).BackupUnit(backupUnit).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("Error creating backup unit: %s, ApiError: %s", err, string(apiResponse.Payload)))
		return diags
	}

	d.SetId(*createdBackupUnit.Id)
	log.Printf("[INFO] Created backup unit: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		backupUnitReady, rsErr := backupUnitReady(client, d, ctx)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("Error while checking readiness status of backup unit %s: %s", d.Id(), rsErr))
			return diags
		}

		if backupUnitReady && rsErr == nil {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func resourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	backupUnit, apiResponse, err := client.BackupUnitApi.BackupunitsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("Error while fetching backup unit %s: %s", d.Id(), err))
		return diags
	}

	contractResources, _, cErr := client.ContractApi.ContractsGet(ctx).Execute()

	if cErr != nil {
		diags := diag.FromErr(fmt.Errorf("Error while fetching contract resources for backup unit %s: %s", d.Id(), cErr))
		return diags
	}

	log.Printf("[INFO] Successfully retreived contract resource for backup unit unit %s: %+v", d.Id(), contractResources)

	if backupUnit.Properties.Name != nil {
		err := d.Set("name", *backupUnit.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("Error while setting name property for backup unit %s: %s", d.Id(), err))
			return diags
		}
	}

	if backupUnit.Properties.Email != nil {
		epErr := d.Set("email", backupUnit.Properties.Email)
		if epErr != nil {
			diags := diag.FromErr(fmt.Errorf("Error while setting email property for backup unit %s: %s", d.Id(), epErr))
			return diags
		}
	}

	if backupUnit.Properties.Name != nil && contractResources.Properties.ContractNumber != nil {
		err := d.Set("login", fmt.Sprintf("%s-%d", *backupUnit.Properties.Name, *contractResources.Properties.ContractNumber))
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("Error while setting login property for backup unit %s: %s", d.Id(), err))
			return diags
		}
	}

	log.Printf("[INFO] Successfully retreived backup unit %s: %+v", d.Id(), backupUnit)

	return nil
}

func resourceBackupUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)
	var diags diag.Diagnostics

	request := ionoscloud.BackupUnit{}
	request.Properties = &ionoscloud.BackupUnitProperties{}

	log.Printf("[INFO] Attempting update backup unit %s", d.Id())

	if d.HasChange("email") {
		oldEmail, newEmail := d.GetChange("email")
		log.Printf("[INFO] backup unit email changed from %+v to %+v", oldEmail, newEmail)

		newEmailStr := newEmail.(string)
		request.Properties.Email = &newEmailStr
	}

	if d.HasChange("password") {
		oldPassword, newPassword := d.GetChange("password")
		log.Printf("[INFO] backup unit password changed from %+v to %+v", oldPassword, newPassword)

		newPasswordStr := newPassword.(string)
		request.Properties.Password = &newPasswordStr

		if !d.HasChange("email") {
			oldEmail, _ := d.GetChange("email")
			oldEmailStr := oldEmail.(string)
			request.Properties.Email = &oldEmailStr
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.BackupUnitApi.BackupunitsPut(ctx, d.Id()).BackupUnit(request).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error while updating backup unit %s: %s", d.Id(), err))
	}

	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())
		time.Sleep(5 * time.Second)

		backupUnitReady, rsErr := backupUnitReady(client, d, ctx)

		if rsErr != nil {
			diags = diag.FromErr(fmt.Errorf("Error while checking readiness status of backup unit %s: %s", d.Id(), rsErr))
			return diags
		}

		if backupUnitReady && rsErr == nil {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func resourceBackupUnitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}
	_, apiResponse, err := client.BackupUnitApi.BackupunitsDelete(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("Error while deleting backup unit %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for backupUnit %s to be deleted...", d.Id())
		time.Sleep(5 * time.Second)

		backupUnitDeleted, dsErr := backupUnitDeleted(client, d, ctx)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("Error while checking deletion status of backup unit %s: %s", d.Id(), dsErr))
			return diags
		}

		if backupUnitDeleted && dsErr == nil {
			log.Printf("[INFO] Successfully deleted backup unit: %s", d.Id())
			break
		}
	}

	return nil
}

func backupUnitReady(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	backupUnit, _, err := client.BackupUnitApi.BackupunitsFindById(c, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("Error checking backup unit status: %s", err)
	}
	return *backupUnit.Metadata.State == "AVAILABLE", nil
}

func backupUnitDeleted(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	_, apiResponse, err := client.BackupUnitApi.BackupunitsFindById(c, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking backup unit deletion status: %s", err)
	}
	return false, nil
}
