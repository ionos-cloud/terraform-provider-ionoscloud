package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func dataSourceBackupUnit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupUnitRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Id of the backup unit.",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Alphanumeric name you want assigned to the backup unit.",
				Optional:    true,
				Computed:    true,
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

func dataSourceBackupUnitRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClientWithFailover()
	if err != nil {
		return diag.FromErr(err)
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("id and name cannot be both specified in the same time"), nil)
	}
	if !idOk && !nameOk {
		return diagutil.ToDiags(d, fmt.Errorf("please provide either the backup unit id or name"), nil)
	}
	var backupUnit ionoscloud.BackupUnit
	var apiResponse *ionoscloud.APIResponse

	if idOk {
		/* search by ID */
		backupUnit, apiResponse, err = BackupUnitFindByID(ctx, id.(string), client)
		logApiRequestTime(apiResponse)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching the backup unit %s: %w", id.(string), err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
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
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching backup unit: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}

		var results []ionoscloud.BackupUnit
		if backupUnits.Items != nil {
			for _, bu := range *backupUnits.Items {
				if bu.Properties != nil && bu.Properties.Name != nil && *bu.Properties.Name == name.(string) {
					tmpBackupUnit, apiResponse, err := BackupUnitFindByID(ctx, *bu.Id, client)
					logApiRequestTime(apiResponse)
					if err != nil {
						return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching backup unit with ID %s: %w", *bu.Id, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
					}
					results = append(results, tmpBackupUnit)
				}

			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no backup unit found with the specified name %s", name), nil)
		} else {
			backupUnit = results[0]
		}

	}

	contractResources, apiResponse, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()
	logApiRequestTime(apiResponse)

	if cErr != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching contract resources for backup unit: %w", cErr), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	if err := setBackupUnitData(d, &backupUnit, &contractResources); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
