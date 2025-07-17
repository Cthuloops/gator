package main

import (
	"log"
	"os"

	"gator/internal/commands"
	"gator/internal/config"
)

func main() {
	state := &config.State{
		Config: config.Read(),
	}
	cmds := commands.NewCommands()
	cmds.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Too few arguments")
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}
	err := cmds.Run(state, commands.NewCommand(cmdName, cmdArgs))
	if err != nil {
		log.Fatal(err)
	}
}
