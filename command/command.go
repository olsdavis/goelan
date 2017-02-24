package command

import (
	"fmt"
)

var (
	ConsoleSender = consoleSender{}
)

type consoleSender struct {
}

func (sender consoleSender) SendMessage(message string) {
	fmt.Println(message)
}

func (sender consoleSender) IsPlayer() bool {
	return false
}

func (sender consoleSender) HasPermission(perm string) bool {
	return true
}

type Command interface {
	// Must return the different labels with which you can call the command
	// (at least one).
	Labels() []string

	// The minimal number of arguments - 0 if not needed.
	MinArgs() int

	// Returns the permission required to execute the command.
	RequiredPermission() string

	// Returns command's help.
	Help() string

	// Returns command's description - used for the help command
	Description() string

	// Called when the command has to be executed.
	// Label is the label used by the sender,
	// arguments are the strings following the command (command <args>),
	// sender is the sender.
	Execute(label string, arguments []string, sender CommandSender)
}

type CommandSender interface {
	// Sends a message to the sender.
	SendMessage(string)

	// Returns true if the command sender is a player.
	IsPlayer() bool

	// Returns true if the command sender has the given permission.
	HasPermission(string) bool
}
