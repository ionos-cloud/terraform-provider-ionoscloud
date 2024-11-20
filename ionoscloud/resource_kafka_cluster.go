package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
				Description: "The desired Kafka Version. Supported version: 3.7.0",
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
	client := meta.(services.SdkBundle).KafkaClient

	createdCluster, _, err := client.CreateCluster(ctx, d)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating Kafka Cluster: %w", err))
		return diags
	}

	d.SetId(*createdCluster.Id)
	log.Printf("[INFO] Created Kafka Cluster: %s", d.Id())

	// Sleep for 5 second to avoid 500 error from the API
	time.Sleep(5 * time.Second)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsClusterAvailable)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while Kafka Cluster waiting to be ready: %w", err))
		return diags
	}

	return resourceKafkaClusterRead(ctx, d, meta)
}

func resourceKafkaClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient

	cluster, apiResponse, err := client.GetClusterByID(ctx, d.Id(), d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching Kafka Cluster %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived Kafka Cluster %s: %+v", d.Id(), cluster)

	if err := client.SetKafkaClusterData(d, &cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKafkaClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).KafkaClient

	apiResponse, err := client.DeleteCluster(ctx, d.Id(), d.Get("location").(string))
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting Kafka Cluster %s: %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsClusterDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("the check for Kafka Cluster deletion failed with the following error: %w", err))
	}

	d.SetId("")

	return nil
}

func resourceKafkaClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected ID in the format location:cluster_id")
	}

	location := parts[0]
	if !slices.Contains(kafka.AvailableLocations, location) {
		return nil, fmt.Errorf("invalid location: %v, location must be one of: %v", location, kafka.AvailableLocations)
	}

	if err := d.Set("location", parts[0]); err != nil {
		return nil, fmt.Errorf("failed to set location Kafka Cluster for import: %w", err)
	}
	d.SetId(parts[1])

	diags := resourceKafkaClusterRead(ctx, d, meta)
	if diags != nil && diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
