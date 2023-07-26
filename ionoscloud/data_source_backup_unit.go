package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

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
	client := meta.(services.SdkBundle).CloudApiClient

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

	if idOk {
		/* search by ID */
		backupUnit, apiResponse, err = client.BackupUnitsApi.BackupunitsFindById(ctx, id.(string)).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the backup unit %s: %w", id.(string), err))
		}
		if backupUnit.Properties != nil {
			log.Printf("[INFO] Got backupUnit [Name=%s] [Id=%s]", *backupUnit.Properties.Name, *backupUnit.Id)
		}
	} else {
		/* search by name */
		var backupUnits ionoscloud.BackupUnits

		backupUnits, apiResponse, err := client.BackupUnitsApi.BackupunitsGet(ctx).Depth(1).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %w", err))
		}

		var results []ionoscloud.BackupUnit
		if backupUnits.Items != nil {
			for _, bu := range *backupUnits.Items {
				if bu.Properties != nil && bu.Properties.Name != nil && *bu.Properties.Name == name.(string) {
					tmpBackupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, *bu.Id).Execute()
					logApiRequestTime(apiResponse)
					if err != nil {
						return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit with ID %s: %w", *bu.Id, err))
					}
					results = append(results, tmpBackupUnit)
				}

			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no backup unit found with the specified name %s", name))
		} else {
			backupUnit = results[0]
		}

	}

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if cErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while fetching contract resources for backup unit %s: %w", d.Id(), cErr))
		return diags
	}

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
