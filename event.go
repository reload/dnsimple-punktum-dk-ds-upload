package function

import (
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
)

func dnsimpleEventName(payload []byte) (string, error) {
	event, err := webhook.ParseEvent(payload)

	if err != nil {
		return "", err
	}

	return event.Name, nil
}

func dnsimpleEvent(payload []byte) (*webhook.DNSSECEventData, error) {
	event, err := webhook.ParseEvent(payload)

	if err != nil {
		return nil, err
	}

	dnssecEvent, ok := event.GetData().(*webhook.DNSSECEventData)

	if !ok {
		return nil, fmt.Errorf("Could not parse event as a DNSSEC event.")
	}

	return dnssecEvent, nil
}
