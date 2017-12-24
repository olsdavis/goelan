package util

// Matcher is a function which tells whether a given object matches
// criteria defined by the Matcher.
type Matcher func(test interface{}) bool

// Completable is an interface which represents an action that can
// be completed, and thus "useless" once it has been (completed).
type Completable interface {
	// Complete completes the action, if it has not been done yet.
	Complete()
}
