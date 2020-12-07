package profitbricks

import "net/http"

// S3Keys type
type S3Keys struct {
	// Enum: [backupunits]
	// Read Only: true
	ID string `json:"id,omitempty"`
	// Enum: [collection]
	// Read Only: true
	Type string `json:"type,omitempty"`
	// Format: uri
	Href string `json:"href"`
	// Read Only: true
	Items []S3Key `json:"items"`
}

// S3Key Object
type S3Key struct {
	// URL to the object representation (absolute path)
	// Read Only: true
	// Format: uri
	Href string `json:"href,omitempty"`

	// The resource's unique identifier.
	// Read Only: true
	ID string `json:"id,omitempty"`

	// The type of object. In this case s3key
	// Read Only: true
	Type string `json:"type,omitempty"`

	// The metadata for the S3 key
	// Read Only: true
	Metadata *Metadata `json:"metadata,omitempty"`

	// The properties of the S3 key
	// Read Only: false
	Properties *S3KeyProperties `json:"properties,omitempty"`
}

// S3KeyProperties object
type S3KeyProperties struct {
	// Read only: yes
	SecretKey string `json:"secretKey,omitempty"`
	// Required: yes
	// Read only: no
	Active bool `json:"active"`
}

// CreateS3Key creates an S3 Key for an user
func (c *Client) CreateS3Key(userID string) (*S3Key, error) {
	rsp := &S3Key{}
	var requestBody interface{}
	err := c.Post(s3KeysPath(userID), requestBody, rsp, http.StatusCreated)
	return rsp, err
}

// ListS3Keys lists all available S3 keys for an user
func (c *Client) ListS3Keys(userID string) (*S3Keys, error) {
	rsp := &S3Keys{}
	return rsp, c.GetOK(s3KeysListPath(userID), rsp)
}

// UpdateS3Key updates an existing S3 key
func (c *Client) UpdateS3Key(userID string, s3KeyID string, s3Key S3Key) (*S3Key, error) {
	rsp := &S3Key{}
	return rsp, c.PutAcc(s3KeyPath(userID, s3KeyID), s3Key, rsp)
}

// DeleteS3Key deletes an existing S3 key
func (c *Client) DeleteS3Key(userID string, s3KeyID string) (*http.Header, error) {
	return c.DeleteAcc(s3KeyPath(userID, s3KeyID))
}

// GetS3Key retrieves an existing S3 key
func (c *Client) GetS3Key(userID string, s3KeyID string) (*S3Key, error) {
	rsp := &S3Key{}
	return rsp, c.GetOK(s3KeyPath(userID, s3KeyID), rsp)
}
