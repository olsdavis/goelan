package util

import (
	"fmt"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]{1,16}$")
)

// Returns the uuid with the hyphens.
func ToHypenUUID(uuid string) string {
	// 8 - 4 - 4 - 4 - 12
	return fmt.Sprintf("%v-%v-%v-%v-%v", uuid[:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:])
}

// Returns true if the given username is valid
func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
