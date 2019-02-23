/*

Package function is a Google Cloud Function receiving webhook events
from DNSimple (https://dnsimple.com/webhooks).

It reacts to `dnssec.rotation_complete` events and passes the new DS
record on to DK Hostmaster via their DS Update protocol
(https://github.com/DK-Hostmaster/dsu-service-specification).

The cloud function needs to be configured through environment variables.

The `TOKEN` environment variable is the access token that should be
added as URL query parameter to the trigger URL (e.g.
`?token=abcdefeghijklmnopqrstuvxyz0123456789`).

For the domains in your DNSimple account that you would like this
cloud function to update in DK Hostmaster you need to add three
environment variables. They should all be prefix with the Domain ID
from DNSimple (e.g. 123456).

- `123456_DOMAIN`: the (apex) domain name in DK Hostmaster.

- `123456_USERID`: the DK Hostmaster handle you use to login to their
self service.

- `123456_PASSWORD`: the DK Hostmaster password you use to login to
their self service.

*/
package function
