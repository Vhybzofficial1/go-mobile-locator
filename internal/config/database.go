package config

import (
	"fmt"
	"log"
	"mobile-locator/internal/embedfiles"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := GetConfig()
	dbPath := GetSQLitePath(cfg)
	// 如果数据库不存在，则从内置资源复制一份
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		data, err := embedfiles.CarrierDB.ReadFile("assets/carrier.db")
		if err != nil {
			log.Fatalf("无法读取内置 carrier.db: %v", err)
		}
		err = os.WriteFile(dbPath, data, 0644)
		if err != nil {
			log.Fatalf("无法写入SQLite数据库: %v", err)
		}
	}
	// GORM 日志配置
	var gormLogLevel logger.LogLevel
	switch cfg.Database.Gorm.LogLevel {
	case 1:
		gormLogLevel = logger.Silent
	case 2:
		gormLogLevel = logger.Error
	case 3:
		gormLogLevel = logger.Warn
	case 4:
		gormLogLevel = logger.Info
	default:
		gormLogLevel = logger.Silent
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	// 获取通用数据库对象 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}
	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.Gorm.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.Gorm.MaxOpenConns)
	if cfg.Database.Gorm.ParsedConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.Database.Gorm.ParsedConnMaxLifetime)
	}
	DB = db
	log.Printf("数据库连接成功: %s", dbPath)
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("数据库未初始化，请先调用 InitDB()")
	}
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移
func AutoMigrate(model ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return DB.AutoMigrate(model...)
}

// IsConnected 检查数据库是否已连接
func IsConnected() bool {
	if DB == nil {
		return false
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}
	return sqlDB.Ping() == nil
}

// GetSQLitePath 获取SQLite地址
func GetSQLitePath(cfg *Config) string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("获取用户配置目录失败: %v", err)
	}
	appDir := filepath.Join(userConfigDir, cfg.App.Name)
	os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, cfg.Database.SQLite.Path)
}
