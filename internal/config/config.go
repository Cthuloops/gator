package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path"
)

const configFileName string = ".gatorconfig.json"

type Config struct {
	DB_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) SetUser(user *string) error {
	if user == nil {
		c.Current_user_name = "harley"
	} else {
		c.Current_user_name = *user
	}
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func Read() Config {
	filePath, err := getConfigPath()
	if err != nil {
		log.Fatal("HOME environment variable not set.")
	}

	gatorConfig, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open %s: %v\n", filePath, err)
	}
	defer gatorConfig.Close()

	stats, err := gatorConfig.Stat()
	if err != nil {
		log.Fatal(err)

	}

	buf := make([]byte, stats.Size())
	n, err := gatorConfig.Read(buf)
	if err != nil {
		if err == io.EOF && n != len(buf) {
			log.Fatalf("Error reading whole file: %v\n", err)
		}
		if err != io.EOF {
			log.Fatalf("Error reading file: %v\n", err)
		}
	}

	data := Config{}
	if err = json.Unmarshal(buf, &data); err != nil {
		log.Fatalf("Error converting json: %v\n", err)

	}

	return data
}

func write(config Config) error {
	filePath, err := getConfigPath()
	if err != nil {
		log.Println(err)
		return err
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	gatorConfig, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer gatorConfig.Close()

	n, err := gatorConfig.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		log.Printf("Expected %d bytes written, wrote %d\n", len(data), n)
		return errors.New("Didn't write full config struct to file")
	}

	return nil
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := path.Join(home, configFileName)
	return filePath, nil
}
