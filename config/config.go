package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var Cfg *Settings

func InitConfig(env string) {
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join("config", "env"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	log.Printf("成功加载 %s 环境配置", env)
}
