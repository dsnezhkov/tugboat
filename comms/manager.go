package comms

import (
	"tugboat/comms/http"
	"tugboat/payloads"
)

type CommsManager struct {
	HttpComm *http.HttpComm
	Payman *payloads.PayloadManager
}

type CommsManagerOptions struct {
	Payman *payloads.PayloadManager
}

var commsManager *CommsManager

func GetCommsManager() *CommsManager {
	return commsManager
}
func CreateCommsManager(options CommsManagerOptions) *CommsManager {

	if  commsManager == nil {

		httpCommsOptions := http.HttpCommOptions{
			Timeout: 10,
			Payman: options.Payman,
		}

		commsManager = &CommsManager{
			HttpComm: http.NewHttpComm(httpCommsOptions),
		}
	}
	return commsManager
}

func (comms *CommsManager) GetHTTPComm() *http.HttpComm{
	return comms.HttpComm
}


