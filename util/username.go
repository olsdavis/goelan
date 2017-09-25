package util

// IsValidUsername returns true if the given username is valid.
func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}
