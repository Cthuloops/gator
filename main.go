package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"gator/internal/commands"
	"gator/internal/config"
	"gator/internal/database"
)

func main() {
	state := &config.State{
		Config: config.Read(),
	}
	cmds := commands.NewCommands()
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)

	// Open connection to postgres.
	db, err := sql.Open("postgres", state.Config.DB_url)
	if err != nil {
		log.Fatal(err)
	}
	// Create queries object.
	dbQueries := database.New(db)
	// Store queries object in the state.
	state.DB = dbQueries

	if len(os.Args) < 2 {
		log.Fatal("Too few arguments")
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}
	err = cmds.Run(state, commands.NewCommand(cmdName, cmdArgs))
	if err != nil {
		log.Fatal(err)
	}
}
