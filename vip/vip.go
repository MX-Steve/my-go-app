package vip

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func GetIniData(key string) string {
	config := viper.New()
	dir, _ := os.Getwd()
	config.AddConfigPath(dir)
	config.SetConfigName("config")
	config.SetConfigType("ini")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("找不到配置文件")
			return ""
		} else {
			log.Println("配置文件出错")
			return ""
		}
	}
	value := config.GetString(key)
	return value
}
