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

	err := cmds.Run(state, commands.NewCommand(os.Args[0], os.Args[1:]))
	if err != nil {
		log.Fatal(err)
	}
}
