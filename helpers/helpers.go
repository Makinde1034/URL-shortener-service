package helpers

import (
	"os"
	// "strings"
)

func RemoveDomainError(url string) bool{
	if url == os.Getenv("DOMAIN") {
		return false
	}

	return true
}

func EnforceHTTP(url string) string{
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}