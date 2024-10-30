package main

import (
	"fmt"
	"github.com/charlierm/whine/govee"
	"github.com/cloudkucooland/go-kasa"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const apiKeyEnvVar = "GOVEE_API_KEY"
const minimumTemperature = 15.0

func getDevice() (*govee.Device, error) {
	apiKey := os.Getenv(apiKeyEnvVar)
	goveeClient := govee.NewClient(apiKey)
	devices, err := goveeClient.ListDevices()
	if err != nil {
		return nil, fmt.Errorf("error listing govee devices: %v", err)
	}
	if len(devices) == 0 || len(devices) > 1 {
		return nil, fmt.Errorf("expected 1 device got %d: %v", len(devices), err)
	}
	return &devices[0], nil
}

func getHeaterSwitch() (*kasa.Device, error) {
	discovery, err := kasa.BroadcastDiscovery(10, 10)
	if err != nil {
		log.Warnf("error discovering switch devices: %v", err)
	}
	if len(discovery) > 1 || len(discovery) == 0 {
		log.Warnf("expected 1 switch device, got %d, will continue monitoring temperature", len(discovery))
	}

	var device *kasa.Device
	for ip := range discovery {
		device, err = kasa.NewDevice(ip)
		if err != nil {
			log.Warnf("issue setting up device: %v", err)
			break
		}
	}
	if device != nil {
		return device, nil
	}
	return nil, fmt.Errorf("no switch device found, cannot turn on heater")
}

func monitorTemperature(ticker *time.Ticker, done <-chan struct{}) {
	var device *govee.Device
	var err error
	var last time.Time

	for {
		select {
		case <-done:
			log.Infof("stopping temperature monitoring")
			return
		case <-ticker.C:
			if device == nil {
				device, err = getDevice()
				if err != nil {
					log.Warnf("error getting device: %v", err)
					continue
				}
			}
			state, err := device.GetState()
			if err != nil {
				log.Warnf("error getting device state: %v", err)
				continue
			}
			log.Infof("current temperature: %.2f°c", state.Temperature)
			if time.Since(last) >= time.Hour {
				// Turn on!
				log.Info("turning on heater")
				heater, err := getHeaterSwitch()
				if err != nil {
					log.Warnf("error getting heater switch: %v", err)
					continue
				}
				err = heater.SetRelayState(true)
				if err != nil {
					log.Warnf("error turning on heater: %v", err)
					continue
				}
				last = time.Now()
				continue
			}
			log.Infof("not ready to turn on %s since last switch on", time.Since(last))

		}
	}
}
