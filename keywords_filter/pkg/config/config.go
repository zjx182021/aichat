package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		IP          string
		Port        int
		AccessToken string
	}
	Log struct {
		Level   string
		LogPath string
	}
}

var conf *Config

func InitConfig(filepath string, typ ...string) {
	v := viper.New()
	v.SetConfigFile(filepath)
	if len(typ) != 0 {
		v.SetConfigType(typ[0])
	}
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	conf = &Config{}
	err = v.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
}
func GetConfig() *Config {
	return conf
}
