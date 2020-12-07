package profitbricks

import (
	"net/http"
)

//ContractResources object
type ContractResources struct {
	ID         string                      `json:"id,omitempty"`
	PBType     string                      `json:"type,omitempty"`
	Href       string                      `json:"href,omitempty"`
	Properties ContractResourcesProperties `json:"properties,omitempty"`
	Response   string                      `json:"Response,omitempty"`
	Headers    *http.Header                `json:"headers,omitempty"`
	StatusCode int                         `json:"statuscode,omitempty"`
}

//ContractResourcesProperties object
type ContractResourcesProperties struct {
	PBContractNumber int64            `json:"contractNumber,omitempty"`
	Owner            string           `json:"owner,omitempty"`
	Status           string           `json:"status,omitempty"`
	ResourceLimits   *ResourcesLimits `json:"resourceLimits,omitempty"`
}

//ResourcesLimits object
type ResourcesLimits struct {
	CoresPerServer        int32 `json:"coresPerServer,omitempty"`
	CoresPerContract      int32 `json:"coresPerContract,omitempty"`
	CoresProvisioned      int32 `json:"coresProvisioned,omitempty"`
	RAMPerServer          int32 `json:"ramPerServer,omitempty"`
	RAMPerContract        int32 `json:"ramPerContract,omitempty"`
	RAMProvisioned        int32 `json:"ramProvisioned,omitempty"`
	HddLimitPerVolume     int64 `json:"hddLimitPerVolume,omitempty"`
	HddLimitPerContract   int64 `json:"hddLimitPerContract,omitempty"`
	HddVolumeProvisioned  int64 `json:"hddVolumeProvisioned,omitempty"`
	SsdLimitPerVolume     int64 `json:"ssdLimitPerVolume,omitempty"`
	SsdLimitPerContract   int64 `json:"ssdLimitPerContract,omitempty"`
	SsdVolumeProvisioned  int64 `json:"ssdVolumeProvisioned,omitempty"`
	ReservableIps         int32 `json:"reservableIps,omitempty"`
	ReservedIpsOnContract int32 `json:"reservedIpsOnContract,omitempty"`
	ReservedIpsInUse      int32 `json:"reservedIpsInUse,omitempty"`
}

// GetContractResources returns list of contract resources
func (c *Client) GetContractResources() (*ContractResources, error) {
	url := contractsPath()
	ret := &ContractResources{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err

}
