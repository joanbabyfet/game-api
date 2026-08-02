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
	Cron     CronConfig     `mapstructure:"cron"`
	Skynet   SkynetConfig   `mapstructure:"skynet"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxLifetime int    `mapstructure:"max_lifetime"`
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
	GameURL string `mapstructure:"game_url"`
	FileURL string `mapstructure:"file_url"`
}

type CronConfig struct {
	Report           string `mapstructure:"report"`
	Jackpot          string `mapstructure:"jackpot"`
	OrderRecover     string `mapstructure:"order_recover"` //每 10 秒扫描一次
	RecoverBatchSize int    `mapstructure:"recover_batch_size"` //每批最多 100 笔
	RecoverLockTTL   int    `mapstructure:"recover_lock_ttl"` //抢占锁 30 秒
	RecoverMaxRetry  int    `mapstructure:"recover_max_retry"` //最多自动重试 3 次
	RecoverMinAge    int    `mapstructure:"recover_min_age"` //订单至少停留 15 秒后才补偿
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
