package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type AppConfig struct {
	ServiceName string   `yaml:"service_name"`
	Language    string   `yaml:"language"`
	Mode        string   `yaml:"mode"`
	DB          DB       `mapstructure:"db"`
	Services    Services `yaml:"services"`
	Core        Core     `yaml:"core"`
}

type Auth struct {
	JWTSecretKey string       `mapstructure:"jwt_secret_key"`
	TTL          TimeDuration `mapstructure:"ttl"`
	OtpPeriod    uint8        `mapstructure:"otp_period"`
	OtpLength    uint8        `mapstructure:"otp_length"`
}

type Core struct {
	Http Http `mapstructure:"http"`
	Auth Auth `yaml:"auth"`
}

type Http struct {
	Address string `yaml:"address"`
	Port    uint16 `yaml:"port"`
	Url     string `yaml:"url"`
}

type DB struct {
	Postgres  Postgres `mapstructure:"postgres"`
	Redis     Redis    `yaml:"redis"`
	RunSeeder bool     `mapstructure:"run_seeder"`
}

type Postgres struct {
	Host               string        `yaml:"host"`
	AppName            string        `yaml:"appName"`
	Port               int           `yaml:"port"`
	Username           string        `yaml:"username"`
	Password           string        `mask:"filled" yaml:"password"`
	Name               string        `yaml:"name"`
	Driver             string        `yaml:"driver"`
	AutoMigration      bool          `mapstructure:"auto_migration"`
	Ssl                string        `yaml:"ssl"`
	MigrationsURL      string        `mapstructure:"migrations_url"`
	TransactionTimeout time.Duration `yaml:"transaction_timeout" mapstructure:"transaction_timeout"`
	MaxIdleConnection  int           `yaml:"max_idle_connection" mapstructure:"max_idle_connection"`
	MaxOpenConnection  int           `yaml:"max_open_connection" mapstructure:"max_open_connection"`
	ConnMaxLifetime    time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	TraceStacks        bool          `yaml:"trace_stacks" mapstructure:"trace_stacks"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `mask:"filled" yaml:"password"`
	DB       int    `yaml:"db"`
}

type Services struct{}

func LoadConfig(configPath string) *AppConfig {
	conf := NewConfig(configPath, &AppConfig{})
	configJson, err := json.Marshal(conf.Internal.(*AppConfig))
	if err != nil {
		panic(fmt.Sprintf("Can't make the json the config file:%v", err))
	}
	fmt.Printf("config value: %v", string(configJson))

	return conf.Internal.(*AppConfig)
}
