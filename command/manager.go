package command

import (
	"bufio"
	"fmt"
	set "github.com/deckarep/golang-set"
	"os"
	"strings"
)

var (
	commands       = make(map[string]Command)
	ConsoleChannel = make(chan string)
)

// Returns all the registered commands.
func Commands() set.Set {
	ret := set.NewThreadUnsafeSet()
	for _, command := range commands {
		ret.Add(command)
	}
	return ret
}

// Registers the given command.
func RegisterCommand(command Command) {
	for _, label := range command.Labels() {
		commands[label] = command
	}
}

// Called when a command is executed by the console or a player.
func ExecuteCommand(line string, sender CommandSender) {
	if len(line) == 0 {
		return
	}

	parts := strings.Split(line, " ")
	label := parts[0]
	if handler, ok := commands[label]; ok {
		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}
		handler.Execute(label, args, sender)
	} else {
		sender.SendMessage(fmt.Sprintf("Could not find the command \"%v\". Type help to see the full list.", label))
	}
}

func ReadInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		ConsoleChannel <- line
	}
}
