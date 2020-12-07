package profitbricks

import (
	"net/http"
	"strconv"
	"time"

	resty "github.com/go-resty/resty/v2"
)

type Client struct {
	*resty.Client
	// AuthApiUrl will be used by methods talking to the auth api by sending absolute urls
	AuthApiUrl  string
	CloudApiUrl string
}

const (
	DefaultApiUrl  = "https://api.ionos.com/cloudapi/v5"
	DefaultAuthUrl = "https://api.ionos.com/auth/v1"
	Version        = "5.0.3"
)

func RestyClient(username, password, token string) *Client {
	c := &Client{
		Client:      resty.New(),
		AuthApiUrl:  DefaultAuthUrl,
		CloudApiUrl: DefaultApiUrl,
	}
	if token == "" {
		c.SetBasicAuth(username, password)
	} else {
		c.SetAuthToken(token)
	}
	c.SetHostURL(DefaultApiUrl)
	c.SetDepth(10)
	c.SetTimeout(3 * time.Minute)
	c.SetUserAgent("ionos-enterprise-sdk-go " + Version)
	c.SetRetryCount(3)
	c.SetRetryMaxWaitTime(10 * time.Minute)
	c.SetRetryWaitTime(1 * time.Second)
	c.SetRetryAfter(func(cl *resty.Client, r *resty.Response) (time.Duration, error) {
		switch r.StatusCode() {
		case http.StatusTooManyRequests:
			if retryAfterSeconds := r.Header().Get("Retry-After"); retryAfterSeconds != "" {
				return time.ParseDuration(retryAfterSeconds + "s")
			}
		}
		return cl.RetryWaitTime, nil
	})
	c.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			switch r.StatusCode() {
			case http.StatusTooManyRequests,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
				http.StatusBadGateway:
				return true
			}
			return false
		})
	return c
}

// SetDebug activates/deactivates resty's debug mode. For better readability
// the pretty print feature is also enabled.
func (c *Client) SetDebug(debug bool) {
	c.Client.SetDebug(debug)
	c.SetPretty(debug)
}

// SetDepth sets the depth of information that will be retrieved by api calls. The
// API accepts values from 0 to 10, a low depth means mostly only IDs and hrefs will be
// returned. Therefore nested structures may be nil.
func (c *Client) SetDepth(depth int) {
	c.Client.SetQueryParam("depth", strconv.Itoa(depth))
}

// SetPretty toggles if the data retrieved from the api will be delivered pretty printed.
// Usually this does not make sense from an sdk perspective, but for debugging it's nice
// therefore it is also set to true, if debug is enabled.
func (c *Client) SetPretty(pretty bool) {
	c.Client.SetQueryParam("pretty", strconv.FormatBool(pretty))
}

// NewClient is a constructor for Client object
func NewClient(username, password string) *Client {
	return RestyClient(username, password, "")
}

// NewClientbyToken is a constructor for Client object using bearer tokens for
// authentication instead of username, password
func NewClientbyToken(token string) *Client {
	return RestyClient("", "", token)
}

// SetUserAgent sets User-Agent request header for all API calls
func (c *Client) SetUserAgent(agent string) {
	c.Client.SetHeader("User-Agent", agent)
}

// GetUserAgent gets User-Agent header
func (c *Client) GetUserAgent() string {
	return c.Client.Header.Get("User-Agent")
}

// SetCloudApiURL sets Cloud API url
func (c *Client) SetCloudApiURL(url string) {
	c.Client.SetHostURL(url)
}

// SetAuthApiUrl sets the Auth API url
func (c *Client) SetAuthApiUrl(url string) {
	c.AuthApiUrl = url
}
