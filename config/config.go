package config

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-zxb/fuxi/internal/model"
	"github.com/spf13/viper"
	"os"
	"time"
)

//go:embed config.yaml
var fileYaml string

var Conf *Config

func GetConfig() *Config {
	return Conf
}

type Config struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Gin    Gin    `mapstructure:"gin" json:"gin" yaml:"gin"`
	Jwt    Jwt    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	GPT    GPT    `mapstructure:"gpt" json:"gpt" yaml:"gpt"`
	viper  *viper.Viper
}

func NewConfig(filePath string) (*Config, error) {
	Conf = &Config{}
	return Conf.InitConfig(filePath)
}

type System struct {
	Name        string `mapstructure:"name" json:"name" yaml:"name"`
	Version     string `mapstructure:"version" json:"version" yaml:"version"`
	Description string `mapstructure:"description" json:"description" yaml:"description"`
}

type Gin struct {
	Mode  string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Host  string `mapstructure:"host" json:"host" yaml:"host"`
	Port  int    `mapstructure:"port" json:"port" yaml:"port"`
	Debug bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
}

type Mysql struct {
	Host                      string `mapstructure:"host" json:"host" yaml:"host"`
	Port                      int    `mapstructure:"port" json:"port" yaml:"port"`
	User                      string `mapstructure:"user" json:"user" yaml:"user"`
	Password                  string `mapstructure:"password" json:"password" yaml:"password"`
	Database                  string `mapstructure:"database" json:"database" yaml:"database"`
	MaxIdleConns              int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns              int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogLevel                  int    `mapstructure:"Log_level" json:"Log_level" yaml:"Log_level"`
	Charset                   string `mapstructure:"charset" json:"charset" yaml:"charset"`
	TimeZone                  string `mapstructure:"time_zone" json:"time_zone" yaml:"time_zone"`
	ParseTime                 bool   `mapstructure:"parse_time" json:"parse_time" yaml:"parse_time"`
	Colorful                  bool   `mapstructure:"colorful" json:"colorful" yaml:"colorful"`
	IgnoreRecordNotFoundError bool   `mapstructure:"ignore_record_not_found_error" json:"ignore_record_not_found_error" yaml:"ignore_record_not_found_error"`
	ParameterizedQueries      bool   `mapstructure:"parameterized_queries" json:"parameterized_queries" yaml:"parameterized_queries"`
}

type Redis struct {
	UserName     string        `mapstructure:"user_name" json:"user_name" yaml:"user_name"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password"`
	Addr         string        `mapstructure:"addr" json:"addr" yaml:"addr"`
	DB           int           `mapstructure:"db" json:"db" yaml:"db"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
}

type Jwt struct {
	SecretKey         string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	ExpirationSeconds int64  `mapstructure:"expiration_seconds" json:"expiration_seconds" yaml:"expiration_seconds"`
	Issuer            string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}

type GPTModel struct {
	Model   string `mapstructure:"model" json:"model" yaml:"model"`
	ApiKey  string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`
	BaseURL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"`
}

type GPT struct {
	Prompt          string          `mapstructure:"prompt" json:"prompt" yaml:"prompt"`
	ChatGPTPlatform model.ModelType `mapstructure:"chat_gpt_platform" json:"chat_gpt_platform" yaml:"chat_gpt_platform"`
	Temperature     float64         `mapstructure:"temperature" json:"temperature" yaml:"temperature"`
	Kimi            GPTModel        `mapstructure:"kimi" json:"kimi" yaml:"kimi"`
	DeepSeek        GPTModel        `mapstructure:"deep_seek" json:"deep_seek" yaml:"deep_seek"`
}

func (c *Config) InitConfig(filePath string) (*Config, error) {
	// 如果给的配置文件不支持则导出一份默认配置文件到本地
	path := filePath
	// 判断配置文件是否存在
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		//用户给的配置文件不存在 则自动使用默认配置文件
		// 判断默认配置文件是否存在 不存在则创建一份
		_, err = os.Stat("config/config_dev.yaml")
		if os.IsNotExist(err) {
			err = os.MkdirAll("config", os.ModePerm)
			if err != nil {
				return nil, errors.New("配置文件初始化失败")
			}
			err = os.WriteFile("config/config_dev.yaml", []byte(fileYaml), os.ModePerm)
			if err != nil {
				return nil, errors.New("配置文件初始化失败")
			}
			path = "config/config_dev.yaml"
		} else {
			path = "config/config_dev.yaml"
		}
	}

	var conf *viper.Viper
	conf = viper.New()
	conf.SetConfigFile(path)
	err = conf.ReadInConfig()
	if err != nil {
		panic(err)
		return nil, err
	}

	err = conf.Unmarshal(c)
	if err != nil {
		panic(err)
	}

	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已更改: ", e.Name)
		if err = conf.Unmarshal(c); err != nil {
			fmt.Println(err)
		}
	})
	c.viper = conf
	return c, nil
}
