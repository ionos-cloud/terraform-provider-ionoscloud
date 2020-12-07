package profitbricks

import (
	"context"
	"net/http"

	resty "github.com/go-resty/resty/v2"
)

// Volume object
type Volume struct {
	ID         string           `json:"id,omitempty"`
	PBType     string           `json:"type,omitempty"`
	Href       string           `json:"href,omitempty"`
	Metadata   *Metadata        `json:"metadata,omitempty"`
	Properties VolumeProperties `json:"properties,omitempty"`
	Response   string           `json:"Response,omitempty"`
	Headers    *http.Header     `json:"headers,omitempty"`
	StatusCode int              `json:"statuscode,omitempty"`
}

// VolumeProperties object
type VolumeProperties struct {
	Name                string   `json:"name,omitempty"`
	Type                string   `json:"type,omitempty"`
	Size                int      `json:"size,omitempty"`
	AvailabilityZone    string   `json:"availabilityZone,omitempty"`
	Image               string   `json:"image,omitempty"`
	ImageAlias          string   `json:"imageAlias,omitempty"`
	ImagePassword       string   `json:"imagePassword,omitempty"`
	SSHKeys             []string `json:"sshKeys,omitempty"`
	Bus                 string   `json:"bus,omitempty"`
	LicenceType         string   `json:"licenceType,omitempty"`
	CPUHotPlug          bool     `json:"cpuHotPlug,omitempty"`
	CPUHotUnplug        bool     `json:"cpuHotUnplug,omitempty"`
	RAMHotPlug          bool     `json:"ramHotPlug,omitempty"`
	RAMHotUnplug        bool     `json:"ramHotUnplug,omitempty"`
	NicHotPlug          bool     `json:"nicHotPlug,omitempty"`
	NicHotUnplug        bool     `json:"nicHotUnplug,omitempty"`
	DiscVirtioHotPlug   bool     `json:"discVirtioHotPlug,omitempty"`
	DiscVirtioHotUnplug bool     `json:"discVirtioHotUnplug,omitempty"`
	DiscScsiHotPlug     bool     `json:"discScsiHotPlug,omitempty"`
	DiscScsiHotUnplug   bool     `json:"discScsiHotUnplug,omitempty"`
	DeviceNumber        int64    `json:"deviceNumber,omitempty"`
}

// Volumes object
type Volumes struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Volume     `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// ListVolumes returns a Collection struct for volumes in the Datacenter
func (c *Client) ListVolumes(dcid string) (*Volumes, error) {
	ret := &Volumes{}
	return ret, c.GetOK(volumesPath(dcid), ret)
}

// GetVolume gets a volume
func (c *Client) GetVolume(dcid string, volumeID string) (*Volume, error) {
	ret := &Volume{}
	return ret, c.GetOK(volumePath(dcid, volumeID), ret)
}

// UpdateVolume updates a volume
func (c *Client) UpdateVolume(dcid string, volid string, request VolumeProperties) (*Volume, error) {
	ret := &Volume{}
	return ret, c.PatchAcc(volumePath(dcid, volid), request, ret)
}

// CreateVolume creates a volume
func (c *Client) CreateVolume(dcid string, request Volume) (*Volume, error) {
	ret := &Volume{}
	return ret, c.PostAcc(volumesPath(dcid), request, ret)
}

// DeleteVolume deletes a volume
func (c *Client) DeleteVolume(dcid, volid string) (*http.Header, error) {
	return c.DeleteAcc(volumePath(dcid, volid))
}

// CreateSnapshot creates a volume snapshot
func (c *Client) CreateSnapshot(dcid string, volid string, name string, description string) (*Snapshot, error) {
	ret := &Snapshot{}
	req := c.Client.R().
		SetFormData(map[string]string{"name": name, "description": description}).
		SetResult(ret)
	return ret, c.DoWithRequest(req, resty.MethodPost, createSnapshotPath(dcid, volid), http.StatusAccepted)
}

// RestoreSnapshot restores a volume with provided snapshot
func (c *Client) RestoreSnapshot(dcid string, volid string, snapshotID string) (*http.Header, error) {
	ret := &Header{}
	req := c.Client.R().
		SetFormData(map[string]string{"snapshotId": snapshotID}).
		SetResult(ret)
	err := c.DoWithRequest(req, resty.MethodPost, restoreSnapshotPath(dcid, volid), http.StatusAccepted)
	return ret.GetHeader(), err
}

// CreateVolumeAndWait creates a volume and waits for the request to complete.
func (c *Client) CreateVolumeAndWait(ctx context.Context, dcid string, request Volume) (*Volume, error) {
	volume, err := c.CreateVolume(dcid, request)
	if err != nil {
		return nil, err
	}
	if err := c.WaitTillProvisionedOrCanceled(ctx, volume.Headers.Get("location")); err != nil {
		return volume, err
	}
	return c.GetVolume(dcid, volume.ID)
}

// CreateSnapshotAndWait creates a volume snapshot and waits for the request to
// complete.
func (c *Client) CreateSnapshotAndWait(ctx context.Context, dcId, volId, name, description string) (*Snapshot, error) {
	snapshot, err := c.CreateSnapshot(dcId, volId, name, description)
	if err != nil {
		return nil, err
	}
	if err := c.WaitTillProvisionedOrCanceled(ctx, snapshot.Headers.Get("location")); err != nil {
		return snapshot, err
	}
	return c.GetSnapshot(snapshot.ID)
}

// RestoreSnapshotAndWait restores a volume with the provided snapshot and
// waits for the request to complete.
func (c *Client) RestoreSnapshotAndWait(ctx context.Context, dcId, volId, snapshotId string) error {
	ret, err := c.RestoreSnapshot(dcId, volId, snapshotId)
	if err != nil {
		return err
	}
	return c.WaitTillProvisionedOrCanceled(ctx, ret.Get("location"))
}

// IsSnapshotCreationRequested checks if an active create-snapshot request (QUEUED or RUNNING)
// exists for the given volume.
// Even though the snapshot lock prevents mulitple requests from being processed at the same time,
// this information is only visible by inspecting the request queue. This method can be used to
// verify that the snapshot lock is not held by any other client request.
func (c *Client) IsSnapshotCreationRequested(dcId, volId string) (bool, error) {
	f := NewRequestListFilter().WithUrl(createSnapshotPath(dcId, volId)).WithMethod(http.MethodPost)
	result, err := c.ListRequestsWithFilter(f.Clone().WithRequestStatus(RequestStatusQueued))
	if err != nil {
		return false, err
	}
	if len(result.Items) > 0 {
		return true, nil
	}
	result, err = c.ListRequestsWithFilter(f.Clone().WithRequestStatus(RequestStatusRunning))
	if err != nil {
		return false, err
	}
	return len(result.Items) > 0, nil
}
