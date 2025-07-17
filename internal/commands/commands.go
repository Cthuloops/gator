package commands

import (
	"fmt"

	"gator/internal/config"
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
