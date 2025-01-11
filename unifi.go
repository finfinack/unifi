package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/finfinack/logger/logging"
	"github.com/finfinack/unifi/pkg/unifi"
)

var (
	host    = flag.String("host", "", "Base URL for the UniFi controller")
	keyFile = flag.String("key", "key", "file containing the API key to access the UniFi controller")

	devices      = flag.Bool("devices", false, "list devices per site when set to true")
	deviceDetail = flag.Bool("deviceDetails", false, "fetch device details when set to true")
	clients      = flag.Bool("clients", false, "list clients per site when set to true")
)

func readAPIKey(file string) (string, error) {
	k, err := os.ReadFile(*keyFile)
	if err != nil {
		return "", fmt.Errorf("unable to read API key from %q: %s", file, err)
	}
	return strings.TrimSpace(string(k)), nil
}

func main() {
	flag.Parse()

	// Set up logging
	logging.SetMinLogLevel(logging.LogLevelInfo) // Info is the default level, just a demo here
	logger := logging.NewLogger("MAIN")
	logger.SetWriter(os.Stdout) // stdout is the default, just a demo here
	defer logger.Shutdown()

	if *host == "" {
		logger.Fatalln("-host needs to be set")
	}
	if *keyFile == "" {
		logger.Fatalln("-key needs to be set")
	}

	key, err := readAPIKey(*keyFile)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	controller := unifi.NewController(*host, key)

	info, err := controller.GetInfo()
	if err != nil {
		logger.Fatalln(err.Error())
	}
	fmt.Printf("Controller: %s\n", info.ApplicationVersion)

	sites, err := controller.ListSites()
	if err != nil {
		logger.Fatalln(err.Error())
	}
	for _, s := range sites {
		fmt.Printf("Site: %s\n", s.Name)

		if *devices {
			devices, err := controller.ListDevices(s.ID)
			if err != nil {
				logger.Fatalln(err.Error())
			}
			for _, d := range devices {
				if *deviceDetail {
					device, err := controller.GetDeviceDetail(s.ID, d.ID)
					if err != nil {
						logger.Fatalln(err.Error())
					}
					fmt.Printf("  * Device: %s, %s, %s, Updatable(%t)\n", d.Name, d.Model, d.State, device.FirmwareUpdatable)
				} else {
					fmt.Printf("  * Device: %s, %s, %s\n", d.Name, d.Model, d.State)
				}
			}

		}

		if *clients {
			clients, err := controller.ListClients(s.ID)
			if err != nil {
				logger.Fatalln(err.Error())
			}
			for _, c := range clients {
				fmt.Printf("  * Client: %s\n", c.Name)
			}
		}
	}
}
