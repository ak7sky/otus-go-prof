package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

var (
	loggerLevel    = "Config.Logger.Level"
	storageType    = "Config.Storage.Type"
	serverHTTPPort = "Config.Server.HTTPPort"

	constraints = map[string]string{
		loggerLevel:    " must be one of 'DEBUG', 'INFO', 'ERROR'",
		storageType:    " must be one of 'memory', 'rdb'",
		serverHTTPPort: " must be in interval [1; 65535]",
	}
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level         string `yaml:"level" validate:"oneof=DEBUG INFO ERROR"`
	IsJSONEnabled bool   `yaml:"isJsonEnabled"`
}

type StorageConf struct {
	Type string `yaml:"type" validate:"oneof=memory rdb"`
	DSN  string `env:"DSN"`
}

type ServerConf struct {
	Host     string `yaml:"host"`
	HTTPPort int    `yaml:"httpPort" validate:"min=1,max=65535"`
}

func NewConfig(configFile string) (Config, error) {
	config := Config{}

	if err := cleanenv.ReadConfig(configFile, &config); err != nil {
		return Config{}, fmt.Errorf("failed to init config: %w", err)
	}

	if err := validate(config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func validate(config Config) error {
	validate := validator.New()

	if err := validate.Struct(config); err != nil {
		var msg strings.Builder
		msg.WriteString("invalid config:\n")

		for field, constraint := range constraints {
			if strings.Contains(err.Error(), field) {
				msg.WriteString(fmt.Sprintf("%s %s\n", field, constraint))
			}
		}
		return errors.New(msg.String())
	}

	return nil
}
