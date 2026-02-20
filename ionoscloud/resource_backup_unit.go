package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
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
				Type:             schema.TypeString,
				Description:      "Alphanumeric name you want assigned to the backup unit.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:             schema.TypeString,
				Description:      "The password you want assigned to the backup unit.",
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"email": {
				Type:             schema.TypeString,
				Description:      "The e-mail address you want assigned to the backup unit.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
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
	client := meta.(bundleclient.SdkBundle).CloudApiClient

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
		return utils.ToDiags(d, fmt.Sprintf("error creating backup unit: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	d.SetId(*createdBackupUnit.Id)
	log.Printf("[INFO] Created backup unit: %s", d.Id())

	if diags := waitForUnitToBeReady(ctx, d, client); diags != nil {
		return diags
	}

	return resourceBackupUnitRead(ctx, d, meta)
}

func resourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(bundleclient.SdkBundle).CloudApiClient

	backupUnit, apiResponse, err := BackupUnitFindByID(ctx, d.Id(), client)
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching backup unit: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	contractResources, contractApiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(contractApiResponse)

	if cErr != nil {
		return utils.ToDiags(d, fmt.Sprintf("error while fetching contract resources for backup unit: %s", cErr), &utils.DiagsOpts{StatusCode: contractApiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retrieved contract resource for backup unit unit %s: %+v", d.Id(), contractResources)

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	log.Printf("[INFO] Successfully retrieved backup unit %s: %+v", d.Id(), backupUnit)

	return nil
}

func resourceBackupUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

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
		return utils.ToDiags(d, fmt.Sprintf("error while updating backup unit: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
			return utils.ToDiags(d, fmt.Sprintf("error while checking readiness status of backup unit: %s", rsErr), nil)
		}

		if backupUnitReady {
			log.Printf("[INFO] backup unit ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			return utils.ToDiags(d, fmt.Sprintf("backup unit readiness check timed out! WARNING: your backup unit will still probably be created/updated "+
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}
	}
	return nil
}

func resourceBackupUnitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	apiResponse, err := client.BackupUnitsApi.BackupunitsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting backup unit: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	for {
		log.Printf("[INFO] Waiting for backupUnit %s to be deleted...", d.Id())

		backupUnitDeleted, dsErr := backupUnitDeleted(client, d, ctx)

		if dsErr != nil {
			return utils.ToDiags(d, fmt.Sprintf("error while checking deletion status of backup unit: %s", dsErr), nil)
		}

		if backupUnitDeleted {
			log.Printf("[INFO] Successfully deleted backup unit: %s", d.Id())
			break
		}

		select {
		case <-time.After(constant.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			return utils.ToDiags(d, fmt.Sprintf("backup unit deletion timed out! WARNING: your backup unit will still probably be deleted "+
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"), nil)
		}
	}

	return nil
}

func resourceBackupUnitImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).CloudApiClient

	buId := d.Id()

	backupUnit, apiResponse, err := BackupUnitFindByID(ctx, d.Id(), client)
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("unable to find Backup Unit %q", buId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("unable to retrieve Backup Unit: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Backup Unit found: %+v", backupUnit)

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if cErr != nil {
		return nil, utils.ToError(d, fmt.Sprintf("error while fetching contract resources for backup unit: %s", cErr), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}

func backupUnitReady(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	backupUnit, apiResponse, err := BackupUnitFindByID(c, d.Id(), client)
	logApiRequestTime(apiResponse)

	if err != nil {
		return true, fmt.Errorf("error checking backup unit status: %w", err)
	}
	return strings.EqualFold(*backupUnit.Metadata.State, constant.Available), nil
}

func backupUnitDeleted(client *ionoscloud.APIClient, d *schema.ResourceData, c context.Context) (bool, error) {
	_, apiResponse, err := BackupUnitFindByID(c, d.Id(), client)
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

		if backupUnit.Properties.Name != nil && contractResources.Items != nil && len(*contractResources.Items) > 0 {
			firstContract := (*contractResources.Items)[0]
			if firstContract.Properties != nil && firstContract.Properties.ContractNumber != nil {
				err := d.Set("login", fmt.Sprintf("%d-%s", *firstContract.Properties.ContractNumber, *backupUnit.Properties.Name))
				if err != nil {
					return fmt.Errorf("error while setting login property for backup unit %s: %w", d.Id(), err)
				}
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

// BackupUnitFindByID simulates a FindByID function by filtering backup units from BackupunitsGet using the given ID.
// This is done because of a temporary bug in the API with the regular FindByID function.
// This is a temporary fix, this function should be replaced after the API bug is fixed.
func BackupUnitFindByID(ctx context.Context, backupUnitID string, client *ionoscloud.APIClient) (ionoscloud.BackupUnit, *ionoscloud.APIResponse, error) {
	backupUnits, apiResponse, err := client.BackupUnitsApi.BackupunitsGet(ctx).Depth(2).Execute()
	var backupUnit ionoscloud.BackupUnit
	logApiRequestTime(apiResponse)
	if err != nil {
		return backupUnit, apiResponse, fmt.Errorf("error while retrieving the list of backup units: %w", err)
	}
	if backupUnits.Items == nil {
		return backupUnit, apiResponse, fmt.Errorf("expected a list of backup units in the response but received 'nil' instead")
	}
	for _, backupUnit := range *backupUnits.Items {
		if backupUnit.Id == nil {
			return backupUnit, apiResponse, fmt.Errorf("expected a backup unit with a valid ID but received 'nil' instead")
		}
		if *backupUnit.Id == backupUnitID {
			return backupUnit, apiResponse, nil
		}
	}
	// If the backup unit wasn't found, simulate a 404 error.
	apiResponse.StatusCode = http.StatusNotFound
	return backupUnit, apiResponse, fmt.Errorf("backup unit with ID %s not found", backupUnitID)
}
