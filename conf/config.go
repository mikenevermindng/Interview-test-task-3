package conf

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type MonitorConfiguration struct {
	MaxMonitorConcurrency int `mapstructure:"maxMonitorConcurrency"`
	MonitorIntervalMs     int `mapstructure:"monitorIntervalMs"`
	MaxRedirect           int `mapstructure:"maxRedirect"`
	RequestTimeout        int `mapstructure:"requestTimeout"`
}

type ApiConfiguration struct {
	AdminSecret string `mapstructure:"adminSecret"`
}

type DbConfiguration struct {
	Sqlite string `mapstructure:"sqlite"`
}

type Configuration struct {
	Monitor  MonitorConfiguration
	Api      ApiConfiguration
	Database DbConfiguration
}

func Initialize() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to read current path")
		panic(err)
	}

	appEnv := os.Getenv("APP_ENV")

	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("%s/env", pwd))
	viper.SetConfigFile(fmt.Sprintf("%s/env/%s.yaml", pwd, appEnv))
	viper.SetConfigType("yaml")
	viper.SetConfigName(appEnv)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unable to read configuration file")
		panic(err)
	}
}

func NewConfiguration() *Configuration {
	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to bind env")
		panic(err)
	}
	return &config
}
