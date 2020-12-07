package profitbricks

import (
	"net/http"
)

//Image object
type Image struct {
	ID         string          `json:"id,omitempty"`
	PBType     string          `json:"type,omitempty"`
	Href       string          `json:"href,omitempty"`
	Metadata   *Metadata       `json:"metadata,omitempty"`
	Properties ImageProperties `json:"properties,omitempty"`
	Response   string          `json:"Response,omitempty"`
	Headers    *http.Header    `json:"headers,omitempty"`
	StatusCode int             `json:"statuscode,omitempty"`
}

//ImageProperties object
type ImageProperties struct {
	Name                string       `json:"name,omitempty"`
	Description         string       `json:"description,omitempty"`
	Location            string       `json:"location,omitempty"`
	Size                float64      `json:"size,omitempty"`
	CPUHotPlug          bool         `json:"cpuHotPlug,omitempty"`
	CPUHotUnplug        bool         `json:"cpuHotUnplug,omitempty"`
	RAMHotPlug          bool         `json:"ramHotPlug,omitempty"`
	RAMHotUnplug        bool         `json:"ramHotUnplug,omitempty"`
	NicHotPlug          bool         `json:"nicHotPlug,omitempty"`
	NicHotUnplug        bool         `json:"nicHotUnplug,omitempty"`
	DiscVirtioHotPlug   bool         `json:"discVirtioHotPlug,omitempty"`
	DiscVirtioHotUnplug bool         `json:"discVirtioHotUnplug,omitempty"`
	DiscScsiHotPlug     bool         `json:"discScsiHotPlug,omitempty"`
	DiscScsiHotUnplug   bool         `json:"discScsiHotUnplug,omitempty"`
	LicenceType         string       `json:"licenceType,omitempty"`
	ImageType           string       `json:"imageType,omitempty"`
	ImageAliases        []string     `json:"imageAliases,omitempty"`
	Public              bool         `json:"public,omitempty"`
	Response            string       `json:"Response,omitempty"`
	Headers             *http.Header `json:"headers,omitempty"`
	StatusCode          int          `json:"statuscode,omitempty"`
}

//Images object
type Images struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Image      `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

//Cdroms object
type Cdroms struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Image      `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// ListImages returns an Collection struct
func (c *Client) ListImages() (*Images, error) {
	url := imagesPath()
	ret := &Images{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// GetImage returns an Instance struct where id ==imageid
func (c *Client) GetImage(imageid string) (*Image, error) {
	url := imagePath(imageid)
	ret := &Image{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}
