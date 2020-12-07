package profitbricks

import (
	"fmt"
	"net/http"
)

//Snapshot object
type Snapshot struct {
	ID         string             `json:"id,omitempty"`
	PBType     string             `json:"type,omitempty"`
	Href       string             `json:"href,omitempty"`
	Metadata   Metadata           `json:"metadata,omitempty"`
	Properties SnapshotProperties `json:"properties,omitempty"`
	Response   string             `json:"Response,omitempty"`
	Headers    *http.Header       `json:"headers,omitempty"`
	StatusCode int                `json:"statuscode,omitempty"`
}

// SnapshotProperties properties
type SnapshotProperties struct {
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	Location            string `json:"location,omitempty"`
	Size                int    `json:"size,omitempty"`
	CPUHotPlug          bool   `json:"cpuHotPlug,omitempty"`
	CPUHotUnplug        bool   `json:"cpuHotUnplug,omitempty"`
	RAMHotPlug          bool   `json:"ramHotPlug,omitempty"`
	RAMHotUnplug        bool   `json:"ramHotUnplug,omitempty"`
	NicHotPlug          bool   `json:"nicHotPlug,omitempty"`
	NicHotUnplug        bool   `json:"nicHotUnplug,omitempty"`
	DiscVirtioHotPlug   bool   `json:"discVirtioHotPlug,omitempty"`
	DiscVirtioHotUnplug bool   `json:"discVirtioHotUnplug,omitempty"`
	DiscScsiHotPlug     bool   `json:"discScsiHotPlug,omitempty"`
	DiscScsiHotUnplug   bool   `json:"discScsiHotUnplug,omitempty"`
	LicenceType         string `json:"licenceType,omitempty"`
}

//Snapshots object
type Snapshots struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Snapshot   `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

//ListSnapshots lists all snapshots
func (c *Client) ListSnapshots() (*Snapshots, error) {
	url := snapshotsPath()
	ret := &Snapshots{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//GetSnapshot gets a specific snapshot
func (c *Client) GetSnapshot(snapshotID string) (*Snapshot, error) {
	url := snapshotPath(snapshotID)
	ret := &Snapshot{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// DeleteSnapshot deletes a specified snapshot
func (c *Client) DeleteSnapshot(snapshotID string) (*http.Header, error) {
	url := snapshotPath(snapshotID)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}

// UpdateSnapshot updates a snapshot
func (c *Client) UpdateSnapshot(snapshotID string, request SnapshotProperties) (*Snapshot, error) {
	url := snapshotPath(snapshotID)
	ret := &Snapshot{}
	err := c.Patch(url, request, ret, http.StatusAccepted)
	return ret, err
}

// DeleteSnapshotAndWait deletes a specified snapshot and waits for the request
// to complete.
func (c *Client) DeleteSnapshotAndWait(snapshotID string) error {
	ret, err := c.DeleteSnapshot(snapshotID)
	if err != nil {
		return err
	}
	return c.WaitTillProvisioned(ret.Get("location"))
}

// ListSnapshotsWithSelector retrieves all snapshots and performs client-side
// filtering according to the selector.
func (c *Client) ListSnapshotsWithSelector(selector SnapshotSelector) ([]Snapshot, error) {
	if selector == nil {
		return nil, fmt.Errorf("missing selector")
	}

	snapshots, err := c.ListSnapshots()
	if err != nil {
		return nil, err
	}

	var result []Snapshot
	for _, snapshot := range snapshots.Items {
		if !selector(&snapshot) {
			continue
		}
		result = append(result, snapshot)
	}
	return result, nil
}

// SnapshotSelector is used to do client-side filtering of a list of Snapshots
type SnapshotSelector func(*Snapshot) bool

// SnapshotByState selects snapshots with the given state
func SnapshotByState(state string) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		return snapshot.Metadata.State == state
	}
}

// SnapshotByName selects snapshots with the given name
func SnapshotByName(name string) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		return snapshot.Properties.Name == name
	}
}

// SnapshotByDescription selects snapshots with the given description
func SnapshotByDescription(description string) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		return snapshot.Properties.Description == description
	}
}

// SelectExactSnapshot concatenates the provided selectors with logical AND.
func SelectExactSnapshot(selectors ...SnapshotSelector) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		for _, selector := range selectors {
			if !selector(snapshot) {
				return false
			}
		}
		return true
	}
}

// SelectAnySnapshot concatenates the provided selectors with logical OR.
func SelectAnySnapshot(selectors ...SnapshotSelector) SnapshotSelector {
	return func(snapshot *Snapshot) bool {
		for _, selector := range selectors {
			if selector(snapshot) {
				return true
			}
		}
		return false
	}
}

// IsSnapshotDeletionRequested checks if there is any active delete snapshot request
// (QUEUED or RUNNING).
func (c *Client) IsSnapshotDeletionRequested(snapshotID string) (bool, error) {
	f := NewRequestListFilter().WithUrl(snapshotPath(snapshotID)).WithMethod(http.MethodDelete)
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

// IsSnapshotInUseError returns true if the given error indicates that a snapshot
// is being referenced by a volume that is currently being created.
func IsSnapshotInUseError(err error) bool {
	apiErr, ok := err.(ApiError)
	if !ok {
		return false
	}
	return apiErr.HTTPStatus == http.StatusUnprocessableEntity &&
		apiErr.HasErrorCode(SnapshotInUseErrorCode)
}
