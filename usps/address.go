package usps

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// ErrMissingUsername indicates credentials have not been supplied.
	ErrMissingUsername = errors.New("usps: missing username")
	// ErrAddressNotFound indicates the specified address was not found.
	ErrAddressNotFound = errors.New("usps: address not found")
)

var (
	strAddressNotFound = "ADDRESS NOT FOUND"
	strSuccess         = "SUCCESS"
)

type Address struct {
	Address1 string `xml:"Address1"`
	Address2 string `xml:"Address2"`
	City     string `xml:"City"`
	State    string `xml:"State"`
	Zip5     string `xml:"Zip5"`
	Zip4     string `xml:"Zip4"`
}

type ZipCode struct {
	Zip5 string `xml:"Zip5"`
}

type ResponseError struct {
	Number      string `xml:"Number"`
	Source      string `xml:"Source"`
	Description string `xml:"Description"`
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("usps: api error: %s", r.Description)
}

type AddressValidateResponse struct {
	Address struct {
		Address1 string         `xml:"Address1"`
		Address2 string         `xml:"Address2"`
		City     string         `xml:"City"`
		State    string         `xml:"State"`
		Zip5     string         `xml:"Zip5"`
		Zip4     string         `xml:"Zip4"`
		Error    *ResponseError `xml:"Error"`
	} `xml:"Address"`
}

type ZipCodeLookupResponse struct {
	Address struct {
		Address1 string         `xml:"Address1"`
		Address2 string         `xml:"Address2"`
		City     string         `xml:"City"`
		State    string         `xml:"State"`
		Zip5     string         `xml:"Zip5"`
		Zip4     string         `xml:"Zip4"`
		Error    *ResponseError `xml:"Error"`
	} `xml:"Address"`
}

type CityStateLookupResponse struct {
	ZipC struct {
		Zip5  string         `xml:"Zip5"`
		City  string         `xml:"City"`
		State string         `xml:"State"`
		Error *ResponseError `xml:"Error"`
	} `xml:"ZipCode"`
}

func (c *Client) AddressVerification(address Address) (*AddressValidateResponse, error) {
	result := &AddressValidateResponse{}
	if c.Username == "" {
		return result, ErrMissingUsername
	}

	xmlOut, err := xml.Marshal(address)
	if err != nil {
		return nil, err
	}

	var requestURL bytes.Buffer
	requestURL.WriteString("Verify&XML=")
	urlToEncode := "<AddressValidateRequest USERID=\"" + c.Username + "\">"
	urlToEncode += string(xmlOut)
	urlToEncode += "</AddressValidateRequest>"
	requestURL.WriteString(url.QueryEscape(urlToEncode))

	fmt.Fprintln(os.Stderr, string(xmlOut))
	fmt.Fprintln(os.Stderr, requestURL.String())

	body, err := c.getHTTP(requestURL.String())
	if body == nil {
		return nil, err
	}

	fmt.Fprintln(os.Stderr, "body:", string(body))

	bodyHeaderless := strings.Replace(string(body), xml.Header, "", 1)
	err = xml.Unmarshal([]byte(bodyHeaderless), &result)
	if err != nil {
		return nil, err
	}

	return result, result.Address.Error
}

func (c *Client) ZipCodeLookup(address Address) (*ZipCodeLookupResponse, error) {
	result := &ZipCodeLookupResponse{}
	if c.Username == "" {
		return result, ErrMissingUsername
	}

	xmlOut, err := xml.Marshal(address)
	if err != nil {
		return nil, err
	}

	var requestURL bytes.Buffer
	requestURL.WriteString("ZipCodeLookup&XML=")
	urlToEncode := "<ZipCodeLookupRequest USERID=\"" + c.Username + "\">"
	urlToEncode += string(xmlOut)
	urlToEncode += "</ZipCodeLookupRequest>"
	requestURL.WriteString(url.QueryEscape(urlToEncode))

	body, err := c.getHTTP(requestURL.String())
	if err != nil {
		return nil, err
	}

	bodyHeaderless := strings.Replace(string(body), xml.Header, "", 1)
	err = xml.Unmarshal([]byte(bodyHeaderless), &result)
	if err != nil {
		return nil, err
	}
	return result, result.Address.Error
}

func (c *Client) CityStateLookup(zipcode ZipCode) (*CityStateLookupResponse, error) {
	result := &CityStateLookupResponse{}
	if c.Username == "" {
		return result, ErrMissingUsername
	}
	xmlOut, err := xml.Marshal(zipcode)
	if err != nil {
		return nil, err
	}

	var requestURL bytes.Buffer
	requestURL.WriteString("CityStateLookup&XML=")
	urlToEncode := "<CityStateLookupRequest USERID=\"" + c.Username + "\">"
	urlToEncode += string(xmlOut)
	urlToEncode += "</CityStateLookupRequest>"
	requestURL.WriteString(url.QueryEscape(urlToEncode))

	body, err := c.getHTTP(requestURL.String())
	if err != nil {
		return nil, err
	}

	bodyHeaderless := strings.Replace(string(body), xml.Header, "", 1)
	err = xml.Unmarshal([]byte(bodyHeaderless), &result)
	if err != nil {
		return nil, err
	}

	return result, result.ZipC.Error
}

type AddressResponse struct {
	AddressLine1    string `json:"addressLine1,omitempty"`
	CarrierRoute    string `json:"carrierRoute,omitempty"`
	CheckDigit      string `json:"checkDigit,omitempty"`
	City            string `json:"city,omitempty"`
	IsCMAR          string `json:"cmar,omitempty"`
	CountyName      string `json:"countyName,omitempty"`
	DefaultFlag     string `json:"defaultFlag,omitempty"`
	DefaultInd      string `json:"defaultInd,omitempty"`
	DeliveryPoint   string `json:"deliveryPoint,omitempty"`
	DpvConfirmation string `json:"dpvConfirmation,omitempty"`
	Elot            string `json:"elot,omitempty"`
	ElotIndicator   string `json:"elotIndicator,omitempty"`
	RecordType      string `json:"recordType,omitempty"`
	State           string `json:"state,omitempty"`
	Zip4            string `json:"zip4,omitempty"`
	Zip5            string `json:"zip5,omitempty"`
}

type ZipByAddressResponse struct {
	AddressList  []AddressResponse `json:"addressList,omitempty"`
	ResultStatus string            `json:"resultStatus,omitempty"`
}

func (z ZipByAddressResponse) Address() AddressResponse {
	if len(z.AddressList) == 0 {
		return AddressResponse{}
	}
	return z.AddressList[0]
}

func (c *Client) ZipByAddress(address Address) (*ZipByAddressResponse, error) {
	result := &ZipByAddressResponse{}
	if c.Username == "" {
		return nil, ErrMissingUsername
	}
	params := url.Values{}
	params.Add("address1", address.Address1)
	params.Add("address2", address.Address2)
	params.Add("city", address.City)
	params.Add("state", address.State)
	params.Add("zip", address.Zip5)
	req, err := http.NewRequest("POST", "https://tools.usps.com/tools/app/ziplookup/zipByAddress", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, err := c.getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ResultStatus == strAddressNotFound {
		return result, ErrAddressNotFound
	}
	if result.ResultStatus != strSuccess {
		return result, fmt.Errorf("usps: unexpected result status: '%s'", result.ResultStatus)
	}
	return result, nil
}
