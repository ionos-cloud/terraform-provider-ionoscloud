package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dataplatformService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dataplatform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceDataplatformCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataplatformClusterCreate,
		ReadContext:   resourceDataplatformClusterRead,
		UpdateContext: resourceDataplatformClusterUpdate,
		DeleteContext: resourceDataplatformClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDataplatformClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:             schema.TypeString,
				Description:      "The UUID of the virtual data center (VDC) in which the cluster is provisioned",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
				Required:         true,
				ForceNew:         true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]). It can contain dashes (-), underscores (_), dots (.), and alphanumerics in-between.",
				ValidateDiagFunc: validation.AllDiag(validation.ToDiagFunc(validation.StringLenBetween(0, 63)),
					validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(constant.DataPlatformNameRegexConstraint), constant.DataPlatformRegexNameError))),
				Required: true,
			},
			"version": {
				Type:             schema.TypeString,
				Description:      "The version of the Data Platform.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 32)),
				Optional:         true,
				Computed:         true,
			},
			"maintenance_window": {
				Type:        schema.TypeList,
				Description: "Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:             schema.TypeString,
							Description:      "Time at which the maintenance should start.",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^(?:[01]\\d|2[0-3]):(?:[0-5]\\d):(?:[0-5]\\d)$"), "")),
							Required:         true,
						},
						"day_of_the_week": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsDayOfTheWeek(true)),
							Required:         true,
						},
					},
				},
			},
			"lans": {
				Type:        schema.TypeSet,
				Description: "A list of LANs you want this node pool to be part of",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The LAN ID of an existing LAN at the related data center",
							Required:    true,
						},
						"dhcp": {
							Type:        schema.TypeBool,
							Description: "Indicates if the Kubernetes node pool LAN will reserve an IP using DHCP. The default value is 'true'",
							Optional:    true,
							Default:     true,
						},
						"routes": {
							Type:        schema.TypeSet,
							Description: "An array of additional LANs attached to worker nodes",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network": {
										Type:             schema.TypeString,
										Description:      "IPv4 or IPv6 CIDR to be routed via the interface",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
										Required:         true,
									},
									"gateway": {
										Type:             schema.TypeString,
										Description:      "IPv4 or IPv6 gateway IP for the route",
										ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPAddress),
										Required:         true,
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

func resourceDataplatformClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DataplatformClient

	id, _, err := client.CreateCluster(ctx, d)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while creating the Dataplatform cluster: %w", err))
		return diags
	}

	d.SetId(id)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while waiting for the Dataplatform cluster with ID: %v to become available: %w", id, err))
		return diags
	}

	return resourceDataplatformClusterRead(ctx, d, meta)
}

func resourceDataplatformClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(bundleclient.SdkBundle).DataplatformClient

	clusterId := d.Id()
	dataplatformCluster, apiResponse, err := client.GetClusterById(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Could not find Dataplatform cluster with ID: %s", clusterId)
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching Dataplatform cluster with ID: %s, error: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved Dataplatform cluster %s: %+v", d.Id(), dataplatformCluster)

	if err := dataplatformService.SetDataplatformClusterData(d, dataplatformCluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDataplatformClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DataplatformClient

	clusterId := d.Id()

	_, err := client.UpdateCluster(ctx, clusterId, d)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating the Dataplatform cluster with ID: %v, error: %w", clusterId, err))
		return diags
	}

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while waiting for the Dataplatform cluster to become available after update, ID: %v, error: %w", clusterId, err))
	}

	return resourceDataplatformClusterRead(ctx, d, meta)
}

func resourceDataplatformClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).DataplatformClient

	clusterId := d.Id()

	apiResponse, err := client.DeleteCluster(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Dataplatform cluster with ID: %v, error: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		diag.FromErr(fmt.Errorf("an error occurred while waiting for the Dataplatform cluster with ID: %v to be deleted, error: %w", err))
	}

	return nil
}

func resourceDataplatformClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).DataplatformClient

	clusterId := d.Id()

	dataplatformCluster, apiResponse, err := client.GetClusterById(ctx, clusterId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("Dataplatform cluster with ID: %v does not exist", clusterId)
		}
		return nil, fmt.Errorf("an error occurred while trying to import the Dataplatform cluster with ID: %v, error: %w", clusterId, err)
	}

	log.Printf("[INFO] Dataplatform cluster found: %+v", dataplatformCluster)

	if err := dataplatformService.SetDataplatformClusterData(d, dataplatformCluster); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
