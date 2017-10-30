package event

type EventHook interface {
	// OnEvent is a function called whenever the event for which
	// this EventHook has been registered is called.
	OnEvent(event *Event)
}

type EventManager struct {

}
