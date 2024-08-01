package cdn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cdn "github.com/ionos-cloud/sdk-go-cdn"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func SetDistributionData(d *schema.ResourceData, distribution cdn.Distribution) error {
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
		if err := d.Set("certificate_id", *distribution.Properties.CertificateId); err != nil {
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
					geoRestrictionsList := make([]interface{}, 0)
					geoRestrictionsList = append(geoRestrictionsList, geoRestrictionsEntry)
					upstreamEntry["geo_restrictions"] = geoRestrictionsList
				}
				upstreamList := make([]interface{}, 0)
				upstreamList = append(upstreamList, upstreamEntry)
				ruleEntry["upstream"] = upstreamList
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

	routingRulesVal := d.Get("routing_rules").([]interface{})
	routingRules := make([]cdn.RoutingRule, 0)
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
			if host, hostOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.host", routingRuleIndex)); hostOk {
				host := host.(string)
				routingRule.Upstream.Host = &host
			}
			if caching, cachingOk := d.GetOkExists(fmt.Sprintf("routing_rules.%d.upstream.0.caching", routingRuleIndex)); cachingOk { //nolint:staticcheck
				caching := caching.(bool)
				routingRule.Upstream.Caching = &caching
			}
			if waf, wafOk := d.GetOkExists(fmt.Sprintf("routing_rules.%d.upstream.0.waf", routingRuleIndex)); wafOk { //nolint:staticcheck
				waf := waf.(bool)
				routingRule.Upstream.Waf = &waf
			}

			if _, geoRestrictionsOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.geo_restrictions", routingRuleIndex)); geoRestrictionsOk {
				routingRule.Upstream.GeoRestrictions = &cdn.UpstreamGeoRestrictions{}
				if allowList, allowListOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.geo_restrictions.0.allow_list", routingRuleIndex)); allowListOk {
					raw := allowList.([]interface{})
					if len(raw) > 0 {
						countries := make([]string, 0)
						for _, rawCountry := range raw {
							if rawCountry != nil {
								ip := rawCountry.(string)
								countries = append(countries, ip)
							}
						}
						if len(countries) > 0 {
							routingRule.Upstream.GeoRestrictions.AllowList = &countries
						}
					}
				}
				if blockList, blockListOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.geo_restrictions.0.block_list", routingRuleIndex)); blockListOk {
					raw := blockList.([]interface{})
					if len(raw) > 0 {
						countries := make([]string, 0)
						for _, rawCountry := range raw {
							if rawCountry != nil {
								ip := rawCountry.(string)
								countries = append(countries, ip)
							}
						}
						if len(countries) > 0 {
							routingRule.Upstream.GeoRestrictions.BlockList = &countries
						}
					}
				}
			}
			if rateLimitClass, rateLimitClassOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.rate_limit_class", routingRuleIndex)); rateLimitClassOk {
				rateLimitClass := rateLimitClass.(string)
				routingRule.Upstream.RateLimitClass = &rateLimitClass
			}
		}

		routingRules = append(routingRules, routingRule)

	}

	return &routingRules, nil
}
