package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommandManager struct {
	commands       map[string]Command
	ConsoleChannel chan string
}

func NewManager() *CommandManager {
	cm := &CommandManager{
		make(map[string]Command),
		make(chan string),
	}
	RegisterBaseCommands(cm)
	return cm
}

// Registers the given command.
func (manager *CommandManager) RegisterCommand(command Command) {
	for _, label := range command.Labels() {
		manager.commands[label] = command
	}
}

// Called when a command is executed by the console or a player.
func (manager *CommandManager) CommandExecute(line string, sender CommandSender) {
	if len(line) == 0 {
		return
	}

	line = line[:len(line)-2] // removes the new line char
	parts := strings.Split(line, " ")
	label := parts[0]
	if handler, ok := manager.commands[label]; ok {
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}
		handler.Execute(label, args, sender, manager)
	} else {
		sender.SendMessage(fmt.Sprintf("Could not find the command \"%v\". Type help to see the full list.", label))
	}
}

func (manager *CommandManager) ReadInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		manager.ConsoleChannel <- line
	}
}
