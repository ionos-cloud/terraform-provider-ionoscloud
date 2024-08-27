package cdn

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

// SetDistributionData sets distribution data from a distribution sdk object
func SetDistributionData(d *schema.ResourceData, distribution cdn.Distribution) error {
	resourceName := "distribution"

	d.SetId(distribution.Id)

	if err := d.Set("domain", distribution.Properties.Domain); err != nil {
		return utils.GenerateSetError(resourceName, "domain", err)
	}

	if err := d.Set("certificate_id", distribution.Properties.GetCertificateId()); err != nil {
		return utils.GenerateSetError(resourceName, "certificate_id", err)
	}

	routingRules := make([]interface{}, 0)
	if len(distribution.Properties.RoutingRules) > 0 {
		routingRules = make([]interface{}, 0)
		for _, rule := range distribution.Properties.RoutingRules {
			ruleEntry := make(map[string]interface{})
			ruleEntry["scheme"] = rule.Scheme
			ruleEntry["prefix"] = rule.Prefix

			upstreamEntry := make(map[string]interface{})
			upstreamEntry["caching"] = rule.Upstream.Caching
			upstreamEntry["waf"] = rule.Upstream.Waf
			upstreamEntry["host"] = rule.Upstream.Host
			upstreamEntry["rate_limit_class"] = rule.Upstream.RateLimitClass

			geoRestrictionsEntry := make(map[string]interface{})
			if rule.Upstream.GeoRestrictions != nil {
				geoRestrictionsEntry["allow_list"] = rule.Upstream.GeoRestrictions.AllowList
				geoRestrictionsEntry["block_list"] = rule.Upstream.GeoRestrictions.BlockList
				geoRestrictionsList := make([]interface{}, 0)
				geoRestrictionsList = append(geoRestrictionsList, geoRestrictionsEntry)
				upstreamEntry["geo_restrictions"] = geoRestrictionsList
			}
			upstreamList := make([]interface{}, 0)
			upstreamList = append(upstreamList, upstreamEntry)
			ruleEntry["upstream"] = upstreamList

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

// GetRoutingRulesData gets distribution routing rules data from terraform
func GetRoutingRulesData(d *schema.ResourceData) (*[]cdn.RoutingRule, error) {

	routingRulesVal := d.Get("routing_rules").([]interface{})
	routingRules := make([]cdn.RoutingRule, 0)
	for routingRuleIndex := range routingRulesVal {

		routingRule := cdn.RoutingRule{}

		if scheme, schemeOk := d.GetOk(fmt.Sprintf("routing_rules.%d.scheme", routingRuleIndex)); schemeOk {
			routingRule.Scheme = scheme.(string)
		}

		if prefix, prefixOk := d.GetOk(fmt.Sprintf("routing_rules.%d.prefix", routingRuleIndex)); prefixOk {
			routingRule.Prefix = prefix.(string)
		}

		if _, upstreamOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream", routingRuleIndex)); upstreamOk {
			routingRule.Upstream = cdn.Upstream{}
			if host, hostOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.host", routingRuleIndex)); hostOk {
				routingRule.Upstream.Host = host.(string)
			}
			if caching, cachingOk := d.GetOkExists(fmt.Sprintf("routing_rules.%d.upstream.0.caching", routingRuleIndex)); cachingOk { //nolint:staticcheck
				routingRule.Upstream.Caching = caching.(bool)
			}
			if waf, wafOk := d.GetOkExists(fmt.Sprintf("routing_rules.%d.upstream.0.waf", routingRuleIndex)); wafOk { //nolint:staticcheck
				routingRule.Upstream.Waf = waf.(bool)
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
							routingRule.Upstream.GeoRestrictions.AllowList = countries
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
							routingRule.Upstream.GeoRestrictions.BlockList = countries
						}
					}
				}
			}
			if rateLimitClass, rateLimitClassOk := d.GetOk(fmt.Sprintf("routing_rules.%d.upstream.0.rate_limit_class", routingRuleIndex)); rateLimitClassOk {
				routingRule.Upstream.RateLimitClass = rateLimitClass.(string)
			}
		}

		routingRules = append(routingRules, routingRule)

	}

	return &routingRules, nil
}

// IsDistributionReady checks if the distribution is ready
func (c *Client) IsDistributionReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	distributionID := d.Id()
	distribution, _, err := c.SdkClient.DistributionsApi.DistributionsFindById(ctx, distributionID).Execute()
	if err != nil {
		return true, fmt.Errorf("status check failed for distribution with ID: %v, error: %w", distributionID, err)
	}

	log.Printf("[INFO] state of the distribution with ID: %v is: %s ", distributionID, distribution.Metadata.State)
	return strings.EqualFold(distribution.Metadata.State, constant.Available), nil
}
