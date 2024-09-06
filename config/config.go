package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func InitialiseConfig(configPath string) {

	// Configurations not shared with the database
	viper.SetDefault("base_url", "http://localhost:8001")
	viper.SetDefault("port", "8001")
	viper.SetDefault("secret_key", "VmtkMFUxUnJNVlpOVnpWcFpXcENURU5u")
	viper.SetDefault("paseto_key", "Vmtab2QxSnJNSGRQVlZaWFZsaE9URU5u")
	viper.SetDefault("paseto_duration", "720h")
	viper.SetDefault("resolution", []int{768, 432})
	viper.Set("db_time_format", "2006-01-02 15:04:05")
	
	viper.SetDefault("allowed_origins", []string{
		"http://localhost:*",
		"https://localhost:*",
		"http://127.0.0.1:*",
		"https://127.0.0.1:*",
	})

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configPath)

	viper.SafeWriteConfig()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Error reading configuration file: %s", err))
	}

	viper.WriteConfig()
}
