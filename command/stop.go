package command

import (
	"github.com/olsdavis/goelan/permission"
	"github.com/olsdavis/goelan/server"
)

type StopCommand struct{}

func (cmd StopCommand) Labels() []string {
	return []string{"stop"}
}

func (cmd StopCommand) MinArgs() int {
	return 0
}

func (cmd StopCommand) RequiredPermission() string {
	return permission.StopServer
}

func (cmd StopCommand) Help() string {
	return "stop"
}

func (cmd StopCommand) Description() string {
	return "Stops the server."
}

func (cmd StopCommand) Execute(label string, args []string, sender CommandSender) {
	server.Get().Stop()
}
