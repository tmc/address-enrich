package usps

import (
	"io/ioutil"
	"net/http"
)

const (
	baseURLInsecureTest string = "http://production.shippingapis.com/ShippingAPITest.dll?API="
	baseURLInsecure     string = "http://production.shippingapis.com/ShippingAPI.dll?API="
	baseURLTest         string = "https://secure.shippingapis.com/ShippingAPITest.dll?API="
	baseURL             string = "https://secure.shippingapis.com/ShippingAPI.dll?API="
)

// Client provides access to USPS APIs.
type Client struct {
	// Username for USPS API.
	Username string
	// Password for USPS API.
	Password string
	// Flag to determine test/prod URLs.
	Production bool `default:"false"`

	// HTTP Client to use (uses http.DefaultClient if not set).
	HTTPClient *http.Client
}

func (c *Client) getHTTPClient() *http.Client {
	if c.HTTPClient == nil {
		return http.DefaultClient
	}
	return c.HTTPClient
}

func (c *Client) getHTTP(requestURL string) ([]byte, error) {
	currentURL := ""
	if c.Production {
		currentURL += baseURL
	} else {
		currentURL += baseURLTest
	}
	currentURL += requestURL

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Get(currentURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
