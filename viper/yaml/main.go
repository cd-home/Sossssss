package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func ReadYaml(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(viper.GetString("HOST"))

	// 监听文件改动
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println(e.Name)
	})
	// 会启协程监听
	viper.WatchConfig()
}

func main() {
	ReadYaml("./yaml/")
}
