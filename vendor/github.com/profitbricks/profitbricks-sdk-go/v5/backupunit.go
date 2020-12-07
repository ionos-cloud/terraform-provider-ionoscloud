package profitbricks

import "net/http"

// BackupUnits type
type BackupUnits struct {
	// Enum: [backupunits]
	// Read Only: true
	ID string `json:"id,omitempty"`
	// Enum: [collection]
	// Read Only: true
	Type string `json:"type,omitempty"`
	// Format: uri
	Href string `json:"href"`
	// Read Only: true
	Items []BackupUnit `json:"items"`
}

// BackupUnit Object
type BackupUnit struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// The type of object. In this case backupunit
	// Read Only: true
	Type string `json:"type,omitempty"`

	// The metadata for the backup unit
	// Read Only: true
	Metadata *Metadata `json:"metadata,omitempty"`

	Properties *BackupUnitProperties `json:"properties,omitempty"`
}

// BackupUnitProperties object
type BackupUnitProperties struct {
	// Required: on create
	// Read only: yes
	Name string `json:"name,omitempty"`
	// Required: yes
	Password string `json:"password,omitempty"`
	// Required: yes
	Email string `json:"email,omitempty"`
}

// BackupUnitSSOURL object
type BackupUnitSSOURL struct {

	// The type of object. In this case backupunit
	// Read Only: true
	Type string `json:"type,omitempty"`
	// SSO URL

	// Read Only: true
	SSOUrl string `json:"ssoURL,omitempty"`
}

//

// CreateBackupUnit creates a Backup Unit
func (c *Client) CreateBackupUnit(backupUnit BackupUnit) (*BackupUnit, error) {
	rsp := &BackupUnit{}
	return rsp, c.PostAcc(backupUnitsPath(), backupUnit, rsp)
}

// ListBackupUnits lists all available backup units
func (c *Client) ListBackupUnits() (*BackupUnits, error) {
	rsp := &BackupUnits{}
	return rsp, c.GetOK(backupUnitsPath(), rsp)
}

// UpdateBackupUnit updates an existing backup unit
func (c *Client) UpdateBackupUnit(backupUnitID string, backupUnit BackupUnit) (*BackupUnit, error) {
	rsp := &BackupUnit{}
	return rsp, c.PutAcc(backupUnitPath(backupUnitID), backupUnit, rsp)
}

// DeleteBackupUnit deletes an existing backup unit
func (c *Client) DeleteBackupUnit(backupUnitID string) (*http.Header, error) {
	return c.DeleteAcc(backupUnitPath(backupUnitID))
}

// GetBackupUnit retrieves an existing backup unit
func (c *Client) GetBackupUnit(backupUnitID string) (*BackupUnit, error) {
	rsp := &BackupUnit{}
	return rsp, c.GetOK(backupUnitPath(backupUnitID), rsp)
}

// GetBackupUnitSSOURL retrieves the SSO URL for an existing backup unit
func (c *Client) GetBackupUnitSSOURL(backupUnitID string) (*BackupUnitSSOURL, error) {
	rsp := &BackupUnitSSOURL{}
	return rsp, c.GetOK(backupUnitSSOURLPath(backupUnitID), rsp)
}
