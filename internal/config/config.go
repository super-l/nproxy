package config

import (
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

type softConfig struct {
	System struct {
		LogLevel string `yaml:"logLevel"` // 日志级别
	} `yaml:"system"`

	Rpc struct {
		Ip       string `yaml:"ip"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"rpc"`

	Proxy struct {
		ApiRate int `yaml:"apiRate"`
	} `yaml:"proxy"`

	Db struct {
		Type  string `yaml:"type"`
		Mysql struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"mysql"`

		Sqlite struct {
			DbPath string `yaml:"dbpath"`
		} `yaml:"sqlite"`
	} `yaml:"db"`
}

func (s softConfig) IsEmpty() bool {
	return reflect.DeepEqual(s, softConfig{})
}

var sconfig softConfig

func InitConfig() error {
	conf, err := readYamlConfig("config.yaml")
	if err != nil {
		return err
	}
	sconfig = conf
	return err
}

func GetConfig() softConfig {
	return sconfig
}

func readYamlConfig(path string) (config softConfig, err error) {
	conf := softConfig{}
	if f, err := os.Open(path); err != nil {
		return config, err
	} else {
		yaml.NewDecoder(f).Decode(&conf)
	}
	return conf, nil
}
