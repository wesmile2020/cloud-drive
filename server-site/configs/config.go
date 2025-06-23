package configs

import (
	"cloud-drive/utils"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServerConfig 服务器配置结构
type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

// Config 总配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(pathUtil *utils.PathUtil) (*Config, error) {
	cfgDir := filepath.Join(pathUtil.GetRootDir(), "configs")
	cfgPath := filepath.Join(cfgDir, "config.yaml")

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置文件
		defaultConfig := Config{
			Server: ServerConfig{
				Port: "8080",
				Mode: "debug",
			},
			Database: DatabaseConfig{
				DSN: "cloud-drive.db",
			},
		}

		file, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config: %w", err)
		}

		utils.CreateDir(cfgDir)
		if err := os.WriteFile(cfgPath, file, 0644); err != nil {
			return nil, fmt.Errorf("failed to write default config file: %w", err)
		}

		return &defaultConfig, nil
	}

	file, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
