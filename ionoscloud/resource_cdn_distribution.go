package ionoscloud

import (
	"context"
	"fmt"
	"log"

	ionoscloud_cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCDNDistribution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCDNDistributionCreate,
		ReadContext:   resourceCDNDistributionRead,
		UpdateContext: resourceCDNDistributionUpdate,
		DeleteContext: resourceCDNDistributionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCDNDistributionImport,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:             schema.TypeString,
				Description:      "The domain of the distribution.",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Description: "The ID of the certificate to use for the distribution.",
				Optional:    true,
			},
			"routing_rules": {
				Type:        schema.TypeList,
				MaxItems:    20,
				Description: "The routing rules for the distribution.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:             schema.TypeString,
							Description:      "The scheme of the routing rule.",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"http", "https", "http/https"}, true)),
						},
						"prefix": {
							Type:        schema.TypeString,
							Description: "The prefix of the routing rule.",
							Required:    true,
						},
						"upstream": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Description: "The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.",
										Required:    true,
									},
									"caching": {
										Type:        schema.TypeBool,
										Description: "Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.",
										Required:    true,
									},
									"waf": {
										Type:        schema.TypeBool,
										Description: "Enable or disable WAF to protect the upstream host.",
										Required:    true,
									},
									"geo_restrictions": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"block_list": {
													Type: schema.TypeList,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional: true,
												},
												"allow_list": {
													Type: schema.TypeList,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional: true,
												},
											},
										},
									},
									"rate_limit_class": {
										Type:             schema.TypeString,
										Description:      "Rate limit class that will be applied to limit the number of incoming requests per IP.",
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"R1", "R5", "R10", "R25", "R50", "R100", "R250", "R500"}, true)),
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

func resourceCDNDistributionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(services.SdkBundle).CDNClient

	distributionDomain := d.Get("domain").(string)

	distribution := ionoscloud_cdn.DistributionCreate{
		Properties: ionoscloud_cdn.DistributionProperties{
			Domain: distributionDomain,
		},
	}

	if attr, ok := d.GetOk("certificate_id"); ok {
		attrStr := attr.(string)
		distribution.Properties.CertificateId = &attrStr
	}

	if routingRules, err := cdnService.GetRoutingRulesData(d); err == nil {
		distribution.Properties.RoutingRules = *routingRules
	} else {
		return diag.FromErr(err)
	}

	createdDistribution, _, err := client.SdkClient.DistributionsApi.DistributionsPost(ctx).DistributionCreate(distribution).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating CDN distribution (%s) (%w)", d.Id(), err))
		return diags
	}
	d.SetId(createdDistribution.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsDistributionReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for the CDN Distribution with ID: %v, error: %w", d.Id(), err))
	}

	log.Printf("[INFO] CDN Distribution Id: %s", d.Id())

	return resourceCDNDistributionRead(ctx, d, meta)
}

func resourceCDNDistributionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CDNClient

	distribution, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching CDN distribution %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved CDN distribution %s: %+v", d.Id(), distribution)

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceCDNDistributionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CDNClient

	distributionDomain := d.Get("domain").(string)

	request := ionoscloud_cdn.DistributionUpdate{
		Id: d.Id(),
		Properties: ionoscloud_cdn.DistributionProperties{
			Domain: distributionDomain,
		},
	}

	if attr, ok := d.GetOk("certificate_id"); ok {
		attrStr := attr.(string)
		request.Properties.CertificateId = &attrStr
	}

	if routingRules, err := cdnService.GetRoutingRulesData(d); err == nil {
		request.Properties.RoutingRules = *routingRules
	} else {
		return diag.FromErr(err)
	}

	_, _, err := client.SdkClient.DistributionsApi.DistributionsPut(ctx, d.Id()).DistributionUpdate(request).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating a CDN distribution ID %s %w",
			d.Id(), err))
		return diags
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsDistributionReady)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error occurred while checking the status for the CDN Distribution with ID: %v, error: %w", d.Id(), err))
	}

	return resourceCDNDistributionRead(ctx, d, meta)
}

func resourceCDNDistributionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CDNClient

	_, err := client.SdkClient.DistributionsApi.DistributionsDelete(ctx, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	d.SetId("")
	return nil
}

func resourceCDNDistributionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := *meta.(services.SdkBundle).CDNClient

	distributionID := d.Id()

	distribution, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsFindById(ctx, distributionID).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("registry does not exist %q", distributionID)
		}
		return nil, fmt.Errorf("an error occurred while trying to fetch the import of CDN distribution %q, error:%w", distributionID, err)
	}

	log.Printf("[INFO] CDN distribution found: %+v", distribution)

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
