package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string   `yaml:"port"`
	Domain   string   `yaml:"domain"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

type Postgres struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	DBName            string `yaml:"db_name"`
	MaxOpenConnection int    `yaml:"max_open_connection"`
	MaxIdleConnection int    `yaml:"max_idle_connection"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func LoadConfig() Config {
	file, err := os.Open("cmd/http/config/local.yml")
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
