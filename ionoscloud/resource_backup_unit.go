package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"strings"
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
				Type:             schema.TypeString,
				Description:      "The e-mail address you want assigned to the backup unit.",
				Required:         true,
				ValidateFunc:     validation.All(validation.StringIsNotWhiteSpace),
				DiffSuppressFunc: utils.DiffToLower,
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
	client := meta.(SdkBundle).CloudApiClient

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

	createdBackupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsPost(ctx).BackupUnit(backupUnit).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating backup unit: %w", err))
		return diags
	}

	d.SetId(*createdBackupUnit.Id)
	log.Printf("[INFO] Created backup unit: %s", d.Id())

	if diags := waitForUnitToBeReady(ctx, d, client); diags != nil {
		return diags
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func resourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching backup unit %s: %w", d.Id(), err))
		return diags
	}

	contractResources, _, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if cErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while fetching contract resources for backup unit %s: %w", d.Id(), cErr))
		return diags
	}

	log.Printf("[INFO] Successfully retreived contract resource for backup unit unit %s: %+v", d.Id(), contractResources)

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retreived backup unit %s: %+v", d.Id(), backupUnit)

	return nil
}

func resourceBackupUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	request := ionoscloud.BackupUnit{}
	request.Properties = &ionoscloud.BackupUnitProperties{}

	log.Printf("[INFO] Attempting update backup unit %s", d.Id())
	oldEmail, newEmail := d.GetChange("email")
	emailStr := oldEmail.(string)
	if d.HasChange("email") {
		log.Printf("[INFO] backup unit email changed from %+v to %+v", oldEmail, newEmail)
		emailStr = newEmail.(string)
	}
	request.Properties.Email = &emailStr

	if d.HasChange("password") {
		_, newPassword := d.GetChange("password")
		log.Printf("[INFO] backup unit password changed")

		newPasswordStr := newPassword.(string)
		request.Properties.Password = &newPasswordStr

		if !d.HasChange("email") {
			oldEmail, _ := d.GetChange("email")
			oldEmailStr := oldEmail.(string)
			request.Properties.Email = &oldEmailStr
		}
	}

	_, apiResponse, err := client.BackupUnitsApi.BackupunitsPut(ctx, d.Id()).BackupUnit(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while updating backup unit %s: %w", d.Id(), err))
		return diags
	}

	if diags := waitForUnitToBeReady(ctx, d, client); diags != nil {
		return diags
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func waitForUnitToBeReady(ctx context.Context, d *schema.ResourceData, client *ionoscloud.APIClient) diag.Diagnostics {
	for {
		log.Printf("[INFO] Waiting for backup unit %s to be ready...", d.Id())

		backupUnitReady, rsErr := backupUnitReady(client, d, ctx)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of backup unit %s: %w", d.Id(), rsErr))
			return diags
		}

		if backupUnitReady {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("backup unit readiness check timed out! WARNING: your backup unit will still probably be created/updated " +
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}
	return nil
}

func resourceBackupUnitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.BackupUnitsApi.BackupunitsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting backup unit %s: %w", d.Id(), err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for backupUnit %s to be deleted...", d.Id())

		backupUnitDeleted, dsErr := backupUnitDeleted(client, d, ctx)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of backup unit %s: %w", d.Id(), dsErr))
			return diags
		}

		if backupUnitDeleted {
			log.Printf("[INFO] Successfully deleted backup unit: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("backup unit deletion timed out! WARNING: your backup unit will still probably be deleted " +
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}

	return nil
}

func resourceBackupUnitImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

	buId := d.Id()

	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find Backup Unit %q", buId)
		}
		return nil, fmt.Errorf("unable to retreive Backup Unit %q", buId)
	}

	log.Printf("[INFO] Backup Unit found: %+v", backupUnit)

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if cErr != nil {
		return nil, fmt.Errorf("error while fetching contract resources for backup unit %q: %s", d.Id(), cErr)
	}

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func backupUnitReady(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	backupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(c, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking backup unit status: %w", err)
	}
	return strings.EqualFold(*backupUnit.Metadata.State, utils.Available), nil
}

func backupUnitDeleted(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	_, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(c, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking backup unit deletion status: %w", err)
	}
	return false, nil
}

func setBackupUnitData(d *schema.ResourceData, backupUnit *ionoscloud.BackupUnit, contractResources *ionoscloud.Contracts) error {

	if backupUnit.Id != nil {
		d.SetId(*backupUnit.Id)
	}

	if backupUnit.Properties != nil {

		if backupUnit.Properties.Name != nil {
			epErr := d.Set("name", *backupUnit.Properties.Name)
			if epErr != nil {
				return fmt.Errorf("error while setting name property for backup unit %s: %w", d.Id(), epErr)
			}
		}

		if backupUnit.Properties.Email != nil {
			epErr := d.Set("email", *backupUnit.Properties.Email)
			if epErr != nil {
				return fmt.Errorf("error while setting email property for backup unit %s: %w", d.Id(), epErr)
			}
		}

		if backupUnit.Properties.Name != nil && contractResources.Items != nil && len(*contractResources.Items) > 0 &&
			(*contractResources.Items)[0].Properties.ContractNumber != nil {
			err := d.Set("login", fmt.Sprintf("%s-%d", *backupUnit.Properties.Name, *(*contractResources.Items)[0].Properties.ContractNumber))
			if err != nil {
				return fmt.Errorf("error while setting login property for backup unit %s: %w", d.Id(), err)
			}
		} else {
			if contractResources.Items == nil || len(*contractResources.Items) == 0 {
				return fmt.Errorf("no contracts found for user")
			}

			props := (*contractResources.Items)[0].Properties
			if props == nil {
				return fmt.Errorf("could not get first contract properties")
			}

			if props.ContractNumber == nil {
				return fmt.Errorf("contract number not set")
			}
		}
	}
	return nil
}
