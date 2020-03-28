package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
	if filename == "" {
		// 返回默认值
		return Config{
			Port:    "9000",
			Host:    "localhost",
			AppName: "goee web framework",
			Statics: []string{"static"},
		}
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
	return config
}
