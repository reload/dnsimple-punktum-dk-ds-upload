package function

import (
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Domain   string `required:"true"`
	UserID   string `required:"true"`
	Password string `required:"true"`
}

func getConfig(domainId int64) (config, error) {
	domain := strconv.FormatInt(domainId, 10)
	config := config{}

	err := envconfig.Process(domain, &config)

	return config, err
}
