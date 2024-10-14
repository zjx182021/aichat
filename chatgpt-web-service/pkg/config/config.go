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
	Chat struct {
		APIKey            string  `mapstructure:"api_key"`
		BaseURL           string  `mapstructure:"base_url"`
		Model             string  `mapstructure:"model"`
		MaxTokens         int     `mapstructure:"max_tokens"`
		Temperature       float64 `mapstructure:"temperature"`
		TopP              float64 `mapstructure:"top_p"`
		FrequencyPenalty  float64 `mapstructure:"frequency_penalty"`
		PresencePenalty   float64 `mapstructure:"presence_penalty"`
		BitDesc           string  `mapstructure:"bit_desc"`
		MinResponseTokens int     `mapstructure:"min_response_tokens"`
		ContextTTL        int     `mapstructure:"context_ttl"`
		ContextLen        int     `mapstructure:"context_len"`
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
