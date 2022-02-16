package utils

import "github.com/spf13/viper"

// Save a value to the configuration file.
func SaveConfig(key string, value interface{}) {
	viper.Set(key, value)
	viper.ReadInConfig()
	viper.WriteConfig()
}
