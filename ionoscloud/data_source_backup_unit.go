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

	if idOk {
		/* search by ID */
		backupUnit, _, err = client.BackupUnitsApi.BackupunitsFindById(ctx, id.(string)).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the backup unit %s: %s", id.(string), err))
		}
	} else {
		/* search by name */
		var backupUnits ionoscloud.BackupUnits

		backupUnits, _, err := client.BackupUnitsApi.BackupunitsGet(ctx).Execute()

		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit: %s", err.Error()))
		}

		if backupUnits.Items != nil {
			for _, bu := range *backupUnits.Items {
				tmpBackupUnit, _, err := client.BackupUnitsApi.BackupunitsFindById(ctx, *bu.Id).Execute()
				if err != nil {
					return diag.FromErr(fmt.Errorf("an error occurred while fetching backup unit with ID %s: %s", *bu.Id, err.Error()))
				}
				if tmpBackupUnit.Properties.Name != nil && *tmpBackupUnit.Properties.Name == name.(string) {
					backupUnit = tmpBackupUnit
					break
				}

			}
		}

	}

	if &backupUnit == nil {
		return diag.FromErr(fmt.Errorf("backup unit not found"))
	}

	contractResources, _, cErr := client.ContractResourcesApi.ContractsGet(ctx).Execute()

	if cErr != nil {
		diags := diag.FromErr(fmt.Errorf("error while fetching contract resources for backup unit %s: %s", d.Id(), cErr))
		return diags
	}

	if err := d.Set("id", *backupUnit.Id); err != nil {
		return diag.FromErr(err)
	}

	if diags := setBackupUnitData(d, &backupUnit, &contractResources); diags != nil {
		return diags
	}

	return nil
}

func setBackupUnitData(d *schema.ResourceData, backupUnit *ionoscloud.BackupUnit, contractResources *ionoscloud.Contracts) diag.Diagnostics {

	if backupUnit.Id != nil {
		d.SetId(*backupUnit.Id)
	}

	if backupUnit.Properties != nil {

		if backupUnit.Properties.Name != nil {
			epErr := d.Set("name", backupUnit.Properties.Name)
			if epErr != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting name property for backup unit %s: %s", d.Id(), epErr))
				return diags
			}
		}

		if backupUnit.Properties.Email != nil {
			epErr := d.Set("email", backupUnit.Properties.Email)
			if epErr != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting email property for backup unit %s: %s", d.Id(), epErr))
				return diags
			}
		}

		if backupUnit.Properties.Name != nil && contractResources.Items != nil && len(*contractResources.Items) > 0 &&
			(*contractResources.Items)[0].Properties.ContractNumber != nil {
			err := d.Set("login", fmt.Sprintf("%s-%d", *backupUnit.Properties.Name, *(*contractResources.Items)[0].Properties.ContractNumber))
			if err != nil {
				diags := diag.FromErr(fmt.Errorf("error while setting login property for backup unit %s: %s", d.Id(), err))
				return diags
			}
		}
	}
	return nil
}
