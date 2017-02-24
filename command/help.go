package command

import (
	"../permission"
	"github.com/gosuri/uitable"
	"strings"
)

type HelpCommand struct{}

func (cmd HelpCommand) Labels() []string {
	return []string{"help", "?"}
}

func (cmd HelpCommand) MinArgs() int {
	return 0
}

func (cmd HelpCommand) RequiredPermission() string {
	return permission.BasePermission
}

func (cmd HelpCommand) Help() string {
	return "help (command_name)"
}

func (cmd HelpCommand) Description() string {
	return "Shows available commands with their descriptions, or, if indicated, a specific's command description."
}

func (cmd HelpCommand) Execute(label string, args []string, sender CommandSender) {
	if sender.IsPlayer() {
		sender.SendMessage("Unimplemented. Lol")
	} else {
		// Console
		table := uitable.New()
		table.MaxColWidth = 110
		table.AddRow("COMMAND LABEL(S)", "DESCRIPTION")

		for cmd := range Commands().Iter() {
			command := cmd.(Command)
			table.AddRow(strings.Join(command.Labels(), ", "), command.Description())
		}
		sender.SendMessage(table.String())
	}
}
