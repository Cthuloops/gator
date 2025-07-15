package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser(nil)
	cfg = config.Read()
	fmt.Println(cfg)
}
