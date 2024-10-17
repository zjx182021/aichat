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
		Temperature       float32 `mapstructure:"temperature"`
		TopP              float32 `mapstructure:"top_p"`
		FrequencyPenalty  float32 `mapstructure:"frequency_penalty"`
		PresencePenalty   float32 `mapstructure:"presence_penalty"`
		BotDesc           string  `mapstructure:"bot_desc"`
		MinResponseTokens int     `mapstructure:"min_response_tokens"`
		ContextTTL        int     `mapstructure:"context_ttl"`
		ContextLen        int     `mapstructure:"context_len"`
	} `mapstructure:"chat"`
	Mysql struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		DBname          string `mapstructure:"dbname"`
		Table           string `mapstructure:"table"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		MaxLifeTime     int    `mapstructure:"max_life_time"`
		MaxConnLifetime int    `mapstructure:"max_conn_lifetime"`
		IdleTimeout     int    `mapstructure:"idle_timeout"`
	} `mapstructure:"mysql"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		PoolSize int    `mapstructure:"pool_size"`
		MinIdle  int    `mapstructure:"min_idle"`
	} `mapstructure:"redis"`

	DependOn struct {
		Sensitive struct {
			Address     string `mapstructure:"address"`
			AccessToken string `mapstructure:"accessToken"`
		} `mapstructure:"sensitive"`
		Keywords struct {
			Address     string `mapstructure:"address"`
			AccessToken string `mapstructure:"accessToken"`
		} `mapstructure:"keywords"`
		Tokenizer struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"tokenizer"`
	} `mapstructure:"dependOn"`
	VectorDB struct {
		Url                string `mapstructure:"url"`
		Username           string `mapstructure:"username"`
		Pwd                string `mapstructure:"pwd"`
		Database           string `mapstructure:"database"`
		Timeout            int    `mapstructure:"timeout"`
		MaxIdleConnPerHost int    `mapstructure:"maxIdleConnPerHost"`
		ReadConsistency    string `mapstructure:"readConsistency"`
		IdleConnTimeout    int    `mapstructure:"idleConnTimeout"`
	} `mapstructure:"vectorDB"`
}

// func init() {
// 	InitConfig("./config.yaml", "yaml")
// }

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
