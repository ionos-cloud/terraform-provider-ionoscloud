package ionoscloud

import (
	"context"
	"fmt"
	"log"

	ionoscloudcdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"

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
			"resource_urn": {
				Type:        schema.TypeString,
				Description: "Unique name of the resource.",
				Computed:    true,
			},
			"public_endpoint_v4": {
				Type:        schema.TypeString,
				Description: "IP of the distribution, it has to be included on the domain DNS Zone as A record.",
				Computed:    true,
			},
			"public_endpoint_v6": {
				Type:        schema.TypeString,
				Description: "IP of the distribution, it has to be included on the domain DNS Zone as AAAA record.",
				Computed:    true,
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
									"sni_mode": {
										Type:        schema.TypeString,
										Description: "The SNI (Server Name Indication) mode of the upstream host. It supports two modes: 'distribution' and 'origin', for more information about these modes please check the resource docs.",
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

	client := meta.(bundleclient.SdkBundle).CDNClient

	distributionDomain := d.Get("domain").(string)

	distribution := ionoscloudcdn.DistributionCreate{
		Properties: ionoscloudcdn.DistributionProperties{
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
		return diagutil.ToDiags(d, err, nil)
	}

	createdDistribution, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsPost(ctx).DistributionCreate(distribution).Execute()
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error creating CDN distribution: (%w)", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	d.SetId(createdDistribution.Id)
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsDistributionReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error occurred while checking the status for the CDN Distribution: %w", err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String()})
	}

	log.Printf("[INFO] CDN Distribution Id: %s", d.Id())

	return resourceCDNDistributionRead(ctx, d, meta)
}

func resourceCDNDistributionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CDNClient

	distribution, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching CDN distribution: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	log.Printf("[INFO] Successfully retrieved CDN distribution %s: %+v", d.Id(), distribution)

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceCDNDistributionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CDNClient

	distributionDomain := d.Get("domain").(string)

	request := ionoscloudcdn.DistributionUpdate{
		Id: d.Id(),
		Properties: ionoscloudcdn.DistributionProperties{
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
		return diagutil.ToDiags(d, err, nil)
	}

	_, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsPut(ctx, d.Id()).DistributionUpdate(request).Execute()

	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating a CDN distribution ID %s %w",
			d.Id(), err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}
	err = utils.WaitForResourceToBeReady(ctx, d, client.IsDistributionReady)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("error occurred while checking the status for the CDN Distribution: %w", err), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String()})
	}

	return resourceCDNDistributionRead(ctx, d, meta)
}

func resourceCDNDistributionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CDNClient

	apiResponse, err := client.SdkClient.DistributionsApi.DistributionsDelete(ctx, d.Id()).Execute()

	if err != nil {
		return diagutil.ToDiags(d, err, &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId("")
	return nil
}

func resourceCDNDistributionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := *meta.(bundleclient.SdkBundle).CDNClient

	distributionID := d.Id()

	distribution, apiResponse, err := client.SdkClient.DistributionsApi.DistributionsFindById(ctx, distributionID).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("registry does not exist %q", distributionID), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to fetch the import of CDN distribution %q, error:%w", distributionID, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	log.Printf("[INFO] CDN distribution found: %+v", distribution)

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}
