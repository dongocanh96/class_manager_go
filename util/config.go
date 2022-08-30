package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	SignUpKeyForTeacher string        `mapstructure:"SIGN_UP_KEY_FOR_TEACHER"`
	PrivateKeyLocation  string        `mapstructure:"PRIVATE_KEY_LOCATION"`
	PublicKeyLocation   string        `mapstructure:"PUBLIC_KEY_LOCATION"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	Asset               string        `mapstructure:"ASSET"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  //path is the location of config file
	viper.SetConfigName("app") // app is the name of config file
	viper.SetConfigType("env") // env is the type of config file

	viper.AutomaticEnv() // automatically override values read from config file if evironment variable exist

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
