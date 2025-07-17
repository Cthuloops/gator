package commands

import (
	"errors"
	"fmt"

	"gator/internal/config"
)

func HandlerLogin(s *config.State, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("Login expects a single argument, the username")
	}

	err := s.Config.SetUser(&cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Failed to set user (%v): %w", cmd.arguments[0], err)
	}

	fmt.Printf("Set user to %s\n", s.Config.Current_user_name)
	return nil
}
