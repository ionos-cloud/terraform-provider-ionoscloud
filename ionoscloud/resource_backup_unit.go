package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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

	createdBackupUnit, _, err := client.BackupUnitsApi.BackupunitsPost(ctx).BackupUnit(backupUnit).Execute()

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating backup unit: %s", err))
		return diags
	}

	d.SetId(*createdBackupUnit.Id)
	log.Printf("[INFO] Created backup unit: %s", d.Id())

	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())

		backupUnitReady, rsErr := backupUnitReady(client, d, ctx)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of backup unit %s: %s", d.Id(), rsErr))
			return diags
		}

		if backupUnitReady {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("backup unit creation timed out! WARNING: your backup unit will still probably be created after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func resourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ionoscloud.APIClient)

	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching backup unit %s: %s", d.Id(), err))
		return diags
	}

	contractResources, _, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()

	if cErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while fetching contract resources for backup unit %s: %s", d.Id(), cErr))
		return diags
	}

	log.Printf("[INFO] Successfully retreived contract resource for backup unit unit %s: %+v", d.Id(), contractResources)

	setBackupUnitData(d, &backupUnit, &contractResources)

	log.Printf("[INFO] Successfully retreived backup unit %s: %+v", d.Id(), backupUnit)

	return nil
}

func resourceBackupUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

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

	_, apiResponse, err := client.BackupUnitsApi.BackupunitsPut(ctx, d.Id()).BackupUnit(request).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating backup unit %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())

		backupUnitReady, rsErr := backupUnitReady(client, d, ctx)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of backup unit %s: %s", d.Id(), rsErr))
			return diags
		}

		if backupUnitReady {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("backup unit update timed out! WARNING: your backup unit will still probably be updated after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func resourceBackupUnitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	apiResponse, err := client.BackupUnitsApi.BackupunitsDelete(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting backup unit %s: %s", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for backupUnit %s to be deleted...", d.Id())

		backupUnitDeleted, dsErr := backupUnitDeleted(client, d, ctx)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of backup unit %s: %s", d.Id(), dsErr))
			return diags
		}

		if backupUnitDeleted {
			log.Printf("[INFO] Successfully deleted backup unit: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("backup unit deletion timed out! WARNING: your backup unit will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return nil
}

func resourceBackupUnitImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ionoscloud.APIClient)

	buId := d.Id()

	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Backup Unit %q", buId)
		}
		return nil, fmt.Errorf("unable to retreive Backup Unit %q", buId)
	}

	log.Printf("[INFO] Backup Unit found: %+v", backupUnit)

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()

	if cErr != nil {
		return nil, fmt.Errorf("error while fetching contract resources for backup unit %q: %s", d.Id(), cErr)
	}

	if contractResources.Items == nil || len(*contractResources.Items) == 0 {
		return nil, fmt.Errorf("no contracts found for user")
	}

	props := (*contractResources.Items)[0].Properties
	if props == nil {
		return nil, fmt.Errorf("could not get first contract properties")
	}

	if props.ContractNumber == nil {
		return nil, fmt.Errorf("contract number not set")
	}

	setBackupUnitData(d, &backupUnit, &contractResources)

	return []*schema.ResourceData{d}, nil
}

func backupUnitReady(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	backupUnit, _, err := client.BackupUnitsApi.BackupunitsFindById(c, d.Id()).Execute()

	if err != nil {
		return true, fmt.Errorf("error checking backup unit status: %s", err)
	}
	return *backupUnit.Metadata.State == "AVAILABLE", nil
}

func backupUnitDeleted(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	_, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(c, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking backup unit deletion status: %s", err)
	}
	return false, nil
}
