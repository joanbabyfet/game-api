package bootstrap

import (
	"fmt"
	"sync"
	"time"

	"game-api/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	mysqlOnce sync.Once
)

// InitMySQL 初始化 MySQL（单例）
func InitMySQL() error {

	var err error

	mysqlOnce.Do(func() {

		cfg := config.Cfg.MySQL

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		sqlDB, e := DB.DB()
		if e != nil {
			err = e
			return
		}

		// 最大空闲连接
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)

		// 最大连接数
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)

		// 连接最大存活时间
		sqlDB.SetConnMaxLifetime(
			time.Duration(cfg.MaxLifetime) * time.Second,
		)

		err = sqlDB.Ping()
	})

	return err
}

// CloseMySQL 关闭 MySQL
func CloseMySQL() error {

	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}