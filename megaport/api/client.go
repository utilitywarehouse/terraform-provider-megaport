package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (c *Client) GetLocations() ([]Location, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/locations", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []Location{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) do(req *http.Request, data interface{}) error {
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
		if err = parseResponseBody(resp, &r); err != nil {
			return err
		}
		return fmt.Errorf("megaport: %s", r.Message)
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
