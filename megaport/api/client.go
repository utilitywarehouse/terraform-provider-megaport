package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	EndpointProduction = "https://api.megaport.com"
	EndpointStaging    = "https://api-staging.megaport.com"

	Version   = "0.1"
	UserAgent = "megaport-api-go-client/" + Version

	ProductStatusDecommissioned = "DECOMMISSIONED"
)

var (
	ErrNotFound = fmt.Errorf("megaport-api: not found")
)

type Client struct {
	c         *http.Client
	BaseURL   string
	Token     string
	UserAgent string

	Port *PortService
}

func NewClient(baseURL string) *Client {
	c := &Client{c: &http.Client{}, BaseURL: baseURL}
	c.Port = NewPortService(c)
	return c
}

func (c *Client) Login(username, password, otp string) error {
	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)
	if otp != "" {
		v.Set("oneTimePassword", otp)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/login?%s", c.BaseURL, v.Encode()), nil)
	if err != nil {
		return err
	}
	data := responseLoginData{}
	if err := c.do(req, &data); err != nil {
		return err
	}
	c.Token = data.Token
	return nil
}

func (c *Client) Logout() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/logout", c.BaseURL), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

func (c *Client) GetLocations() ([]*Location, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/locations", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []*Location{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetMegaports() ([]*Megaport, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/dropdowns/partner/megaports", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []*Megaport{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetInternetExchanges(locationId uint64) ([]InternetExchange, error) {
	v := url.Values{}
	v.Set("locationId", strconv.FormatUint(locationId, 10))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/product/ix/types?%s", c.BaseURL, v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	data := []InternetExchange{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetMegaportPrice(locationId, speed, term uint64, productUid string, buyoutPort bool) (*Charges, error) {
	v := url.Values{}
	v.Set("locationId", strconv.FormatUint(locationId, 10))
	v.Set("speed", strconv.FormatUint(speed, 10))
	v.Set("term", strconv.FormatUint(term, 10))
	v.Set("buyoutPort", strconv.FormatBool(buyoutPort))
	if productUid != "" {
		v.Set("productUid", productUid) // TODO: can we just set to empty?
	}
	return c.getCharges("megaport", v)
}

func (c *Client) GetMCR1Price(locationId, speed uint64, productUid string) (*Charges, error) {
	v := url.Values{}
	v.Set("locationId", strconv.FormatUint(locationId, 10))
	v.Set("speed", strconv.FormatUint(speed, 10))
	if productUid != "" {
		v.Set("productUid", productUid) // TODO: can we just set to empty?
	}
	return c.getCharges("mcr", v)
}

func (c *Client) GetMCR2Price(locationId, speed uint64, productUid string) (*Charges, error) {
	v := url.Values{}
	v.Set("locationId", strconv.FormatUint(locationId, 10))
	v.Set("speed", strconv.FormatUint(speed, 10))
	if productUid != "" {
		v.Set("productUid", productUid) // TODO: can we just set to empty?
	}
	return c.getCharges("mcr2", v)
}

func (c *Client) GetVxcPrice(aLocationId, bLocationId, speed uint64) (*Charges, error) {
	v := url.Values{}
	v.Set("aLocationId", strconv.FormatUint(aLocationId, 10))
	v.Set("bLocationId", strconv.FormatUint(bLocationId, 10))
	v.Set("speed", strconv.FormatUint(speed, 10))
	return c.getCharges("vxc", v)
}

func (c *Client) GetIxPrice(ixType string, locationId, speed uint64) (*Charges, error) {
	v := url.Values{}
	v.Set("ixType", ixType)
	v.Set("portLocationId", strconv.FormatUint(locationId, 10))
	v.Set("speed", strconv.FormatUint(speed, 10))
	return c.getCharges("ix", v)
}

func (c *Client) getCharges(product string, v url.Values) (*Charges, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/pricebook/%s?%s", c.BaseURL, product, v.Encode()), nil)
	if err != nil {
		return nil, err
	}
	data := &Charges{}
	if err := c.do(req, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) do(req *http.Request, data interface{}) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		r := megaportResponse{}
		if err := parseResponseBody(resp, &r); err != nil {
			return err
		}
		err := responseDataToError(r.Data)
		if err != nil {
			return fmt.Errorf("megaport-api: %s: %w", r.Message, err)
		} else {
			return fmt.Errorf("megaport-api: %s", r.Message)
		}
	}
	return parseResponseBody(resp, &megaportResponse{Data: data})
}

func responseDataToError(d interface{}) error {
	switch e := d.(type) {
	case string:
		return fmt.Errorf("%s", e)
	case map[string]interface{}:
		errData := &strings.Builder{}
		for k, v := range e {
			if _, err := fmt.Fprintf(errData, "%s=%#v ", k, v); err != nil {
				return err
			}
		}
		return fmt.Errorf("%s", strings.TrimSpace(errData.String()))
	case []interface{}:
		errors := make([]string, len(e))
		for i, v := range e {
			errors[i] = responseDataToError(v).Error()
		}
		return fmt.Errorf("%d errors: ['%s']", len(e), strings.Join(errors, "', '"))
	case nil:
		return nil
	default:
		return fmt.Errorf("cannot process error data of type %T: %#v", e, e)
	}
}

func parseResponseBody(resp *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	return nil
}

func (c *Client) userAgent() string {
	if c.UserAgent == "" {
		return UserAgent
	}
	return UserAgent + " " + c.UserAgent
}
