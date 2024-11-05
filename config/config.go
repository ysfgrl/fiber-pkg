package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	env := os.Getenv("ENV")
	dir, _ := os.Getwd()
	fmt.Printf("Current Path %s\n", dir)
	fmt.Printf("Loading config from environment variable %s\n", env)
	viper.SetConfigFile(env)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func Read(conf interface{}) {
	err := viper.Unmarshal(conf)
	if err != nil {
		panic(err)
	}
}
