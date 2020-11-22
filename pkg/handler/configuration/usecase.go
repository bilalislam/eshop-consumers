package configuration

import (
	"os"

	"github.com/spf13/viper"
)

type ConfigurationClient struct {
	Filename string
	Path     string
}

func NewConfigClient(filename, path string) *ConfigurationClient {
	return &ConfigurationClient{
		Filename: filename,
		Path:     path,
	}
}

func (c *ConfigurationClient) GetConfig() (*viper.Viper, error) {
	v := viper.New()

	environment := os.Getenv("ENV_FILE")
	if environment != "" {
		c.Filename = c.Filename + "." + environment
	}

	v.SetConfigName(c.Filename)
	v.AddConfigPath(c.Path)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
