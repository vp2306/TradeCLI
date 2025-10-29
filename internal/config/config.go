package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Alpaca AlpacaConfig `mapstructure:"alpaca"`
	Trading TradingConfig `mapstructure:"trading"`
}

type AlpacaConfig struct {
	APIKey    string `mapstructure:"api_key"`
	APISecret string `mapstructure:"api_secret"`
	BaseURL   string `mapstructure:"base_url"`
	DataURL   string `mapstructure:"data_url"`
}

type TradingConfig struct {
	MaxPositions     int     `mapstructure:"max_positions"`
	DefaultOrderSize int     `mapstructure:"default_order_size"`
	RiskPerTrade     float64 `mapstructure:"risk_per_trade"`
	MaxOrderValue    float64 `mapstructure:"max_order_value"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}