package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	// Config -.
	Config struct {
		Env      string `mapstructure:"env"`
		App      `mapstructure:"app"`
		HTTP     `mapstructure:"http"`
		Log      `mapstructure:"logger"`
		Postgres `mapstructure:"postgres"`
		Redis    `mapstructure:"redis"`
		JWT      `mapstructure:"jwt"`
		Kafka    `mapstructure:"kafka"`
	}

	// App -.
	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	}

	// HTTP -.
	HTTP struct {
		Port string `mapstructure:"port"`
	}

	// Log -.
	Log struct {
		Level string `mapstructure:"log_level"`
	}

	// Posgres -.
	Postgres struct {
		Host     string `mapstructure:"host"`
		DbName   string `mapstructure:"db_name"`
		User     string `mapstructure:"user"`
		SSLMode  string `mapstructure:"ssl_mode"`
		Password string `mapstructure:"password"`
		Port     int    `mapstructure:"port"`
		Schema   string `mapstructure:"schema"`
	}

	// Redis -.
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}

	// Kafka -.
	JWT struct {
		Secret string `mapstructure:"secret"`
	}

	// Kafka -.
	Kafka struct {
		Brokers    string `mapstructure:"broker"`
		GroupID    string `mapstructure:"group_id"`
		OrderGroup `mapstructure:"order_group"`
	}

	//Order group-.
	OrderGroup struct {
		ConfirmTopic string `mapstructure:"confirm_topic"`
		GroupID      string `mapstructure:"group_id"`
	}
)

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s search_path=%s", p.Host, p.User, p.Password, p.DbName, p.Port, p.SSLMode, p.Schema)
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	viper.SetConfigName("config")

	viper.SetConfigType("yaml")

	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
