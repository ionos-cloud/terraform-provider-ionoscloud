package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nfs "github.com/ionos-cloud/sdk-go-nfs"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func dataSourceNFSShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNFSShareRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network File Storage Share.",
				Optional:    true,
				Computed:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network File Storage Cluster.",
				Required:    true,
			},
			"location": {
				Type: schema.TypeString,
				Description: fmt.Sprintf("The location of the Network File Storage Cluster. "+
					"Available locations: '%s'", strings.Join(ValidNFSLocations, ", '")),
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ValidNFSLocations, false)),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Network File Storage Share",
				Optional:    true,
				// Computed:    true,
			},
			"nfs_path": {
				Type:        schema.TypeString,
				Description: "Path to the NFS export. The NFS path is the path to the directory being exported.",
				Computed:    true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "The quota in MiB for the export. The quota can restrict the amount of data that can be stored within the export. The quota can be disabled using `0`.",
				Optional:    true,
				Computed:    true,
			},
			"gid": {
				Type:        schema.TypeInt,
				Description: "The group ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.",
				Optional:    true,
				Computed:    true,
			},
			"uid": {
				Type:        schema.TypeInt,
				Description: "The user ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.",
				Optional:    true,
				Computed:    true,
			},
			"client_groups": {
				Type:        schema.TypeList,
				Description: "The groups of clients are the systems connecting to the Network File Storage cluster.",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Description: "Optional description for the clients groups.",
							Optional:    true,
							Computed:    true,
						},
						"ip_networks": {
							Type:        schema.TypeList,
							Description: "The allowed host or network to which the export is being shared. The IP address can be either IPv4 or IPv6 and has to be given with CIDR notation.",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hosts": {
							Type:        schema.TypeList,
							Description: "A singular host allowed to connect to the share. The host can be specified as IP address and can be either IPv4 or IPv6.",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"nfs": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"squash": {
										Type:        schema.TypeString,
										Description: "The squash mode for the export. The squash mode can be: none - No squash mode. no mapping, root-anonymous - Map root user to anonymous uid, all-anonymous - Map all users to anonymous uid.",
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using the name filter.",
				Default:     false,
				Optional:    true,
			},
		},
	}
}

func dataSourceNFSShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	idValue, idOk := d.GetOk("id")
	nameValue, nameOk := d.GetOk("name")
	partialMatch := d.Get("partial_match").(bool)
	id := idValue.(string)
	name := nameValue.(string)
	location := d.Get("location").(string)
	clusterID := d.Get("cluster_id").(string)

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be both specified at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the NFS Share ID or name"))
	}

	var share nfs.ShareRead
	var err error
	if idOk {
		share, _, err = client.GetNFSShareByID(ctx, clusterID, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the NFS Share with ID: %s, error: %w", idValue, err))
		}
	} else {
		shares, _, err := client.ListNFSShares(ctx, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching NFS Shares: %w", err))
		}

		var results []nfs.ShareRead
		for _, sh := range *shares.Items {
			if sh.Properties != nil && sh.Properties.Name != nil && utils.NameMatches(*sh.Properties.Name, name, partialMatch) {
				results = append(results, sh)
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no NFS Share found with the specified name: %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one NFS Share found with the specified name: %s", name))
		} else {
			share = results[0]
		}
	}

	if err = client.SetNFSShareData(d, share); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
