package function

import (
	"context"
	"strconv"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"golang.org/x/oauth2"
)

func dsRecords(oauthToken string, domain string) ([]dnsimple.DelegationSignerRecord, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oauthToken})
	tc := oauth2.NewClient(context.Background(), ts)

	// new client
	client := dnsimple.NewClient(tc)

	// get the current authenticated account (if you don't know who you are)
	whoamiResponse, err := client.Identity.Whoami()
	if err != nil {
		return nil, err
	}

	// either assign the account ID or fetch it from the response
	// if you are authenticated with an account token
	accountID := strconv.FormatInt(whoamiResponse.Data.Account.ID, 10)

	// get the list of domains
	// domainsResponse, err := client.Domains.ListDomains(accountID, nil)
	//	domainsResponse, err := client.Domains.ListDomains(accountID, &dnsimple.DomainListOptions{NameLike: domain})
	dsRecords, err := client.Domains.ListDelegationSignerRecords(accountID, domain, &dnsimple.ListOptions{})

	if err != nil {
		return nil, err
	}

	return dsRecords.Data, nil
}
