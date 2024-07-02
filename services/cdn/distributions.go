package cdn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cdn "github.com/ionos-cloud/sdk-go-cdn"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func SetDistributionData(d *schema.ResourceData, distribution cdn.DistributionRead) error {
	resourceName := "distribution"

	if distribution.Id != nil {
		d.SetId(*distribution.Id)
	}

	if distribution.Properties.Domain != nil {
		if err := d.Set("domain", *distribution.Properties.Domain); err != nil {
			return utils.GenerateSetError(resourceName, "domain", err)
		}
	}

	if distribution.Properties.CertificateId != nil {
		if err := d.Set("certificateId", *distribution.Properties.CertificateId); err != nil {
			return utils.GenerateSetError(resourceName, "certificateId", err)
		}
	}

	routingRules := make([]interface{}, 0)
	if distribution.Properties.RoutingRules != nil && len(*distribution.Properties.RoutingRules) > 0 {
		routingRules = make([]interface{}, 0)
		for _, rule := range *distribution.Properties.RoutingRules {
			ruleEntry := make(map[string]interface{})

			if rule.Scheme != nil {
				ruleEntry["scheme"] = *rule.Scheme
			}

			if rule.Prefix != nil {
				ruleEntry["prefix"] = *rule.Prefix
			}

			if rule.Upstream != nil {
				upstreamEntry := make(map[string]interface{})
				if rule.Upstream.Caching != nil {
					upstreamEntry["caching"] = *rule.Upstream.Caching
				}
				if rule.Upstream.Waf != nil {
					upstreamEntry["waf"] = *rule.Upstream.Waf
				}

				if rule.Upstream.Host != nil {
					upstreamEntry["host"] = *rule.Upstream.Host
				}

				if rule.Upstream.RateLimitClass != nil {
					upstreamEntry["rate_limit_class"] = *rule.Upstream.RateLimitClass
				}

				if rule.Upstream.GeoRestrictions != nil {
					geoRestrictionsEntry := make(map[string]interface{})

					if rule.Upstream.GeoRestrictions.AllowList != nil {
						geoRestrictionsEntry["allow_list"] = *rule.Upstream.GeoRestrictions.AllowList
					}
					if rule.Upstream.GeoRestrictions.BlockList != nil {
						geoRestrictionsEntry["block_list"] = *rule.Upstream.GeoRestrictions.BlockList
					}
					upstreamEntry["geo_restrictions"] = geoRestrictionsEntry
				}

				ruleEntry["upstream"] = upstreamEntry
			}

			routingRules = append(routingRules, ruleEntry)
		}
	}

	if len(routingRules) > 0 {
		if err := d.Set("routing_rules", routingRules); err != nil {
			return fmt.Errorf("error while setting routing_rules property for distribution  %s: %w", d.Id(), err)
		}
	}

	return nil
}

func GetRoutingRulesData(d *schema.ResourceData) (*[]cdn.RoutingRule, error) {
	var routingRules []cdn.RoutingRule

	routingRulesVal := d.Get("routing_rules").([]interface{})

	for routingRuleIndex := range routingRulesVal {

		routingRule := cdn.RoutingRule{}

		if scheme, schemeOk := d.GetOk(fmt.Sprintf("routing_rules.%d.scheme", routingRuleIndex)); schemeOk {
			scheme := scheme.(string)
			routingRule.Scheme = &scheme
		}

		if prefix, prefixOk := d.GetOk(fmt.Sprintf("routing_rules.%d.prefix", routingRuleIndex)); prefixOk {
			prefix := prefix.(string)
			routingRule.Prefix = &prefix
		}

		if _, upstreamOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream", routingRuleIndex)); upstreamOk {
			routingRule.Upstream = &cdn.Upstream{}
			if host, hostOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.host", routingRuleIndex)); hostOk {
				host := host.(string)
				routingRule.Upstream.Host = &host
			}
			if caching, cachingOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.caching", routingRuleIndex)); cachingOk {
				caching := caching.(bool)
				routingRule.Upstream.Caching = &caching
			}
			if waf, wafOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.waf", routingRuleIndex)); wafOk {
				waf := waf.(bool)
				routingRule.Upstream.Waf = &waf
			}

			if _, geo_restrictionsOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.geo_restrictions", routingRuleIndex)); geo_restrictionsOk {
				routingRule.Upstream.GeoRestrictions = &cdn.UpstreamGeoRestrictions{}
				if allowList, allowListOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.geo_restrictions.allow_list", routingRuleIndex)); allowListOk {
					allowList := allowList.([]string)
					routingRule.Upstream.GeoRestrictions.AllowList = &allowList
				}
				if blockList, blockListOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.geo_restrictions.block_list", routingRuleIndex)); blockListOk {
					blockList := blockList.([]string)
					routingRule.Upstream.GeoRestrictions.BlockList = &blockList
				}
			}
			if rateLimitClass, rateLimitClassOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.rateLimitClass", routingRuleIndex)); rateLimitClassOk {
				rateLimitClass := rateLimitClass.(string)
				routingRule.Upstream.RateLimitClass = &rateLimitClass
			}
		}

		routingRules = append(routingRules, routingRule)

	}

	return &routingRules, nil
}
