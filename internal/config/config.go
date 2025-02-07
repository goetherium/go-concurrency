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
	Wal        Wal        `yaml:"wal"`
}

type engine struct {
	Type string `yaml:"type"`
}

type connection struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize int           `yaml:"max_message_size"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type logger struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type Wal struct {
	FlushBatchSize    int           `yaml:"flush_batch_size"`
	FlushBatchTimeout time.Duration `yaml:"flush_batch_timeout"`
	FileMaxSize       uint64        `yaml:"file_max_size"`
	DataDir           string        `yaml:"data_dir"`
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
