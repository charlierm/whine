package govee

type Device struct {
	Sku      string
	DeviceId string
	Name     string
	client   Govee
}

func (d Device) GetState() (*State, error) {
	state, err := d.client.GetDeviceState(d.Sku, d.DeviceId)
	if err != nil {
		return nil, err
	}
	return state, nil
}
