package profitbricks

import "net/http"

// PrivateCrossConnect type
type PrivateCrossConnect struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// metadata
	Metadata *Metadata `json:"metadata,omitempty"`

	// properties
	// Required: true
	Properties *PrivateCrossConnectProperties `json:"properties"`

	// The type of object
	// Read Only: true
	// Enum: [pcc]
	PBType string `json:"type,omitempty"`
}

// PrivateCrossConnects type
type PrivateCrossConnects struct {
	// URL to the collection representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// Unique representation for private cros-connect as a collection on a resource.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// Slice of items in that collection
	// Read Only: true
	Items []PrivateCrossConnect `json:"items"`

	// The type of resource within a collection
	// Read Only: true
	// Enum: [collection]
	PBType string `json:"type,omitempty"`
}

// PrivateCrossConnectProperties type
type PrivateCrossConnectProperties struct {
	// The desired name for the PrivateCrossConnect
	// Required: true
	Name string `json:"name,omitempty"`
	// A description for this PrivateCrossConnect
	// Required: true
	Description string `json:"description,omitempty"`
	// The peers of the PrivateCrossConnect
	// Required: false
	// Readonly: true
	Peers *[]PCCPeer `json:"peers,omitempty"`
	// The Connectable VDC's
	// Required: false
	// Readonly: true
	ConnectableDatacenters *[]PCCConnectableDataCenter `json:"connectableDatacenters,omitempty"`
}

// PCCPeer type
type PCCPeer struct {
	// The id of the cross-connected LAN
	// Required: false
	LANId string `json:"id,omitempty"`
	// The name of the cross-connected LAN
	// Required: false
	LANName string `json:"name,omitempty"`
	// The id of the cross-connected VDC
	// Required: false
	DataCenterID string `json:"datacenterId,omitempty"`
	// The name of the cross-connected VDC
	// Required: false
	DataCenterName string `json:"datacenterName,omitempty"`
	// The location of the cross-connected VDC
	// Required: false
	Location string `json:"location,omitempty"`
}

// PCCConnectableDataCenter type
type PCCConnectableDataCenter struct {
	// The id of the cross-connectable VDC
	// Required: false
	ID string `json:"id,omitempty"`
	// The name of the cross-connectable VDC
	// Required: false
	Name string `json:"name,omitempty"`
	// The name of the cross-connectable VDC
	// Required: false
	Location string `json:"location,omitempty"`
}

// ListPrivateCrossConnects gets a list of all private cross-connects
func (c *Client) ListPrivateCrossConnects() (*PrivateCrossConnects, error) {
	rsp := &PrivateCrossConnects{}
	return rsp, c.GetOK(PrivateCrossConnectsPath(), rsp)
}

// GetPrivateCrossConnect gets a private cross-connect with given id
func (c *Client) GetPrivateCrossConnect(pccID string) (*PrivateCrossConnect, error) {
	rsp := &PrivateCrossConnect{}
	return rsp, c.GetOK(PrivateCrossConnectPath(pccID), rsp)
}

// CreatePrivateCrossConnect creates a private cross-connect
func (c *Client) CreatePrivateCrossConnect(pcc PrivateCrossConnect) (*PrivateCrossConnect, error) {
	rsp := &PrivateCrossConnect{}
	return rsp, c.PostAcc(PrivateCrossConnectsPath(), pcc, rsp)
}

// UpdatePrivateCrossConnect updates a private cross-connect
func (c *Client) UpdatePrivateCrossConnect(pccID string, pcc PrivateCrossConnect) (*PrivateCrossConnect, error) {
	rsp := &PrivateCrossConnect{}
	return rsp, c.PatchAcc(PrivateCrossConnectPath(pccID), pcc.Properties, rsp)
}

// DeletePrivateCrossConnect deletes a private cross-connect by its id
func (c *Client) DeletePrivateCrossConnect(pccID string) (*http.Header, error) {
	h := &http.Header{}
	return h, c.Delete(PrivateCrossConnectPath(pccID), h, http.StatusAccepted)
}
