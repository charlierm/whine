package govee

func newDeviceStateRequest() *deviceStateRequest {
	req := &deviceStateRequest{
		RequestId: "uuid",
	}
	req.Payload.Sku = "sku"
	return req
}

type State struct {
	IsOnline    bool
	Temperature float64
	Humidity    float64
}

type deviceStateRequest struct {
	RequestId string `json:"requestId"`
	Payload   struct {
		Sku    string `json:"sku"`
		Device string `json:"device"`
	} `json:"payload"`
}

type devicesResponse struct {
	Data []struct {
		Sku        string `json:"sku"`
		Device     string `json:"device"`
		DeviceName string `json:"deviceName"`
	} `json:"data"`
}

type deviceStateResponse struct {
	Payload struct {
		Sku          string `json:"sku"`
		Device       string `json:"device"`
		Capabilities []struct {
			Type     string `json:"type"`
			Instance string `json:"instance"`
			State    struct {
				Value interface{} `json:"value"`
			} `json:"state"`
		} `json:"capabilities"`
	} `json:"payload"`
}
