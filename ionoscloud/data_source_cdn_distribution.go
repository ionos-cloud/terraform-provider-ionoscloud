package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	cdnService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cdn"
)

func dataSourceCDNDistribution() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDNDistributionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"partial_match": {
				Type:        schema.TypeBool,
				Description: "Whether partial matching is allowed or not when using domain argument.",
				Default:     false,
				Optional:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain of the distribution.",
				Optional:    true,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Description: "The ID of the certificate to use for the distribution.",
				Computed:    true,
			},
			"routing_rules": {
				Type:        schema.TypeList,
				Description: "The routing rules for the distribution.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:        schema.TypeString,
							Description: "The scheme of the routing rule.",
							Computed:    true,
						},
						"prefix": {
							Type:        schema.TypeString,
							Description: "The prefix of the routing rule.",
							Computed:    true,
						},
						"upstream": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Description: "The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.",
										Computed:    true,
									},
									"caching": {
										Type:        schema.TypeBool,
										Description: "Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.",
										Computed:    true,
									},
									"waf": {
										Type:        schema.TypeBool,
										Description: "Enable or disable WAF to protect the upstream host.",
										Computed:    true,
									},
									"sni_mode": {
										Type:        schema.TypeString,
										Description: "The SNI (Server Name Indication) mode of the upstream host. It supports two modes: 'distribution' and 'origin', for more information about these modes please check the data source docs.",
										Computed:    true,
									},
									"geo_restrictions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"block_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed: true,
												},
												"allow_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed: true,
												},
											},
										},
									},
									"rate_limit_class": {
										Type:        schema.TypeString,
										Description: "Rate limit class that will be applied to limit the number of incoming requests per IP.",
										Computed:    true,
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

func dataSourceCDNDistributionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).CDNClient

	idValue, idOk := d.GetOk("id")
	domainValue, domainOk := d.GetOk("domain")

	id := idValue.(string)
	domain := domainValue.(string)

	if idOk && domainOk {
		diags := diag.FromErr(errors.New("id and domain cannot be both specified in the same time"))
		return diags
	}
	if !idOk && !domainOk {
		diags := diag.FromErr(errors.New("please provide the distribution id or domain"))
		return diags
	}

	var distribution cdn.Distribution
	var err error

	if idOk {
		/* search by ID */
		distribution, _, err = client.SdkClient.DistributionsApi.DistributionsFindById(ctx, id).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching the distribution with ID %s: %w", id, err))
			return diags
		}
	} else {
		var results []cdn.Distribution

		distributions, _, err := client.SdkClient.DistributionsApi.DistributionsGet(ctx).Execute()
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occurred while fetching container distributions: %w", err))
			return diags
		}

		results = distributions.Items
		if domainOk {
			partialMatch := d.Get("partial_match").(bool)

			log.Printf("[INFO] Using data source for container registry by domain with partial_match %t and domain: %s", partialMatch, domain)

			if len(distributions.Items) > 0 {
				var distributionsByDomain []cdn.Distribution

				for _, distributionItem := range distributions.Items {
					if isValidDomain(distributionItem.Properties, domain, partialMatch) {
						distributionsByDomain = append(distributionsByDomain, distributionItem)
					}
				}
				if len(distributionsByDomain) > 0 {
					results = distributionsByDomain
				} else {
					return diag.FromErr(fmt.Errorf("no distribution found with the specified criteria: domain = %v", domain))
				}
			}
		}

		switch {
		case len(results) == 0:
			return diag.FromErr(fmt.Errorf("no CDN distribution found with the specified criteria: domain = %s", domain))
		case len(results) > 1:
			return diag.FromErr(fmt.Errorf("more than one CDN distribution found with the specified criteria: domain = %s", domain))
		default:
			distribution = results[0]
		}
	}

	if err := cdnService.SetDistributionData(d, distribution); err != nil {
		return diag.FromErr(err)
	}

	return nil

}

func isValidDomain(properties cdn.DistributionProperties, domain string, partialMatch bool) bool {
	if partialMatch {
		return strings.Contains(properties.Domain, domain)
	}
	return strings.EqualFold(properties.Domain, domain)
}
