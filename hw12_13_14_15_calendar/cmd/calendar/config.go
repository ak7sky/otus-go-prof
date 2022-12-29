package main

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
		loggerLevel:    " must be one of 'TRACE', 'DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL'",
		storageType:    " must be one of 'memory', 'rdb'",
		serverHTTPPort: " must be in interval [1; 65535]",
	}
)

// Config
// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string `yaml:"level" validate:"oneof=TRACE DEBUG INFO WARN ERROR FATAL"`
}

type StorageConf struct {
	Type string `yaml:"type" validate:"oneof=memory rdb"`
}

type ServerConf struct {
	Host     string `yaml:"host"`
	HTTPPort int    `yaml:"httpPort" validate:"min=1,max=65535"`
}

func NewConfig() (Config, error) {
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
