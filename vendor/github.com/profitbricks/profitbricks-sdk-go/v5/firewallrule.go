package profitbricks

import (
	"net/http"
)

//FirewallRule object
type FirewallRule struct {
	ID         string                 `json:"id,omitempty"`
	PBType     string                 `json:"type,omitempty"`
	Href       string                 `json:"href,omitempty"`
	Metadata   *Metadata              `json:"metadata,omitempty"`
	Properties FirewallruleProperties `json:"properties,omitempty"`
	Response   string                 `json:"Response,omitempty"`
	Headers    *http.Header           `json:"headers,omitempty"`
	StatusCode int                    `json:"statuscode,omitempty"`
}

//FirewallruleProperties object
type FirewallruleProperties struct {
	Name           string  `json:"name"`
	Protocol       string  `json:"protocol,omitempty"`
	SourceMac      *string `json:"sourceMac"`
	SourceIP       *string `json:"sourceIp"`
	TargetIP       *string `json:"targetIp"`
	IcmpCode       *int    `json:"icmpCode"`
	IcmpType       *int    `json:"icmpType"`
	PortRangeStart *int    `json:"portRangeStart"`
	PortRangeEnd   *int    `json:"portRangeEnd"`
}

//FirewallRules object
type FirewallRules struct {
	ID         string         `json:"id,omitempty"`
	PBType     string         `json:"type,omitempty"`
	Href       string         `json:"href,omitempty"`
	Items      []FirewallRule `json:"items,omitempty"`
	Response   string         `json:"Response,omitempty"`
	Headers    *http.Header   `json:"headers,omitempty"`
	StatusCode int            `json:"statuscode,omitempty"`
}

//ListFirewallRules lists all firewall rules
func (c *Client) ListFirewallRules(dcID string, serverID string, nicID string) (*FirewallRules, error) {
	url := firewallRulesPath(dcID, serverID, nicID)
	ret := &FirewallRules{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//GetFirewallRule gets a firewall rule
func (c *Client) GetFirewallRule(dcID string, serverID string, nicID string, fwID string) (*FirewallRule, error) {
	url := firewallRulePath(dcID, serverID, nicID, fwID)
	ret := &FirewallRule{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//CreateFirewallRule creates a firewall rule
func (c *Client) CreateFirewallRule(dcID string, serverID string, nicID string, fw FirewallRule) (*FirewallRule, error) {
	url := firewallRulesPath(dcID, serverID, nicID)
	ret := &FirewallRule{}
	err := c.Post(url, fw, ret, http.StatusAccepted)
	return ret, err
}

// UpdateFirewallRule updates a firewall rule.
// You need to pass all wanted properties, not just those you want to change.
func (c *Client) UpdateFirewallRule(dcID string, serverID string, nicID string, fwID string, obj FirewallruleProperties) (*FirewallRule, error) {
	url := firewallRulePath(dcID, serverID, nicID, fwID)
	ret := &FirewallRule{}
	err := c.Patch(url, obj, ret, http.StatusAccepted)
	return ret, err
}

//DeleteFirewallRule deletes a firewall rule
func (c *Client) DeleteFirewallRule(dcID string, serverID string, nicID string, fwID string) (*http.Header, error) {
	url := firewallRulePath(dcID, serverID, nicID, fwID)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}
