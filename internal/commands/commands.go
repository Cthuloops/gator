package commands

import (
	"fmt"

	"gator/internal/config"
)

type command struct {
	name      string
	arguments []string
}

type Commands struct {
	command_list map[string]func(*config.State, command) error
}

func (c *Commands) run(s *config.State, cmd command) error {
	user_cmd, ok := c.command_list[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s does not exist", cmd.name)
	}
	return user_cmd(s, cmd)
}

func (c *Commands) register(name string, f func(*config.State, command) error) {
	c.command_list[name] = f
}
