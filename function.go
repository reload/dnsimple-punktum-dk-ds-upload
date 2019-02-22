package function

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
)

// F is the entrypoint.
func Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if !isAuthorized(r.URL.Query().Get("token")) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	defer r.Body.Close()
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Could not read HTTP Body: %s", err.Error())
		http.Error(w, "Could not read HTTP Body", http.StatusBadRequest)

		return
	}

	event := &webhook.DNSSECEvent{}
	err = webhook.ParseDNSSECEvent(event, payload)

	if err != nil {
		log.Printf("Could not parse event as a DNSSEC event: %s", err.Error())
		http.Error(w, "Could not parse event as a DNSSEC event", http.StatusBadRequest)

		return
	}

	if event.Name != "dnssec.rotation_complete" {
		log.Printf("Not a `dnssec.rotation_complete` event: %s", event.Name)
		http.Error(w, "Not a `dnssec.rotation_complete` event", http.StatusUnprocessableEntity)

		return
	}

	// We get the config for DK Hostmaster from the
	// environment. They are prefixed with DNSimple domain ID.
	config, err := getConfig(event.DelegationSignerRecord.DomainID)

	if err != nil {
		log.Printf("Config problem: %s", err.Error())
		http.Error(w, "Missing DK Hostmaster credentials", http.StatusNotImplemented)

		return
	}

	body, err := dsUpload(config, event.DelegationSignerRecord)

	if err != nil {
		log.Printf("Upload problem: %s", err.Error())
		http.Error(w, "Internal server error uploading to DK Hostmaster", http.StatusInternalServerError)

		return
	}

	log.Printf("Upload succeeded: %s", body)
	_, _ = w.Write(body)
}
