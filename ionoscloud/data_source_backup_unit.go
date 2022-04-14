package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"strings"
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
				Description: "The login associated with the backup unit. Derived from the contract number.",
				Computed:    true,
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using name argument.",
				Default:     false,
				Optional:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")

	id := idValue.(string)
	name := nameValue.(string)

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
		backupUnit, apiResponse, err = client.BackupUnitsApi.BackupunitsFindById(ctx, id).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the backup unit while searching by id %s: %w", id, err))
		}
		if backupUnit.Properties != nil {
			log.Printf("[INFO] Got backupUnit [Name=%s] [Id=%s]", *backupUnit.Properties.Name, *backupUnit.Id)
		}
	} else {
		/* search by name */
		var results []ionoscloud.BackupUnit

		partialMatch := d.Get("partial_match").(bool)

		if partialMatch {
			backupUnits, apiResponse, err := client.BackupUnitsApi.BackupunitsGet(ctx).Depth(1).Filter("name", name).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching backup units while searching by partial name: %s, %w", name, err))
			}

			results = *backupUnits.Items
		} else {
			backupUnits, apiResponse, err := client.BackupUnitsApi.BackupunitsGet(ctx).Depth(1).Execute()
			logApiRequestTime(apiResponse)

			if err != nil {
				return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
			}

			if backupUnits.Items != nil {
				for _, bu := range *backupUnits.Items {
					if bu.Properties != nil && bu.Properties.Name != nil && strings.EqualFold(*bu.Properties.Name, name) {
						tmpBackupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, *bu.Id).Execute()
						logApiRequestTime(apiResponse)
						if err != nil {
							return diag.FromErr(fmt.Errorf("an error occurred while fetching the backup unit with ID: %s while searching by full name: %s: %w", *bu.Id, name, err))
						}
						results = append(results, tmpBackupUnit)
					}

				}
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no backup unit found with the specified name %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one backup unit found with the specified criteria: name = %s", name))
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
