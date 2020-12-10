package profitbricks

import (
	"net/http"
)

// Groups object
type Groups struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Group      `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// Group object
type Group struct {
	ID         string          `json:"id,omitempty"`
	PBType     string          `json:"type,omitempty"`
	Href       string          `json:"href,omitempty"`
	Properties GroupProperties `json:"properties,omitempty"`
	Entities   *GroupEntities  `json:"entities,omitempty"`
	Response   string          `json:"Response,omitempty"`
	Headers    *http.Header    `json:"headers,omitempty"`
	StatusCode int             `json:"statuscode,omitempty"`
}

// GroupProperties object
type GroupProperties struct {
	Name                 string `json:"name,omitempty"`
	CreateDataCenter     *bool  `json:"createDataCenter,omitempty"`
	CreateSnapshot       *bool  `json:"createSnapshot,omitempty"`
	ReserveIP            *bool  `json:"reserveIp,omitempty"`
	AccessActivityLog    *bool  `json:"accessActivityLog,omitempty"`
	CreateBackupUnit     *bool  `json:"createBackupUnit,omitempty"`
	CreateInternetAccess *bool  `json:"createInternetAccess,omitempty"`
	CreateK8sCluster     *bool  `json:"createK8sCluster,omitempty"`
	CreatePcc            *bool  `json:"createPcc,omitempty"`
	S3Privilege          *bool  `json:"s3Privilege,omitempty"`
}

// GroupEntities object
type GroupEntities struct {
	Users     Users     `json:"users,omitempty"`
	Resources Resources `json:"resources,omitempty"`
}

// Users object
type Users struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []User       `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// User object
type User struct {
	ID         string          `json:"id,omitempty"`
	PBType     string          `json:"type,omitempty"`
	Href       string          `json:"href,omitempty"`
	Metadata   *Metadata       `json:"metadata,omitempty"`
	Properties *UserProperties `json:"properties,omitempty"`
	Entities   *UserEntities   `json:"entities,omitempty"`
	Response   string          `json:"Response,omitempty"`
	Headers    *http.Header    `json:"headers,omitempty"`
	StatusCode int             `json:"statuscode,omitempty"`
}

// UserProperties object
type UserProperties struct {
	Firstname         string `json:"firstname,omitempty"`
	Lastname          string `json:"lastname,omitempty"`
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	Administrator     bool   `json:"administrator,omitempty"`
	ForceSecAuth      bool   `json:"forceSecAuth,omitempty"`
	SecAuthActive     bool   `json:"secAuthActive,omitempty"`
	Active            *bool  `json:"active,omitempty"`
	S3CanonicalUserID string `json:"s3CanonicalUserId,omitempty"`
}

// UserEntities object
type UserEntities struct {
	Groups Groups `json:"groups,omitempty"`
	Owns   Owns   `json:"owns,omitempty"`
}

// Resources object
type Resources struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Resource   `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// Resource object
type Resource struct {
	ID         string            `json:"id,omitempty"`
	PBType     string            `json:"type,omitempty"`
	Href       string            `json:"href,omitempty"`
	Metadata   *Metadata         `json:"metadata,omitempty"`
	Entities   *ResourceEntities `json:"entities,omitempty"`
	Response   string            `json:"Response,omitempty"`
	Headers    *http.Header      `json:"headers,omitempty"`
	StatusCode int               `json:"statuscode,omitempty"`
}

// ResourceEntities object
type ResourceEntities struct {
	Groups Groups `json:"groups,omitempty"`
}

// Owns object
type Owns struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Entity     `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// Entity object
type Entity struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Metadata   *Metadata    `json:"metadata,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// Shares object
type Shares struct {
	ID         string       `json:"id,omitempty"`
	PBType     string       `json:"type,omitempty"`
	Href       string       `json:"href,omitempty"`
	Items      []Share      `json:"items,omitempty"`
	Response   string       `json:"Response,omitempty"`
	Headers    *http.Header `json:"headers,omitempty"`
	StatusCode int          `json:"statuscode,omitempty"`
}

// Share object
type Share struct {
	ID         string          `json:"id,omitempty"`
	PBType     string          `json:"type,omitempty"`
	Href       string          `json:"href,omitempty"`
	Properties ShareProperties `json:"properties,omitempty"`
	Response   string          `json:"Response,omitempty"`
	Headers    *http.Header    `json:"headers,omitempty"`
	StatusCode int             `json:"statuscode,omitempty"`
}

// ShareProperties object
type ShareProperties struct {
	EditPrivilege  *bool `json:"editPrivilege,omitempty"`
	SharePrivilege *bool `json:"sharePrivilege,omitempty"`
}

//ListGroups lists all groups
func (c *Client) ListGroups() (*Groups, error) {
	url := groupsPath()
	ret := &Groups{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//GetGroup gets a group
func (c *Client) GetGroup(groupid string) (*Group, error) {
	url := groupPath(groupid)
	ret := &Group{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//CreateGroup creates a group
func (c *Client) CreateGroup(grp Group) (*Group, error) {
	url := groupsPath()
	ret := &Group{}
	err := c.Post(url, grp, ret, http.StatusAccepted)
	return ret, err
}

//UpdateGroup updates a group
func (c *Client) UpdateGroup(groupid string, obj Group) (*Group, error) {
	url := groupPath(groupid)
	ret := &Group{}
	err := c.Put(url, obj, ret, http.StatusAccepted)
	return ret, err
}

//DeleteGroup deletes a group
func (c *Client) DeleteGroup(groupid string) (*http.Header, error) {
	url := groupPath(groupid)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}

//ListShares lists all shares
func (c *Client) ListShares(grpid string) (*Shares, error) {
	url := sharesPath(grpid)
	ret := &Shares{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// GetShare gets a share
func (c *Client) GetShare(groupid string, resourceid string) (*Share, error) {
	url := sharePath(groupid, resourceid)
	ret := &Share{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// AddShare adds a share
func (c *Client) AddShare(groupid string, resourceid string, share Share) (*Share, error) {
	url := sharePath(groupid, resourceid)
	ret := &Share{}
	err := c.Post(url, share, ret, http.StatusAccepted)
	return ret, err
}

// UpdateShare updates a share
func (c *Client) UpdateShare(groupid string, resourceid string, obj Share) (*Share, error) {
	url := sharePath(groupid, resourceid)
	ret := &Share{}
	err := c.Put(url, obj, ret, http.StatusAccepted)
	return ret, err
}

// DeleteShare deletes a share
func (c *Client) DeleteShare(groupid string, resourceid string) (*http.Header, error) {
	url := sharePath(groupid, resourceid)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}

//ListGroupUsers lists Users in a group
func (c *Client) ListGroupUsers(groupid string) (*Users, error) {
	url := groupUsersPath(groupid)
	ret := &Users{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// AddUserToGroup adds a user to a group
func (c *Client) AddUserToGroup(groupid string, userid string) (*User, error) {
	var usr User
	usr.ID = userid
	url := groupUsersPath(groupid)
	ret := &User{}
	err := c.Post(url, usr, ret, http.StatusAccepted)
	return ret, err
}

// DeleteUserFromGroup removes a user from a group
func (c *Client) DeleteUserFromGroup(groupid string, userid string) (*http.Header, error) {
	url := groupUserPath(groupid, userid)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}

//ListUsers lists all users
func (c *Client) ListUsers() (*Users, error) {
	url := usersPath()
	ret := &Users{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

// GetUser gets a user
func (c *Client) GetUser(usrid string) (*User, error) {
	url := userPath(usrid)
	ret := &User{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//CreateUser creates a user
func (c *Client) CreateUser(usr User) (*User, error) {
	url := usersPath()
	ret := &User{}
	err := c.Post(url, usr, ret, http.StatusAccepted)
	return ret, err
}

//UpdateUser updates user information
func (c *Client) UpdateUser(userid string, obj User) (*User, error) {
	url := userPath(userid)
	ret := &User{}
	err := c.Put(url, obj, ret, http.StatusAccepted)
	return ret, err
}

//DeleteUser deletes the specified user
func (c *Client) DeleteUser(userid string) (*http.Header, error) {
	url := userPath(userid)
	ret := &http.Header{}
	err := c.Delete(url, ret, http.StatusAccepted)
	return ret, err
}

//ListResources lists all resources
func (c *Client) ListResources() (*Resources, error) {
	url := resourcesPath()
	ret := &Resources{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//GetResourceByType gets a resource by type
func (c *Client) GetResourceByType(resourcetype string, resourceid string) (*Resource, error) {
	url := resourcePath(resourcetype, resourceid)
	ret := &Resource{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}

//ListResourcesByType list resources by type
func (c *Client) ListResourcesByType(resourcetype string) (*Resources, error) {
	url := resourcesTypePath(resourcetype)
	ret := &Resources{}
	err := c.Get(url, ret, http.StatusOK)
	return ret, err
}
