package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Http struct {
		Host        string
		Port        int
		AccessToken string `mapstructure:"access_token"`
		Mode        string
	}
	Chat struct {
		APIKeys []string `mapstructure:"api_keys"`
		BaseURL string   `mapstructure:"base_url"`
	}
	Log struct {
		Level   string
		LogPath string `mapstructure:"log_path"`
	}
}

var conf *Config

func InitConfig(filePath string, typ ...string) {
	v := viper.New()
	v.SetConfigFile(filePath)
	if len(typ) > 0 {
		v.SetConfigType(typ[0])
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	conf = &Config{}
	err = v.Unmarshal(conf)
	if err != nil {
		log.Fatal(err)
	}

}

func GetConfig() *Config {
	return conf
}
