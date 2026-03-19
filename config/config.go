package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Tink     TinkConfig     `mapstructure:"tink"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpiryHours int    `mapstructure:"expiry_hours"`
	Issuer      string `mapstructure:"issuer"`
	JWKSUri     string `mapstructure:"jwks_uri"`
	Audience    string `mapstructure:"audience"`
}

type TinkConfig struct {
	Keyset string `mapstructure:"keyset"`
}

// Load reads config from file + env vars using Viper.
func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "logs/app.log")
	viper.SetDefault("jwt.expiry_hours", 24)

	// Explicit env bindings for nested keys (AutomaticEnv + Unmarshal is unreliable)
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "SERVER_MODE")
	_ = viper.BindEnv("database.dsn", "DATABASE_DSN")
	_ = viper.BindEnv("log.level", "LOG_LEVEL")
	_ = viper.BindEnv("log.filename", "LOG_FILENAME")
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")
	_ = viper.BindEnv("jwt.expiry_hours", "JWT_EXPIRY_HOURS")
	_ = viper.BindEnv("jwt.issuer", "JWT_ISSUER")
	_ = viper.BindEnv("jwt.jwks_uri", "JWT_JWKS_URI")
	_ = viper.BindEnv("jwt.audience", "JWT_AUDIENCE")
	_ = viper.BindEnv("tink.keyset", "TINK_KEYSET")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return &cfg
}
