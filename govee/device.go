package govee

import "net/http"

type Device struct {
	Sku        string
	DeviceId   string
	Name       string
	apiKey     string
	httpClient *http.Client
}

func (d Device) GetState() {

}
