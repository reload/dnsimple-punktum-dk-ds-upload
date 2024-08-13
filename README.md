# DNSimple Punktum DS upload

> [!IMPORTANT]
> Deprecated: Punktum.dk has closed the DS-update Service, see
> <https://punktum.dk/artikler/breaking-changes>

[![Go reference](https://pkg.go.dev/badge/github.com/reload/dnsimple-punktum-dk-ds-upload)](https://pkg.go.dev/github.com/reload/dnsimple-punktum-dk-ds-upload)

Package function is a Google Cloud Function receiving webhook events
from DNSimple (https://dnsimple.com/webhooks).

It reacts to `dnssec.rotation_start` and `dnssec.rotation_complete`
events and passes the new DS record on to Punktum.dk via their DS
Update protocol
(https://github.com/Punktum-dk/dsu-service-specification).

The cloud function needs to be configured through environment variables.

The `TOKEN` environment variable is the access token that should be
added as URL query parameter to the trigger URL (e.g.
`?token=abcdefeghijklmnopqrstuvxyz0123456789`).

The `DNSIMPLE_TOKEN` environment variable is a DNSimple API token that
is used to retrieve DS records from DNsimple.

For the domains in your DNSimple account that you would like this
cloud function to update in Punktum.dk you need to add three
environment variables. They should all be prefix with the Domain ID
from DNSimple (e.g. 123456).

`123456_DOMAIN`: the (apex) domain name in Punktum.dk.

`123456_USERID`: the Punktum.dk handle you use to login to their
self service.

`123456_PASSWORD`: the Punktum.dk password you use to login to
their self service.




