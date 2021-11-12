package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func dataSourceBackupUnit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupUnitRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Id of the backup unit.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Alphanumeric name you want assigned to the backup unit.",
				Optional:    true,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "The e-mail address you want assigned to the backup unit.",
				Computed:    true,
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

func dataSourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ionoscloud.APIClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("id and name cannot be both specified in the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the backup unit id or name"))
	}
	var backupUnit ionoscloud.BackupUnit
	var err error
	var apiResponse *ionoscloud.APIResponse

	found := false

	if idOk {
		/* search by ID */
		backupUnit, apiResponse, err = client.BackupUnitsApi.BackupunitsFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the backup unit %s: %s", id.(string), err))
		}
		found = true
	} else {
		/* search by name */
		var backupUnits ionoscloud.BackupUnits

		backupUnits, apiResponse, err := client.BackupUnitsApi.BackupunitsGet(ctx).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
		}

		if backupUnits.Items != nil {
			for _, bu := range *backupUnits.Items {
				tmpBackupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, *bu.Id).Execute()
				logApiRequestTime(apiResponse)
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit with ID %s: %s", *bu.Id, err.Error()))
				}
				if tmpBackupUnit.Properties.Name != nil && *tmpBackupUnit.Properties.Name == name.(string) {
					backupUnit = tmpBackupUnit
					found = true
					break
				}

			}
		}

	}

	if !found {
		return diag.FromErr(fmt.Errorf("backup unit not found"))
	}

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if cErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while fetching contract resources for backup unit %s: %s", d.Id(), cErr))
		return diags
	}

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
