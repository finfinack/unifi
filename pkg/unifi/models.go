package unifi

import (
	"net"
	"time"
)

type Info struct {
	ApplicationVersion string `json:"applicationVersion"`
}

type ListSiteResponse struct {
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
	Count      int     `json:"count"`
	TotalCount int     `json:"totalCount"`
	Data       []*Site `json:"data"`
}

type Site struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListClientsResponse struct {
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	Count      int       `json:"count"`
	TotalCount int       `json:"totalCount"`
	Data       []*Client `json:"data"`
}

type Client struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ConnectedAt time.Time `json:"connectedAt"`
	IPAddress   net.IP    `json:"ipAddress"`
}

type ListDevicesResponse struct {
	Offset     int       `json:"offset"`
	Limit      int       `json:"limit"`
	Count      int       `json:"count"`
	TotalCount int       `json:"totalCount"`
	Data       []*Device `json:"data"`
}

type Device struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Model      string   `json:"model"`
	MACAddress string   `json:"macAddress"`
	IPAddress  net.IP   `json:"ipAddress"`
	State      string   `json:"state"`
	Features   []string `json:"features"`
	Interfaces []string `json:"interfaces"`
}

type FullDevice struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Model             string            `json:"model"`
	Supported         bool              `json:"supported"`
	MACAddress        string            `json:"macAddress"`
	IPAddress         net.IP            `json:"ipAddress"`
	State             string            `json:"state"`
	FirmwareVersion   string            `json:"firmwareVersion"`
	FirmwareUpdatable bool              `json:"firmwareUpdatable"`
	AdoptedAt         time.Time         `json:"adoptedAt"`
	ProvisionedAt     time.Time         `json:"provisionedAt"`
	ConfigurationID   string            `json:"configurationId"`
	Uplink            *DeviceUplink     `json:"uplink"`
	Features          *DeviceFeatures   `json:"features"`
	Interfaces        *DeviceInterfaces `json:"interfaces"`
}

type DeviceUplink struct {
	DeviceID string `json:"deviceId"`
}

type DeviceFeatures struct {
	Switching   interface{} `json:"switching"`
	AccessPoint interface{} `json:"accessPoint"`
}

type DeviceInterfaces struct {
	Ports  []*DevicePort  `json:"ports"`
	Radios []*DeviceRadio `json:"radios"`
}

type DevicePort struct {
	Index        int    `json:"idx"`
	State        string `json:"state"`
	Connector    string `json:"connector"`
	MaxSpeedMbps int    `json:"maxSpeedMbps"`
	SpeedMbps    int    `json:"speedMbps"`
}

type DeviceRadio struct {
	WLANStandard    string  `json:"wlanStandard"`
	FrequencyGHz    float32 `json:"frequencyGHz"`
	ChannelWidthMHz int     `json:"channelWidthMHz"`
	Channel         int     `json:"channel"`
}
