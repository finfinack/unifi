package unifi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/finfinack/logger/logging"
)

const (
	apiBasePath              = "/proxy/network/integrations"
	infoPath                 = "/v1/info"
	sitesPath                = "/v1/sites"
	siteClientsPathTmpl      = "/v1/sites/%s/clients"
	siteDevicesPathTmpl      = "/v1/sites/%s/devices"
	siteDeviceDetailPathTmpl = "/v1/sites/%s/devices/%s"
)

type Controller struct {
	Host   string
	APIKey string

	Client *http.Client
	Logger *logging.Logger
}

func NewController(host, key string) *Controller {
	return &Controller{
		Host:   host,
		APIKey: key,
		Client: http.DefaultClient,
		Logger: logging.NewLogger("UNIF"),
	}
}

func (c *Controller) callAPI(method, path string, params url.Values) ([]byte, error) {
	u, err := url.JoinPath(c.Host, apiBasePath, path)
	if err != nil {
		return nil, fmt.Errorf("unable to create host path: %s", err)
	}

	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %s", err)
	}
	req.Header.Set("X-API-KEY", c.APIKey)
	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = params.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("site response is not OK: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %s", err)
	}
	defer resp.Body.Close()

	return body, nil
}

func (c *Controller) GetInfo() (*Info, error) {
	resp, err := c.callAPI(http.MethodGet, infoPath, nil)
	if err != nil {
		return nil, err
	}
	infoResponse := &Info{}
	if err := json.Unmarshal(resp, infoResponse); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response: %s", err)
	}
	return infoResponse, nil
}

func (c *Controller) ListSites() ([]*Site, error) {
	var sites []*Site
	offset := 0
	limit := 25
	for {
		params := url.Values{}
		params.Add("offset", strconv.Itoa(offset))
		params.Add("limit", strconv.Itoa(limit))
		resp, err := c.callAPI(http.MethodGet, sitesPath, params)
		if err != nil {
			return nil, err
		}

		siteResponse := &ListSiteResponse{}
		if err := json.Unmarshal(resp, siteResponse); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response: %s", err)
		}
		sites = append(sites, siteResponse.Data...)

		if len(sites) >= siteResponse.TotalCount {
			return sites, nil
		}
		offset += limit
	}
}

func (c *Controller) ListClients(site string) ([]*Client, error) {
	var clients []*Client
	offset := 0
	limit := 25
	for {
		params := url.Values{}
		params.Add("siteId", site)
		params.Add("offset", strconv.Itoa(offset))
		params.Add("limit", strconv.Itoa(limit))
		resp, err := c.callAPI(http.MethodGet, fmt.Sprintf(siteClientsPathTmpl, site), params)
		if err != nil {
			return nil, err
		}

		clientResponse := &ListClientsResponse{}
		if err := json.Unmarshal(resp, clientResponse); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response: %s", err)
		}
		clients = append(clients, clientResponse.Data...)

		if len(clients) >= clientResponse.TotalCount {
			return clients, nil
		}
		offset += limit
	}
}

func (c *Controller) ListDevices(site string) ([]*Device, error) {
	var devices []*Device
	offset := 0
	limit := 25
	for {
		params := url.Values{}
		params.Add("siteId", site)
		params.Add("offset", strconv.Itoa(offset))
		params.Add("limit", strconv.Itoa(limit))
		resp, err := c.callAPI(http.MethodGet, fmt.Sprintf(siteDevicesPathTmpl, site), params)
		if err != nil {
			return nil, err
		}

		devicesResponse := &ListDevicesResponse{}
		if err := json.Unmarshal(resp, devicesResponse); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response: %s", err)
		}
		devices = append(devices, devicesResponse.Data...)

		if len(devices) >= devicesResponse.TotalCount {
			return devices, nil
		}
		offset += limit
	}
}

func (c *Controller) GetDeviceDetail(site, device string) (*FullDevice, error) {
	params := url.Values{}
	params.Add("siteId", site)
	params.Add("deviceId", device)
	resp, err := c.callAPI(http.MethodGet, fmt.Sprintf(siteDeviceDetailPathTmpl, site, device), params)
	if err != nil {
		return nil, err
	}

	deviceResponse := &FullDevice{}
	if err := json.Unmarshal(resp, deviceResponse); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response: %s", err)
	}
	return deviceResponse, nil
}
