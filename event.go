package function

import (
	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
)

func dnsimpleEventName(payload []byte) (string, error) {
	name, err := webhook.ParseName(payload)

	if err != nil {
		return "", err
	}

	return name, nil
}

func dnsimpleEvent(payload []byte) (*webhook.DNSSECEvent, error) {
	event := &webhook.DNSSECEvent{}
	err := webhook.ParseDNSSECEvent(event, payload)

	if err != nil {
		return nil, err
	}

	return event, nil
}
