package ionoscloud

import (
	"context"
	"fmt"
	"log"

	ionoscloud_cdn "github.com/ionos-cloud/sdk-go-cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDistribution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDistributionCreate,
		ReadContext:   resourceDistributionRead,
		UpdateContext: resourceDistributionUpdate,
		DeleteContext: resourceDistributionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDistributionImport,
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
				Required:    true,
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
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Description: "The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.",
										Required:    true,
									},
									"caching": {
										Type:        schema.TypeString,
										Description: "Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.",
										Required:    true,
									},
									"waf": {
										Type:        schema.TypeString,
										Description: "Enable or disable WAF to protect the upstream host.",
										Required:    true,
									},
									"geo_restrictions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"block_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional: true,
												},
												"allow_list": {
													Type: schema.TypeSet,
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
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"none", "R10", "R100", "R1000", "R10000"}, true)),
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

func resourceDistributionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(services.SdkBundle).CdnClient

	distributionDomain := d.Get("domain").(string)

	distribution := ionoscloud_cdn.DistributionCreate{
		Properties: &ionoscloud_cdn.Distribution{
			Domain: &distributionDomain,
		},
	}

	if attr, ok := d.GetOk("certificate_id"); ok {
		attrStr := attr.(string)
		distribution.Properties.CertificateId = &attrStr
	}

	if routingRules, err := cdnService.GetRoutingRulesData(d); err == nil {
		distribution.Properties.RoutingRules = routingRules
	} else {
		return diag.FromErr(err)
	}

	createdDistribution, _, err := client.DistributionsApi.DistributionsPost(ctx).DistributionCreate(distribution).Execute()
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating distribution (%s) (%w)", d.Id(), err))
		return diags
	}
	d.SetId(*createdDistribution.Id)

	log.Printf("[INFO] Distribution Id: %s", d.Id())

	return resourceDistributionRead(ctx, d, meta)
}

func resourceDistributionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CdnClient

	distribution, apiResponse, err := client.DistributionsApi.DistributionsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching distribution %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved distribution %s: %+v", d.Id(), distribution)

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDistributionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CdnClient

	request := ionoscloud_cdn.DistributionEnsure{
		Properties: &ionoscloud_cdn.Distribution{},
	}

	if d.HasChange("domain") {
		_, v := d.GetChange("domain")
		vStr := v.(string)
		request.Properties.Domain = &vStr
	}

	if d.HasChange("certificate_id") {
		_, v := d.GetChange("certificate_id")
		vStr := v.(string)
		request.Properties.CertificateId = &vStr
	}

	if d.HasChange("routing_rules") {
		if routingRules, err := cdnService.GetRoutingRulesData(d); err == nil {
			request.Properties.RoutingRules = routingRules
		} else {
			return diag.FromErr(err)
		}
	}

	_, _, err := client.DistributionsApi.DistributionsPut(ctx, d.Id()).DistributionEnsure(request).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating a distribution ID %s %w",
			d.Id(), err))
		return diags
	}

	return resourceDistributionRead(ctx, d, meta)
}

func resourceDistributionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CdnClient

	_, err := client.DistributionsApi.DistributionsDelete(ctx, d.Id()).Execute()

	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	d.SetId("")
	return nil
}

func resourceDistributionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).CdnClient

	distributionId := d.Id()

	distribution, apiResponse, err := client.DistributionsApi.DistributionsFindById(ctx, distributionId).Execute()

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("registry does not exist %q", distributionId)
		}
		return nil, fmt.Errorf("an error occurred while trying to fetch the import of distribution %q, error:%w", distributionId, err)
	}

	log.Printf("[INFO] distribution found: %+v", distribution)

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
