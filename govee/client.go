package govee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	goveeHost = "https://openapi.api.govee.com"
)

type Govee struct {
	client *http.Client
	apiKey string
}

// NewClient creates a new Govee client with the given API key.
func NewClient(apiKey string) Govee {
	govee := Govee{}
	govee.client = &http.Client{}
	govee.apiKey = apiKey
	return govee
}

func (g Govee) ListDevices() ([]Device, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/router/api/v1/user/devices", goveeHost), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListDevices request: %w", err)
	}
	setHeaders(req, g.apiKey)
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user/devices endpoint: %w", err)
	}
	defer resp.Body.Close()

	r := devicesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return g.mapResponseToDevices(r), nil
}

func (g Govee) GetDeviceState(sku string, device string) (*State, error) {
	stateReq := newDeviceStateRequest()
	stateReq.Payload.Device = device
	stateReq.Payload.Sku = sku

	marshalled, err := json.Marshal(stateReq)
	if err != nil {
		return nil, fmt.Errorf("impossible to marshall request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/router/api/v1/device/state", goveeHost), bytes.NewReader(marshalled))
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDeviceState request: %w", err)
	}
	setHeaders(req, g.apiKey)
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call device/state endpoint: %w", err)
	}
	defer resp.Body.Close()

	r := deviceStateResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return mapResponseToState(r), nil
}
func convertFahrenheitToCelsius(fahrenheit float64) float64 {
	return (fahrenheit - 32) * 5 / 9

}

func (g Govee) mapResponseToDevices(response devicesResponse) []Device {
	devices := make([]Device, len(response.Data))
	for i, d := range response.Data {
		devices[i] = Device{
			Sku:      d.Sku,
			DeviceId: d.Device,
			Name:     d.DeviceName,
			client:   g,
		}
	}
	return devices
}

func mapResponseToState(response deviceStateResponse) *State {
	state := State{}
	for _, c := range response.Payload.Capabilities {
		switch c.Instance {
		case "sensorTemperature":
			state.Temperature = convertFahrenheitToCelsius(c.State.Value.(float64)) // convert to celsius
		case "sensorHumidity":
			state.Humidity = c.State.Value.(map[string]interface{})["currentHumidity"].(float64) //
		case "IsOnline":
			state.IsOnline = c.State.Value.(bool)
		}
	}
	return &state
}

func setHeaders(r *http.Request, apiKey string) *http.Request {
	r.Header.Set("Govee-API-Key", apiKey)
	r.Header.Set("Content-Type", "application/json")
	return r
}
