package common

import (
	"github.com/spf13/viper"
	"os"
)

func init() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("failed to read config, err: " + err.Error())
	}
}
