package utils

import "regexp"

const addressLength = 40

var addressRegex = regexp.MustCompile("^[a-f0-9]+$")

// ValidateAddress returns bool identifying is address is valid or not
func ValidateAddress(address string) bool {
	return addressRegex.MatchString(address) && len(address) == addressLength
}
