package main

import (
	"github.com/charlierm/whine/govee"
	"github.com/charlierm/whine/weather"
	log "github.com/sirupsen/logrus"
)
import "github.com/cloudkucooland/go-kasa"

func main() {
	getConditions()
	discovery, err := kasa.BroadcastDiscovery(10, 10)
	if err != nil {
		log.Warnf("error discovering switch devices: %v", err)
	}
	if len(discovery) > 1 || len(discovery) == 0 {
		log.Warnf("expected 1 switch device, got %d, will continue monitoring temperature", len(discovery))
	}

	var device *kasa.Device
	for ip, _ := range discovery {
		device, err = kasa.NewDevice(ip)
		if err != nil {
			log.Warnf("issue setting up device: %v", err)
			break
		}
	}

	// Get the current energy usage
	emeter, err := device.GetEmeter()
	if err != nil {
		log.Warnf("failed to get energy meter: %v", err)
	}
	log.Infof("current energy usage: %v", emeter.CurrentMA)

	println(discovery)
	println(weather.GetTemperature())
}

func getConditions() {
	goveeClient := govee.NewClient("")
	_, err := goveeClient.ListDevices()
	//devices[0].
	if err != nil {
		log.Warnf("error listing govee devices: %v", err)
	}
}
