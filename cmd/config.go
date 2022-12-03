package cmd

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type config struct {
	Access struct {
		ClientID     string `yaml:"clientID"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"access"`
	Server string `yaml:"server"`
}

func initConfig() {
	if err := initConfigE(); err != nil {
		log.Fatal(err)
	}
}

func initConfigE() error {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	defaultPath := filepath.Join(confDir, "shrt", "config.yaml")

	if confPath == "" {
		confPath = defaultPath
	}

	f, err := os.Open(confPath)
	if err != nil {
		return err
	}

	defer f.Close()

	var c config

	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return err
	}

	serverAddr = c.Server
	clientID = c.Access.ClientID
	clientSecret = c.Access.ClientSecret

	return nil
}
