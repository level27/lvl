package utils

import "github.com/spf13/viper"

func SaveConfig(key string, value interface{}) {
	viper.Set(key, value)
	viper.ReadInConfig()
	viper.WriteConfig()
}
