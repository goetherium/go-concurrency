package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type App struct {
	Engine     engine     `yaml:"engine"`
	Connection connection `yaml:"connection"`
	Logger     logger     `yaml:"logger"`
}

type engine struct {
	Type string `yaml:"type"`
}

type connection struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize int           `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type logger struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func Setup(path string) App {
	payload, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var app App
	if err = yaml.Unmarshal(payload, &app); err != nil {
		panic(err)
	}

	return app
}
