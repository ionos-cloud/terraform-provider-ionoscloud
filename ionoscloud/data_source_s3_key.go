package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceObjectStorageKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectStorageKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Id of the key.",
				Optional:    true,
			},
			"user_id": {
				Type:             schema.TypeString,
				Description:      "The ID of the user that owns the key.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The Secret key.",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether this key should be active or not.",
				Optional:    true,
				Default:     true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceObjectStorageKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id, idOk := d.GetOk("id")

	if !idOk {
		return diag.FromErr(fmt.Errorf("please provide the object storage key id"))
	}
	d.SetId(id.(string))

	if diags := resourceS3KeyRead(ctx, d, meta); diags != nil {
		return diags
	}

	return nil
}
