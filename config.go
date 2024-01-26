package main

import (
	telegram "github.com/pashandor789/broadcaster/bot"
	"github.com/pashandor789/broadcaster/http"
	"github.com/spf13/viper"
	"net/url"
)

func SetupConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetDatabaseURL() string {
	u := url.URL{
		Host:   viper.GetString("database.address") + ":" + viper.GetString("database.port"),
		Scheme: viper.GetString("database.scheme"),
		User:   url.UserPassword(viper.GetString("database.username"), viper.GetString("database.password")),
		Path:   viper.GetString("database.name"),
	}
	return u.String()
}

func GetBotConfig() telegram.BotConfig {
	return telegram.BotConfig{
		Token: viper.GetString("bot.token"),
	}
}

func GetServerConfig() http.ServerConfig {
	return http.ServerConfig{
		Port: viper.GetUint16("server.port"),
	}
}
