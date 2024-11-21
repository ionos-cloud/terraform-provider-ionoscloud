package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/nsg"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"
)

func resourceDatacenterNSGSelection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatacenterNSGSelectionCreate,
		ReadContext:   resourceDatacenterNSGSelectionRead,
		UpdateContext: resourceDatacenterNSGSelectionUpdate,
		DeleteContext: resourceDatacenterNSGSelectionDelete,
		Schema: map[string]*schema.Schema{

			"datacenter_id": {
				Type:             schema.TypeString,
				Description:      "ID of the Datacenter to which the NSG will be attached.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"nsg_id": {
				Type:             schema.TypeString,
				Description:      "ID of the NSG which will be attached to the datacenter. If an empty string is specified and a NSG was attached previously, it will be unset.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.Any(validation.IsUUID, validation.StringIsEmpty)),
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}
func resourceDatacenterNSGSelectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	nsgID := d.Get("nsg_id").(string)

	ns := nsg.Service{Client: meta.(services.SdkBundle).CloudAPIClient, Meta: meta, D: d}
	if diags := ns.SetDefaultDatacenterNSG(ctx, dcID, nsgID); diags.HasError() {
		return diags
	}

	d.SetId(uuidgen.ResourceUuid().String())
	return resourceDatacenterNSGSelectionRead(ctx, d, meta)
}

func resourceDatacenterNSGSelectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudAPIClient
	dcID := d.Get("datacenter_id").(string)
	nsgID := d.Get("nsg_id").(string)

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcID).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = setDatacenterDefaultNSGSelectionData(d, &datacenter); err != nil {
		return diag.FromErr(fmt.Errorf("error reading default NSG for datacenter, dcId: %s, sId: %s, (%w)", dcID, nsgID, err))
	}
	return nil
}

func resourceDatacenterNSGSelectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)

	if d.HasChange("nsg_id") {
		_, newID := d.GetChange("nsg_id")
		ns := nsg.Service{Client: meta.(services.SdkBundle).CloudAPIClient, Meta: meta, D: d}
		if diags := ns.SetDefaultDatacenterNSG(ctx, dcID, newID.(string)); diags.HasError() {
			return diags
		}
	}

	return resourceDatacenterNSGSelectionRead(ctx, d, meta)

}

func resourceDatacenterNSGSelectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcID := d.Get("datacenter_id").(string)
	ns := nsg.Service{Client: meta.(services.SdkBundle).CloudAPIClient, Meta: meta, D: d}
	if diags := ns.SetDefaultDatacenterNSG(ctx, dcID, ""); diags.HasError() {
		return diags
	}
	d.SetId("")

	return nil
}

func setDatacenterDefaultNSGSelectionData(d *schema.ResourceData, datacenter *ionoscloud.Datacenter) error {

	if datacenter.Properties.DefaultSecurityGroupId != nil {
		if err := d.Set("nsg_id", *datacenter.Properties.DefaultSecurityGroupId); err != nil {
			return err
		}
		return nil
	}
	return nil
}
