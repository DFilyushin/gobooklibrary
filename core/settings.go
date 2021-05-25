package core

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server 		string `mapstructure:"SERVER"`
	Port 		string `mapstructure:"PORT"`
	User 		string `mapstructure:"USER"`
	Password 	string `mapstructure:"PASSWORD"`
	Database	string `mapstructure:"DATABASE"`
}


func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}