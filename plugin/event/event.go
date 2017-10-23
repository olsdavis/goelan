package event

type Event interface {
	// IsCancellable returns true if the current event may be cancelled.
	IsCancellable() bool

	// SetCancelled makes the event cancelled or not. The value is taken
	// into account only if the event is cancellable - you can determine
	// it by reading its documentation or calling the IsCancellable function.
	SetCancelled(bool)

	// IsCancelled returns true if the current event will be cancelled.
	// (If the event is cancellable.)
	IsCancelled() bool
}
