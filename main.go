package main

import (
	"./command"
	"./log"
	"./server"
	"os"
)

func main() {
	log.Init()

	defer log.Info("Goodbye!")
	log.Info("Starting up server...")
	// start up process
	srv := server.CreateServerFromProperties()
	go srv.CommandManager.ReadInput()
	go srv.Start()

	for srv.IsRunning() {
		select {
		case <-srv.ExitChan:
			break
		case line := <-srv.CommandManager.ConsoleChannel:
			srv.CommandManager.CommandExecute(line, command.ConsoleSender)
		}
	}

	os.Exit(0)
}
