package function

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

// DS Upload Sub-status codes.
// See: https://github.com/DK-Hostmaster/dsu-service-specification#http-sub-status-codes)
const (
	StatusUserIDNotSpecified                                   = 480
	StatusPasswordNotSpecified                                 = 481
	StatusMissingAParameter                                    = 482
	StatusDomainNameNotSpecified                               = 483
	StatusInvalidDomainName                                    = 484
	StatusInvalidUserID                                        = 485
	StatusInvalidDigestAndDigestTypeCombination                = 486
	StatusTheContentsOfAtLeastOneParameterIsSyntacticallyWrong = 487
	StatusAtLeastOneDSKeyHasAnInvalidAlgorithm                 = 488
	StatusInvalidSequenceOfSets                                = 489
	StatusUnknownParameterGiven                                = 495
	StatusUnknownUserID                                        = 496
	StatusUnknownDomainName                                    = 497
	StatusAuthenticationFailed                                 = 531
	StatusAuthorizationFailed                                  = 532
	StatusAuthenticatingUsingThisPasswordTypeIsNotSupported    = 533
)

var statusText = map[int]string{
	StatusUserIDNotSpecified:                                   "Userid not specified",
	StatusPasswordNotSpecified:                                 "Password not specified",
	StatusMissingAParameter:                                    "Missing a parameter",
	StatusDomainNameNotSpecified:                               "Domain name not specified",
	StatusInvalidDomainName:                                    "Invalid domain name",
	StatusInvalidUserID:                                        "Invalid userid",
	StatusInvalidDigestAndDigestTypeCombination:                "Invalid digest and digest_type combination",
	StatusTheContentsOfAtLeastOneParameterIsSyntacticallyWrong: "The contents of at least one parameter is syntactically wrong",
	StatusAtLeastOneDSKeyHasAnInvalidAlgorithm:                 "At least one DS key has an invalid algorithm",
	StatusInvalidSequenceOfSets:                                "Invalid sequence of sets",
	StatusUnknownParameterGiven:                                "Unknown parameter given",
	StatusUnknownUserID:                                        "Unknown userid",
	StatusUnknownDomainName:                                    "Unknown domain name",
	StatusAuthenticationFailed:                                 "Authentication failed",
	StatusAuthorizationFailed:                                  "Authorization failed",
	StatusAuthenticatingUsingThisPasswordTypeIsNotSupported:    "Authenticating using this password type is not supported",
}

func dsUpload(conf config, ds *dnsimple.DelegationSignerRecord) ([]byte, error) {
	form := url.Values{}
	form.Set("keytag1", ds.Keytag)
	form.Set("algorithm1", ds.Algorithm)
	form.Set("digest_type1", ds.DigestType)
	form.Set("digest1", ds.Digest)

	form.Set("domain", conf.Domain)
	form.Set("userid", conf.UserID)
	form.Set("password", conf.Password)

	resp, err := http.PostForm("https://dsu.dk-hostmaster.dk/1.0", form)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		return body, nil
	}

	subStatus := resp.Header.Get("X-DSU")

	if subStatus != "" {
		i, err := strconv.Atoi(subStatus)

		if err != nil {
			return nil, fmt.Errorf("DS Upload sub-status: %s", subStatus)
		}

		return nil, fmt.Errorf("DS Upload error: %s", statusText[i])
	}

	return nil, fmt.Errorf("DS Upload error: %s", resp.Status)
}
