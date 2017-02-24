package main

import (
	"./command"
	"./log"
	"./server"
	"os"
)

func main() {
	log.Init()

	log.Info("Starting up server...")
	// start up process
	srv := server.CreateServerFromProperties()
	command.RegisterBaseCommands()
	go command.ReadInput()
	go srv.Start()

	for srv.IsRunning() {
		select {
		case <-srv.ExitChan:
			break
		case line := <-command.ConsoleChannel:
			command.CommandExecute(line, command.ConsoleSender)
		}
	}

	log.Info("Goodbye!")

	os.Exit(0)
}
