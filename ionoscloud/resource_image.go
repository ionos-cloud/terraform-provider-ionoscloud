package ionoscloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/nfs"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: resourceImageImport,
		},
		Schema: map[string]*schema.Schema{
			"location": {
				Type: schema.TypeString,
				Description: fmt.Sprintf(
					"The location of the Network File Storage Cluster. "+
						"Available locations: '%s'", strings.Join(nfs.ValidNFSLocations, ", '"),
				),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(nfs.ValidNFSLocations, false)),
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network File Storage Cluster.",
				Computed:    true,
			},
