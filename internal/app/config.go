package app

import (
	"github.com/joomcode/errorx"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

var Errors = errorx.NewNamespace("config")
var ParseErr = Errors.NewType("parse")
var FileErr = Errors.NewType("file")

type Config struct {
	Client      ClientConfig `yaml:"client"`
	Bot         BotConfig    `yaml:"bot"`
	FilterRegex string       `yaml:"filterRegex"`
}

type BotConfig struct {
	Token   string `yaml:"token"`
	Timeout int    `yaml:"timeout"`
}

type ClientConfig struct {
	ApiId              int         `yaml:"apiId"`
	ApiHash            string      `yaml:"apiHash"`
	Phone              string      `yaml:"phone"`
	SystemLanguageCode string      `yaml:"systemLanguageCode"`
	SystemVersion      string      `yaml:"systemVersion"`
	DeviceModel        string      `yaml:"deviceModel"`
	ApplicationVersion string      `yaml:"applicationVersion"`
	FilesDirectory     string      `yaml:"filesDirectory"`
	DatabaseDirectory  string      `yaml:"databaseDirectory"`
	UseFileDatabase    bool        `yaml:"useFileDatabase"`
	CheckCode          string      `yaml:"checkCode"`
	Password           string      `yaml:"password"`
	Proxy              ProxyConfig `yaml:"proxy"`
}

type ProxyConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

func LoadConfig(r io.Reader) (c *Config, err error) {
	d := yaml.NewDecoder(r)
	c = &Config{}
	err = d.Decode(c)
	if err != nil {
		err = ParseErr.Wrap(err, "config parse failed")
	}
	return
}

func LoadConfigFile(path string) (c *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		err = FileErr.Wrap(err, "failed to open file: "+path)
		return
	}
	c, err = LoadConfig(file)
	return
}
