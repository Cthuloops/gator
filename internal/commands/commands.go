package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"gator/internal/config"
	"gator/internal/database"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	command_list map[string]func(*config.State, command) error
}

func (c *commands) Run(s *config.State, cmd command) error {
	user_cmd, ok := c.command_list[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s does not exist", cmd.name)
	}
	return user_cmd(s, cmd)
}

func (c *commands) Register(name string, f func(*config.State, command) error) {
	c.command_list[name] = f
}

func NewCommand(name string, arguments []string) command {
	return command{
		name:      name,
		arguments: arguments,
	}
}

func NewCommands() commands {
	return commands{
		command_list: make(map[string]func(*config.State, command) error),
	}
}

///////////////////////////////////////////////////////////////////////////////

func HandlerLogin(s *config.State, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("Login expects a single argument, the username")
	}

	_, err := s.DB.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Fatalf("User %s does not exist", cmd.arguments[0])
		} else {
			return err
		}
	}

	err = s.Config.SetUser(&cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Failed to set user (%v): %w", cmd.arguments[0], err)
	}

	fmt.Printf("Set user to %s\n", s.Config.Current_user_name)
	return nil
}

func HandlerRegister(s *config.State, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("Register expects a single argument, the users name")
	}

	// Check if the user exists
	exists := true
	_, err := s.DB.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists = false
		} else {
			return err
		}
	}

	if exists {
		log.Fatalf("Trying to register a user (%s) that already exists",
			cmd.arguments[0])
	}

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	user, err := s.DB.CreateUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("Something went wrong %w", err)
	}

	// Set the current user to the newly added one.
	err = s.Config.SetUser(&cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s created", cmd.arguments[0])
	log.Printf(`
		ID: %v,
		CreatedAt: %v,
		UpdatedAt: %v,
		Name: %s
`, user.ID, user.CreatedAt, user.UpdatedAt, user.Name)

	return nil
}
