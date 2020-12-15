package ultis

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port       int32
	DbUsername string
	DbPassword string
	DbUrl      string
	DbName     string
}

func LoadConfigFile() (*Config, error) {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config File (config.json.json) not found")
			panic(err)
		} else {
			fmt.Println("Config File error:")
			panic(err)
		}
	}

	c := &Config{}
	c.Port = viper.GetInt32("Port")
	c.DbUsername = viper.GetString("DbUsername")
	c.DbPassword = viper.GetString("DbPassword")
	c.DbUrl = viper.GetString("DbUrl")
	c.DbName = viper.GetString("DbName")

	return c, nil
}
