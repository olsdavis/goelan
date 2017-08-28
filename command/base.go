package command

func RegisterBaseCommands() {
	RegisterCommand(BanCommand{})
	RegisterCommand(HelpCommand{})
	RegisterCommand(StopCommand{})
}
