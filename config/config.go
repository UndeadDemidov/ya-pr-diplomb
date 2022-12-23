package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ErrConfigFileNotFound error.
var ErrConfigFileNotFound = errors.New("config file not found")

// App includes all set of configurations.
type App struct {
	Server
	Postgres
	Metrics
	Logger
	Jaeger
}

// Server contains configurations variables for GRPC and HTTP servers.
type Server struct {
	Mode              string
	AppVersion        string
	Port              string
	PprofPort         string
	JwtSecretKey      string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CtxDefaultTimeout time.Duration
	SSL               bool
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

// Logger contains configurations variables for logging.
type Logger struct {
	Development bool
	Level       string
}

// Postgres contains configurations variables for postgres database.
type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  bool
	Driver   string
}

// Session contains time to live value for sessions duration.
type Session struct {
	TTL int64
}

// Metrics contains configurations variables for metrics service.
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger contains configurations variables for jaeger service.
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Load loading config file from given path.
// ToDo add arguments parsing.
func Load(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var cfgFileNotFoundErr *viper.ConfigFileNotFoundError
		if errors.As(err, &cfgFileNotFoundErr) {
			return nil, ErrConfigFileNotFound
		}

		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	return v, nil
}

// Parse parsing config file.
func Parse(v *viper.Viper) (*App, error) {
	var c App

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &c, nil
}

// Get gets app configuration from given path.
func Get(configPath string) (*App, error) {
	cfgFile, err := Load(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := Parse(cfgFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetPath gets config path for local or docker.
func GetPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
