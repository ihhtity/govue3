package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件已更改:", e.Name)
		if err := viper.Unmarshal(&Cfg); err != nil {
			log.Printf("配置热更新失败: %v", err)
		}
	})
}
