package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server        ServerConfig
	Redis         RedisConfig
	RedisSentinel RedisSentinelConfig
	Logger        Logger
	Service       ServiceConfig
	JwtConfig     JwtConfig
	Cors          CorsConfig
	SecretToken   SecretToken
	RateLimit     RateLimit
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
	ServiceName       string
	TimeConvert       string
	JwtCertFile       string
	WhiteList         string
}

// Logger config
type Logger struct {
	DevMode  bool
	Encoder  string
	Encoding string
	LogLevel string
}

type RateLimit struct {
	Limit  int
	Expire int
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}
type RedisSentinelConfig struct {
	Addr             string
	RouteByLatency   bool // Allows routing read-only commands to the closest master or slave node.
	RouteRandomly    bool // Allows routing read-only commands to the random master or slave node.
	Username         string
	Password         string
	SentinelUsername string
	SentinelPassword string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	PoolFIFO         bool // PoolFIFO uses FIFO mode for each node connection pool GET/PUT (default LIFO).
	PoolSize         int
	MinIdleConns     int
	MaxRetries       int
	MinRetryBackoff  time.Duration
	MaxRetryBackoff  time.Duration
	DialTimeout      time.Duration
	PoolTimeout      time.Duration
	MasterName       string
}

// Service config
type ServiceConfig struct {
	CustomerServiceUrl string
	OrderServiceUrl    string
}

// JwtKey config
type JwtConfig struct {
	SecretKey string
	PublicKey []string
	TTL       int
}

// Cors config
type CorsConfig struct {
	AllowOrigins []string
	AllowHeaders []string
}

type SecretToken string

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	if cfg.RateLimit.Limit == 0 {
		cfg.RateLimit.Limit = 10
	}

	if cfg.RateLimit.Expire == 0 {
		cfg.RateLimit.Expire = 60
	}

	return cfg, nil
}

// Get config
func InitConfig(env string) (*Config, error) {
	var configPath string
	switch env {
	case "qc":
		configPath = "./config/qc"
	case "staging":
		configPath = "./config/staging"
	case "prod":
		configPath = "./config/prod"
	default:
		configPath = "./config/local"
	}

	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
