package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	ServerAddr string `mapstructure:"SERVER_ADDRESS"`
}

func MustLoadConfig(path string) (config Config) {
	var err error
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not read config file %v\n", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("could not read config file %v\n", err)
	}
	return
}
