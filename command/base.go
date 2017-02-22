package command

import (
	"../permission"
	"github.com/gosuri/uitable"
	"strings"
)

type helpCommand struct{}

func (cmd helpCommand) Labels() []string {
	return []string{"help", "?"}
}

func (cmd helpCommand) MinArgs() int {
	return 0
}

func (cmd helpCommand) RequiredPermission() string {
	return permission.BasePermission
}

func (cmd helpCommand) Help() string {
	return "help (command_name)"
}

func (cmd helpCommand) Description() string {
	return "Shows available commands with their descriptions, or, if indicated, a specific's command description."
}

func (cmd helpCommand) Execute(label string, args []string, sender CommandSender, manager *CommandManager) {
	if sender.IsPlayer() {
		sender.SendMessage("Unimplemented. Lol")
	} else {
		// Console
		table := uitable.New()
		table.MaxColWidth = 110
		table.AddRow("COMMAND LABEL(S)", "DESCRIPTION")
		for _, command := range manager.commands {
			table.AddRow(strings.Join(command.Labels(), ", "), command.Description())
		}
		sender.SendMessage(table.String())
	}
}

func RegisterBaseCommands(manager *CommandManager) {
	manager.RegisterCommand(helpCommand{})
}
