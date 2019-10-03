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
)

type Client struct {
	c       *http.Client
	baseURL string
	token   string
}

func NewClient(baseURL string) *Client {
	return &Client{
		c:       &http.Client{},
		baseURL: baseURL,
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) Login(username, password string) error {
	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)
	// v.Set("oneTimePassword", "")
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/login?%s", c.baseURL, v.Encode()), nil)
	if err != nil {
		return err
	}
	data := responseLoginData{}
	if err := c.do(req, &data); err != nil {
		return err
	}
	c.token = data.Token
	return nil
}

func (c *Client) Logout() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/logout", c.baseURL), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

func (c *Client) GetLocations() ([]*Location, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/locations", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []*Location{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetMegaports() ([]Megaport, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/dropdowns/partner/megaports", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []Megaport{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) GetInternetExchanges(locationId uint64) ([]InternetExchange, error) {
	v := url.Values{}
	v.Set("locationId", strconv.FormatUint(locationId, 10))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/product/ix/types?%s", c.baseURL, v.Encode()), nil)
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/pricebook/%s?%s", c.baseURL, product, v.Encode()), nil)
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
	if c.token != "" {
		req.Header.Set("X-Auth-Token", c.token)
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	r := megaportResponse{}
	if resp.StatusCode != http.StatusOK {
		r.Data = map[string]interface{}{}
		if err = parseResponseBody(resp, &r); err != nil {
			return err
		}
		errData := &strings.Builder{}
		for k, v := range r.Data.(map[string]interface{}) {
			if _, err := fmt.Fprintf(errData, "%s=%#v ", k, v); err != nil {
				return err
			}
		}
		return fmt.Errorf("megaport-api: %s (%s)", r.Message, strings.TrimSpace(errData.String()))
	}
	r.Data = data
	return parseResponseBody(resp, &r)
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
