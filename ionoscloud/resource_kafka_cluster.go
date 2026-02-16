package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceKafkaCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKafkaClusterCreate,
		ReadContext:   resourceKafkaClusterRead,
		DeleteContext: resourceKafkaClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKafkaClusterImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Kafka Cluster.",
				Computed:    true,
			},
			"name": {
				Type:             schema.TypeString,
				Description:      "The name of your Kafka Cluster. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("The location of your Kafka Cluster. Supported locations: %s", strings.Join(kafka.AvailableLocations, ", ")),
				Optional:    true,
				ForceNew:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The desired Kafka Version. Supported versions: 3.9.0",
				ForceNew:    true,
				Required:    true,
			},
			"size": {
				Type:             schema.TypeString,
				Description:      "The size of your Kafka Cluster. The size of the Kafka Cluster is given in T-shirt sizes. Valid values are: XS, S",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"connections": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Description: "The network connection for your Kafka Cluster. Only one connection is allowed.",
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:             schema.TypeString,
							Description:      "The datacenter to connect your Kafka Cluster to.",
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
						},
						"lan_id": {
							Type:        schema.TypeString,
							Description: "The numeric LAN ID to connect your Kafka Cluster to.",
							Required:    true,
							ForceNew:    true,
						},
						"broker_addresses": {
							Type:        schema.TypeList,
							Description: "The broker addresses of the Kafka Cluster. Can be empty, but must be present.",
							Required:    true,
							ForceNew:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"broker_addresses": {
				Type:        schema.TypeList,
				Description: "IP address and port of cluster brokers.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceKafkaClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).KafkaClient

	createdCluster, apiResponse, err := client.CreateCluster(ctx, d)
	if err != nil {
		d.SetId("")
		return utils.ToDiags(d, fmt.Sprintf("error creating Kafka Cluster: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	d.SetId(createdCluster.Id)
	log.Printf("[INFO] Created Kafka Cluster: %s", d.Id())

	// Sleep for 5 second to avoid 500 error from the API
	time.Sleep(5 * time.Second)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterAvailable)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while Kafka Cluster waiting to be ready: %s", err), nil)
	}

	return resourceKafkaClusterRead(ctx, d, meta)
}

func resourceKafkaClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).KafkaClient

	cluster, apiResponse, err := client.GetClusterByID(ctx, d.Id(), d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching Kafka Cluster: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retreived Kafka Cluster %s: %+v", d.Id(), cluster)

	if err := client.SetKafkaClusterData(d, &cluster); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceKafkaClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).KafkaClient

	apiResponse, err := client.DeleteCluster(ctx, d.Id(), d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting Kafka Cluster: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("the check for Kafka Cluster deletion failed with the following error: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}

	d.SetId("")

	return nil
}

func resourceKafkaClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "expected ID in the format location:cluster_id", nil)
	}

	location := parts[0]
	if err := d.Set("location", location); err != nil {
		return nil, utils.ToError(d, fmt.Sprintf("failed to set location Kafka Cluster for import: %s", err), nil)
	}
	d.SetId(parts[1])

	diags := resourceKafkaClusterRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, utils.ToError(d, diags[0].Summary, nil)
	}

	return []*schema.ResourceData{d}, nil
}
