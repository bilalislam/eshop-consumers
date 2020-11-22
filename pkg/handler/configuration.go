package configuration

import (
	"github.com/spf13/viper"
)

type ConfigurationUseCase interface {
	GetConfig() (*viper.Viper, error)
}

type ConfigHandler struct {
	ConfigClient ConfigurationUseCase
}

func (c ConfigHandler) NewConfigHandler() (*viper.Viper, error) {
	return c.ConfigClient.GetConfig()
}
