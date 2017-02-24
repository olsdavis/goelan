package command

func RegisterBaseCommands() {
	RegisterCommand(HelpCommand{})
	RegisterCommand(StopCommand{})
}
