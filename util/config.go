package util

import (
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"log"
)

type Config struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	UserServerAddress string `mapstructure:"USER_SERVER_ADDRESS"`
}

func LoadClientOption() option.ClientOption {
	return option.WithCredentialsFile("serviceAccount.json")
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = viper.Unmarshal(&config)
	return
}
