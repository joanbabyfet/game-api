package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Cfg Config

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	Provider ProviderConfig `mapstructure:"provider"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Skynet 	 SkynetConfig 	`mapstructure:"skynet"`
}

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	MaxIdleConn int `mapstructure:"max_idle_conn"`
	MaxOpenConn int `mapstructure:"max_open_conn"`
	MaxLifetime  int `mapstructure:"max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type ProviderConfig struct {
	Port int `mapstructure:"port"`
}

type AdminConfig struct {
	Port int `mapstructure:"port"`
}

type WorkerConfig struct {
	Cron CronConfig `mapstructure:"cron"`
}

type CronConfig struct {
	Report  string `mapstructure:"report"`
	Jackpot string `mapstructure:"jackpot"`
}

type SkynetConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

// Load 支持加载多个 yaml，后面的覆盖前面的
func Load(files ...string) error {

	v := viper.New()

	for _, file := range files {

		v.SetConfigFile(file)

		if err := v.MergeInConfig(); err != nil {
			return fmt.Errorf("load %s failed: %w", file, err)
		}
	}

	if err := v.Unmarshal(&Cfg); err != nil {
		return err
	}

	return nil
}