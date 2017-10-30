package plugin

// Plugin struct represents an extension for goelan which can handle events,
// create new commands, etc. In short: customize the server thanks to an API
// for the used language.
type Plugin struct {
	Language Language
}