package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type SectionInternal interface{}

type APICredential struct {
	AccessKey string
	SecretKey string
}

type FileConfig struct {
	FilePath string `mapstructure:"file_path"`
}

type ArrayConfig struct {
	Items string `yaml:"items"`
}

type MapConfig struct {
	Items string `yaml:"items"`
}

type Config struct {
	configPath string
	Internal   SectionInternal `yaml:"internal"`
}

func (c *Config) configureViper() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("$HOME/.")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
}

func (c *Config) loadConf() (SectionInternal, error) {
	c.configureViper()
	if err := c.readConf(); err != nil {
		return nil, err
	}

	configSection := c.Internal
	err := viper.Unmarshal(&configSection)
	if err != nil {
		log.Printf("unable to decode into config struct: %v", err)
		return nil, err
	}

	return configSection, nil
}

func (c *Config) readConfFromFile() error {
	if c.configPath != "" {
		content, err := os.ReadFile(c.configPath)
		if err != nil {
			log.Printf("File does not exist: %s", c.configPath)
			return err
		}
		if err = viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			log.Printf("Can't read config file: %s | error:%s", c.configPath, err)
			return err
		}
	} else {
		if err := viper.MergeInConfig(); err == nil {
			log.Printf("Using config file:%v", viper.ConfigFileUsed())
		} else {
			log.Printf("Config file not found.")
		}
	}
	return nil
}

func (c *Config) readConf() error {
	if c.configPath != "" {
		content, err := os.ReadFile(c.configPath)
		if err != nil {
			log.Printf("File does not exist: %s | error:%s", c.configPath, err)
			return err
		}
		// load default config
		if err = viper.MergeConfig(bytes.NewBuffer(content)); err != nil {
			log.Printf("Can't merge config viper: %s", err)
			return err
		}
	}

	if err := c.readConfFromFile(); err != nil {
		return err
	}

	return nil
}

func NewConfig(path string, internalSection interface{}) *Config {
	conf := Config{
		Internal:   internalSection,
		configPath: path,
	}
	cfg, err := conf.loadConf()
	if err != nil {
		log.Fatalf("Load yaml config file error: '%v'", err)
		return nil
	}
	conf.Internal = cfg
	return &conf
}

func (s *FileConfig) GetValue() string {
	apiKey, err := os.ReadFile(s.FilePath)
	if err != nil {
		log.Panicf(fmt.Sprintf("Error to read file in path %v with error: %v", s.FilePath, err))
	}
	return strings.TrimSpace(string(apiKey))
}

func (s *ArrayConfig) GetItems() []string {
	return strings.Split(s.Items, ",")
}

func (s *MapConfig) GetKeyValues() map[string]string {
	result := make(map[string]string)
	for _, item := range strings.Split(s.Items, ",") {
		kv := strings.Split(strings.TrimSpace(item), ":")
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}

func (s *FileConfig) GetAPICredentialValue() *APICredential {
	dump := strings.Split(strings.TrimSpace(s.GetValue()), " ")
	if len(dump) < 2 {
		log.Panic("invalid secret file, expected 2 string separated by space")
	}

	res := APICredential{}
	res.AccessKey = dump[0]
	res.SecretKey = dump[1]
	return &res
}

type TimeDuration string

func (t TimeDuration) Duration() time.Duration {
	timeDuration, err := time.ParseDuration(string(t))
	if err != nil {
		log.Fatalf("Error parsing duration: %v", err)
	}

	return timeDuration
}
