package function

import (
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Domain   string `required:"true"`
	UserID   string `required:"true"`
	Password string `required:"true"`
	// We want the DNSimple token in an unprefixed environment
	// variable because it can serve multiple domains.
	DnsimpleToken string `required:"true" envconfig:"DNSIMPLE_TOKEN"`
}

func getConfig(domainId int64) (config, error) {
	domain := strconv.FormatInt(domainId, 10)
	config := config{}

	err := envconfig.Process(domain, &config)

	return config, err
}
