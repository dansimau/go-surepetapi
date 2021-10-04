package surepetapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dansimau/go-surepetapi/pkg/httpclient"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const (
	apiURL = "https://app.api.surehub.io/api"
)

type Client struct {
	cfg  *Config
	http *httpclient.Client
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.AuthToken == "" && (cfg.EmailAddress == "" || cfg.Password == "") {
		return nil, errors.New("missing auth token, or email/password")
	}
	return &Client{cfg: &cfg}, nil
}

func (c *Client) request(params httpclient.RequestParams) (*http.Response, error) {
	if c.http == nil {
		c.http = httpclient.New(
			httpclient.WithDebugLogging(false),
		)
	}

	if c.http.BaseURL == "" {
		c.http.BaseURL = apiURL
	}

	return c.http.Request(params)
}

func (c *Client) authenticatedRequest(params httpclient.RequestParams) (*http.Response, error) {
	if c.cfg.AuthToken == "" {
		deviceID, err := uuid.NewV1()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate device ID")
		}

		authRes, _, err := c.AuthLogin(AuthLoginRequest{
			EmailAddress: c.cfg.EmailAddress,
			Password:     c.cfg.Password,
			DeviceID:     deviceID.String(),
		})
		if err != nil {
			return nil, err
		}

		c.cfg.AuthToken = authRes.Data.Token
	}

	if params.Headers == nil {
		params.Headers = map[string]string{}
	}
	params.Headers["Authorization"] = fmt.Sprintf("Bearer %s", c.cfg.AuthToken)

	httpRes, err := c.request(params)
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != 200 {
		return httpRes, errors.New("non-200 status code")
	}

	return httpRes, nil
}

type AuthLoginRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
	DeviceID     string `json:"device_id"`
}

type AuthLoginResponse struct {
	Data struct {
		User  User
		Token string
	}
}

func (c *Client) AuthLogin(req AuthLoginRequest) (res *AuthLoginResponse, httpRes *http.Response, err error) {
	httpRes, err = c.request(httpclient.RequestParams{
		Method:   "POST",
		URL:      "/auth/login",
		BodyJSON: req,
	})
	if err != nil {
		return nil, httpRes, err
	}

	b, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, httpRes, errors.Wrap(err, "error reading body")
	}

	if err := json.Unmarshal(b, &res); err != nil {
		return nil, httpRes, err
	}

	return res, httpRes, err
}

type DevicesResponse struct {
	Data []Device
}

func (c *Client) ListDevices() (res *DevicesResponse, httpRes *http.Response, err error) {
	httpRes, err = c.authenticatedRequest(httpclient.RequestParams{
		Method: "GET",
		URL:    "/device?with%5B%5D=children&with%5B%5D=tags&with%5B%5D=control&with%5B%5D=status",
	})
	if err != nil {
		return nil, httpRes, err
	}

	b, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, httpRes, errors.Wrap(err, "error reading body")
	}

	if err := json.Unmarshal(b, &res); err != nil {
		return nil, httpRes, err
	}

	return res, httpRes, err
}

func (c *Client) SetLockState(deviceID string, lockState LockState) (res *DeviceControl, httpRes *http.Response, err error) {
	httpRes, err = c.authenticatedRequest(httpclient.RequestParams{
		Method: "PUT",
		URL:    fmt.Sprintf("/device/%s/control", deviceID),
		BodyJSON: struct {
			Locking LockState `json:"locking"`
		}{
			Locking: lockState,
		},
	})
	if err != nil {
		return nil, httpRes, err
	}

	b, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return nil, httpRes, errors.Wrap(err, "error reading body")
	}

	if err := json.Unmarshal(b, &res); err != nil {
		return nil, httpRes, err
	}

	return res, httpRes, err
}
