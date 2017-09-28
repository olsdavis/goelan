package util

import "regexp"

var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]{1,16}$")

// IsValidUsername returns true if the given username is valid.
func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
