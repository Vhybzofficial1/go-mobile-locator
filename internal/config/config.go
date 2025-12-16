package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	SQLite struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"sqlite"`

	Gorm struct {
		LogLevel              int           `mapstructure:"logLevel"`
		MaxIdleConns          int           `mapstructure:"maxIdleConns"`
		MaxOpenConns          int           `mapstructure:"maxOpenConns"`
		ConnMaxLifetime       string        `mapstructure:"connMaxLifetime"`
		ParsedConnMaxLifetime time.Duration `mapstructure:"-"`
	} `mapstructure:"gorm"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
}

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	App      AppConfig      `mapstructure:"app"`
}

var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig(configPath string) error {
	GlobalConfig = &Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// 添加配置文件路径
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
		v.AddConfigPath("../configs")
		if exe, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exe)
			v.AddConfigPath(exeDir)
			v.AddConfigPath(filepath.Join(exeDir, "configs"))
		}
	}
	// 设置默认值
	setDefaultValues(v)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("配置文件未找到，使用默认配置")
		} else {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}
	// 解析到结构体
	if err := v.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}
	// 解析 GORM 时间字段
	if err := parseTimeFields(); err != nil {
		return err
	}
	// 验证配置
	if err := validateConfig(); err != nil {
		return err
	}
	log.Println("配置加载成功")
	log.Printf("应用: %s v%s", GlobalConfig.App.Name, GlobalConfig.App.Version)
	log.Printf("数据库: %s", GlobalConfig.Database.SQLite.Path)
	return nil
}

// setDefaultValues 设置默认值
func setDefaultValues(v *viper.Viper) {
	v.SetDefault("app.name", "MobileLocator")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.debug", true)
	v.SetDefault("database.sqlite.path", "app.db")
	v.SetDefault("database.gorm.logLevel", 1)
	v.SetDefault("database.gorm.maxIdleConns", 10)
	v.SetDefault("database.gorm.maxOpenConns", 100)
	v.SetDefault("database.gorm.connMaxLifetime", "1h")
}

// parseTimeFields 解析 GORM 时间字段
func parseTimeFields() error {
	if GlobalConfig.Database.Gorm.ConnMaxLifetime != "" {
		duration, err := time.ParseDuration(GlobalConfig.Database.Gorm.ConnMaxLifetime)
		if err != nil {
			return fmt.Errorf("解析 connMaxLifetime 失败: %w", err)
		}
		GlobalConfig.Database.Gorm.ParsedConnMaxLifetime = duration
	}
	return nil
}

// validateConfig 验证配置
func validateConfig() error {
	if GlobalConfig.Database.SQLite.Path == "" {
		return fmt.Errorf("数据库路径不能为空")
	}
	if GlobalConfig.Database.Gorm.MaxIdleConns <= 0 {
		return fmt.Errorf("maxIdleConns 必须大于 0")
	}
	if GlobalConfig.Database.Gorm.MaxOpenConns <= 0 {
		return fmt.Errorf("maxOpenConns 必须大于 0")
	}
	if GlobalConfig.Database.Gorm.LogLevel < 1 || GlobalConfig.Database.Gorm.LogLevel > 4 {
		return fmt.Errorf("logLevel 必须在 1-4 之间")
	}
	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if GlobalConfig == nil {
		log.Println("警告: 配置未初始化，使用默认值初始化")
		InitConfig("")
	}
	return GlobalConfig
}

// ReloadConfig 重新加载配置
func ReloadConfig() error {
	return InitConfig("")
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.App.Debug
}

// GetAppName 获取应用名称
func (c *Config) GetAppName() string {
	return c.App.Name
}

// GetDBPath 获取数据库路径
func (c *Config) GetDBPath() string {
	return c.Database.SQLite.Path
}
