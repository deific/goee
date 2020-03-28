package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

// 定义程序配置类型
type Config struct {
	Port    string   `yaml:"port"`
	Host    string   `yaml:"host"`
	AppName string   `yaml:"app-name"`
	Statics []string `yaml:"statics"`
}

// 加载配置文件
func LoadConfig(filename string) Config {
	// 判断配置文件是否存在
	_, err := os.Lstat(filename)
	if os.IsNotExist(err) {
		// 返回默认值
		defaultConf := Config{
			Port:    "9000",
			Host:    "localhost",
			AppName: "goee web framework",
			Statics: []string{"static"},
		}
		log.Printf("config file:%s not exist,use default config data:%s", filename, defaultConf)
		return defaultConf
	}

	fByte, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("load config file [%s]error:/n%s", filename, err)
		return Config{}
	}
	// 声明对象
	var config Config
	err = yaml.Unmarshal(fByte, &config)
	if err != nil {
		log.Fatalf("load config file [%s]error:/n%s", filename, err)
	}
	log.Printf("use config file:%s config data:%s", filename, config)
	return config
}
