package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
)

func resourceNFSShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNFSShareCreate,
		ReadContext:   resourceNFSShareRead,
		UpdateContext: resourceNFSShareUpdate,
		DeleteContext: resourceNFSShareDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNFSShareImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network File Storage Share.",
				Computed:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "The ID of the Network File Storage Cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"location": {
				Type: schema.TypeString,
				Description: fmt.Sprintf("The location of the Network File Storage Cluster. "+
					"Available locations: '%s'", strings.Join(ValidNFSLocations, ", '")),
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ValidNFSLocations, false)),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The directory being exported",
				Required:    true,
			},
			"quota": {
				Type:        schema.TypeInt,
				Description: "The quota in MiB for the export. The quota can restrict the amount of data that can be stored within the export. The quota can be disabled using `0`.",
				Optional:    true,
				Default:     0,
			},
			"gid": {
				Type:        schema.TypeInt,
				Description: "The group ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.",
				Optional:    true,
				// Default:     512, // max 65534
			},
			"uid": {
				Type:        schema.TypeInt,
				Description: "The user ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.",
				Optional:    true,
				// Default:     512, // max 65534
			},
			"client_groups": {
				Type:        schema.TypeList,
				Description: "The groups of clients are the systems connecting to the Network File Storage cluster.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Description: "Optional description for the clients groups.",
							Optional:    true,
						},
						"ip_networks": {
							Type:        schema.TypeList,
							Description: "The allowed host or network to which the export is being shared. The IP address can be either IPv4 or IPv6 and has to be given with CIDR notation.",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hosts": {
							Type:        schema.TypeList,
							Description: "A singular host allowed to connect to the share. The host can be specified as IP address and can be either IPv4 or IPv6.",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"nfs": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"squash": {
										Type:        schema.TypeString,
										Description: "The squash mode for the export. The squash mode can be: none - No squash mode. no mapping, root-anonymous - Map root user to anonymous uid, all-anonymous - Map all users to anonymous uid.",
										Optional:    true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
											"none",
											"root-anonymous",
											"all-anonymous",
										}, false)),
									},
								},
							},
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNFSShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	response, _, err := client.CreateNFSShare(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating NFS Share: %w", err))
	}
	shareID := *response.Id
	d.SetId(shareID)

	if err := client.SetNFSShareData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNFSShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	clusterID := d.Get("cluster_id").(string)
	location := d.Get("location").(string)
	shareID := d.Id()

	share, _, err := client.GetNFSShareById(ctx, clusterID, shareID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding NFS Share: %w", err))
	}

	if err := client.SetNFSShareData(d, share); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set NFS Share data: %w", err))
	}
	return nil
}

func resourceNFSShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	response, _, err := client.UpdateNFSShare(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating NFS Share: %w", err))
	}

	if err := client.SetNFSShareData(d, response); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNFSShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).NFSClient
	clusterID := d.Get("cluster_id").(string)
	location := d.Get("location").(string)
	shareID := d.Id()

	_, err := client.DeleteNFSShare(ctx, clusterID, shareID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting NFS Share with ID: %v, error: %w", shareID, err))
	}
	return nil
}

func resourceNFSShareImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).NFSClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID: %q, expected ID in the format '<location>:<cluster_id>:<share_id>'", d.Id())
	}
	location := parts[0]
	clusterID := parts[1]
	shareID := parts[2]

	d.Set("location", location)
	d.Set("cluster_id", clusterID)
	d.Set("id", shareID)

	share, _, err := client.GetNFSShareById(ctx, clusterID, shareID, location)
	if err != nil {
		return nil, fmt.Errorf("error finding NFS Share: %w", err)
	}

	if err := client.SetNFSShareData(d, share); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
