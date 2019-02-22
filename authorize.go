package function

import "os"

func isAuthorized(tokenParam string) bool {
	token, ok := os.LookupEnv("TOKEN")

	if !ok {
		return false
	}

	if tokenParam == token {
		return true
	}

	return false
}
