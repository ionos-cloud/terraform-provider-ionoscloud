package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/nsg"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/uuidgen"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceDatacenterDefaultNSGSelection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatacenterDefaultNSGSelectionCreate,
		ReadContext:   resourceDatacenterDefaultNSGSelectionRead,
		UpdateContext: resourceDatacenterDefaultNSGSelectionUpdate,
		DeleteContext: resourceDatacenterDefaultNSGSelectionDelete,
		Schema: map[string]*schema.Schema{

			"datacenter_id": {
				Type:             schema.TypeString,
				Description:      "ID of the Datacenter for which the default NSG will be set.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"nsg_id": {
				Type:             schema.TypeString,
				Description:      "ID of the NSG which will be set as default for the datacenter. If empty string is specified, any the default NSG will be unset",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.Any(validation.IsUUID, validation.StringIsEmpty)),
			},
		},

		Timeouts: &resourceDefaultTimeouts,
	}
}
func resourceDatacenterDefaultNSGSelectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)
	nsgId := d.Get("nsg_id").(string)

	ns := nsg.Service{Client: meta.(services.SdkBundle).CloudApiClient, Meta: meta, D: d}
	if diags := ns.SetDefaultDatacenterNSG(ctx, dcId, nsgId); diags.HasError() {
		return diags
	}

	d.SetId(uuidgen.ResourceUuid().String())
	return resourceDatacenterDefaultNSGSelectionRead(ctx, d, meta)
}

func resourceDatacenterDefaultNSGSelectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CloudApiClient
	dcId := d.Get("datacenter_id").(string)
	nsgId := d.Get("nsg_id").(string)

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()
	apiResponse.LogInfo()
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err = setDatacenterDefaultNSGSelectionData(d, &datacenter); err != nil {
		return diag.FromErr(fmt.Errorf("error reading default NSG for datacenter, dcId: %s, sId: %s, (%w)", dcId, nsgId, err))
	}
	return nil
}

func resourceDatacenterDefaultNSGSelectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)

	if d.HasChange("nsg_id") {
		_, newId := d.GetChange("nsg_id")
		ns := nsg.Service{Client: meta.(services.SdkBundle).CloudApiClient, Meta: meta, D: d}
		if diags := ns.SetDefaultDatacenterNSG(ctx, dcId, newId.(string)); diags.HasError() {
			return diags
		}
	}

	return resourceDatacenterDefaultNSGSelectionRead(ctx, d, meta)

}

func resourceDatacenterDefaultNSGSelectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dcId := d.Get("datacenter_id").(string)
	ns := nsg.Service{Client: meta.(services.SdkBundle).CloudApiClient, Meta: meta, D: d}
	if diags := ns.SetDefaultDatacenterNSG(ctx, dcId, ""); diags.HasError() {
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
