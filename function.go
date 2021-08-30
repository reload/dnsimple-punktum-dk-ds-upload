package function

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"arnested.dk/go/dsupdate"
	"github.com/containrrr/shoutrrr"
	"github.com/dnsimple/dnsimple-go/dnsimple/webhook"
)

// Handle is the entrypoint for the Google Cloud Function.
func Handle(w http.ResponseWriter, r *http.Request) {
	// We log to Google Cloud Functions and don't need a timestamp
	// since it will be present in the log anyway. On the other
	// hand a reference to file and line number would be nice.
	log.SetFlags(log.Lshortfile)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if !isAuthorized(r.URL.Query().Get("token")) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	services, _ := notifyConfig()
	notify, err := shoutrrr.CreateSender(services.Services...)
	if err != nil {
		log.Printf("Error creating notification sender(s): %s", err.Error())
	}

	defer r.Body.Close()
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not parse webhook payload: %s", err.Error())
		notify.Send(fmt.Sprintf("Could not parse webhook payload: %s", err.Error()), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	event, err := webhook.ParseEvent(payload)

	log.Printf("Processing DNSimple event with request ID %q", event.RequestID)

	if err != nil {
		log.Printf("Could not parse webhook name: %s", err.Error())
		notify.Send(fmt.Sprintf("Could not parse webhook name: %s", err.Error()), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if event.Name != "dnssec.rotation_start" && event.Name != "dnssec.rotation_complete" {
		log.Printf("Not a rotation event: %s", event.Name)
		// It's OK if this is not a DNSSEC rotation event. We
		// send a 200 OK so DNSimple will not retry.
		http.Error(w, "Not a rotation event", http.StatusOK)

		return
	}

	dnssecEvent, ok := event.GetData().(*webhook.DNSSECEventData)

	if !ok {
		log.Printf("Could not parse webhook DNSSEC rotation event: %s", err.Error())
		notify.Send(fmt.Sprintf("Could not parse webhook DNSSEC rotation event: %s", err.Error()), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	config, err := envConfig(dnssecEvent.DelegationSignerRecord.DomainID)
	if err != nil {
		log.Printf("No DK Hostmaster / DNSimple config for %d: %s", dnssecEvent.DelegationSignerRecord.DomainID, err.Error())
		// It's OK if there is no configuration. It could be a
		// domain not handled by DK Hostmaster and/or DNSSEC.
		// We send a 200 OK so DNSimple will not retry.
		http.Error(w, "Missing DK Hostmaster / DNSimple credentials config", http.StatusOK)

		return
	}

	client := dsupdate.Client{
		Domain:   config.Domain,
		UserID:   config.UserID,
		Password: config.Password,
	}

	records, err := dsRecords(config.DnsimpleToken, config.Domain)
	if err != nil {
		log.Printf("Could not get DS records from DNSimple for %q: %s", config.Domain, err.Error())
		notify.Send(fmt.Sprintf("Could not get DS records from DNSimple for %q: %s", config.Domain, err.Error()), nil)
		http.Error(w, "Could not get DS records from DNSimple", http.StatusInternalServerError)

		return
	}

	// We'll set a 50 second timeout in the deletion using the
	// context package.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	resp, err := client.Update(ctx, records)
	if err != nil {
		log.Printf("Could not update DS records for %q: %s", config.Domain, err.Error())
		notify.Send(fmt.Sprintf("Could not update DS records for %q: %s", config.Domain, err.Error()), nil)
		http.Error(w, "Could not update DS records", http.StatusInternalServerError)

		return
	}

	log.Printf("Successful update of DS records for %q: %s", config.Domain, resp)

	errors := notify.Send(fmt.Sprintf("Successful update of DS records for %q: %s", config.Domain, resp), nil)
	if countErrors(errors) > 0 {
		log.Printf("Could not send Shoutrrr status: %v", err)
	}

	_, _ = w.Write(resp)
}

// countErrors but not nils.
func countErrors(slice []error) int {
	i := 0

	for _, elem := range slice {
		if elem != nil {
			i++
		}
	}

	return i
}
