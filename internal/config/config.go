package config

import "github.com/spf13/viper"
import "strings"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Cors     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	TTL      string `mapstructure:"ttl"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

func Load() (Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	// Override with environment variables if set
	v.SetEnvPrefix("")
	v.BindEnv("server.base_url", "BASE_URL")
	v.BindEnv("database.url", "DATABASE_URL")
	v.BindEnv("redis.url", "REDIS_URL")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("redis.ttl", "REDIS_TTL")

	cfg.Server.BaseURL = v.GetString("server.base_url")
	cfg.Database.URL = v.GetString("database.url")
	cfg.Redis.URL = v.GetString("redis.url")
	cfg.Redis.Password = v.GetString("redis.password")
	cfg.Redis.TTL = v.GetString("redis.ttl")

	return cfg, nil
}
