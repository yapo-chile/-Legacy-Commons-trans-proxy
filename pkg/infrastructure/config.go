package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// RuntimeConfig config to start the app
type RuntimeConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8080"`
	// Profiling if the service should add profiling endpoints with net/http/pprof
	Profiling bool   `env:"PROFILING" envDefault:"true"`
	APIKey    string `env:"API_KEY" envDefault:"test"`
}

// Addresss return the address of the service with host and port
func (c RuntimeConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoggerConf holds configuration for logging
// LogLevel definition:
//   0 - Debug
//   1 - Info
//   2 - Warning
//   3 - Error
//   4 - Critic
type LoggerConf struct {
	SyslogIdentity string `env:"SYSLOG_IDENTITY"`
	SyslogEnabled  bool   `env:"SYSLOG_ENABLED" envDefault:"false"`
	StdlogEnabled  bool   `env:"STDLOG_ENABLED" envDefault:"true"`
	LogLevel       int    `env:"LOG_LEVEL" envDefault:"0"`
}

// PrometheusConf holds configuration to report to Prometheus
type PrometheusConf struct {
	Port    string `env:"PORT" envDefault:"8877"`
	Enabled bool   `env:"ENABLED" envDefault:"false"`
}

// ProfileConf holds configuration to send http request to profile
// CorsConf holds cors headers
type CorsConf struct {
	Enabled bool   `env:"ENABLED" envDefault:"false"`
	Origin  string `env:"ORIGIN" envDefault:"*"`
	Methods string `env:"METHODS" envDefault:"GET, OPTIONS"`
	Headers string `env:"HEADERS" envDefault:"Accept,Content-Type,Content-Length,If-None-Match,Accept-Encoding,User-Agent"`
}

// GetHeaders return map of cors used
func (cc CorsConf) GetHeaders() map[string]string {
	if !cc.Enabled {
		return map[string]string{}
	}

	return map[string]string{
		"Origin":  cc.Origin,
		"Methods": cc.Methods,
		"Headers": cc.Headers,
	}
}

// InBrowserCacheConf Used to handle browser cache
type InBrowserCacheConf struct {
	Enabled bool `env:"ENABLED" envDefault:"false"`
	// Cache max age in secs(use browser cache)
	MaxAge time.Duration `env:"MAX_AGE" envDefault:"720h"`
	Etag   int64
}

// InitEtag use current epoc to config etag
func (chc *InBrowserCacheConf) InitEtag() int64 {
	chc.Etag = time.Now().Unix()
	return chc.Etag
}

// TransConf transaction server connection.
type TransConf struct {
	// AllowedCommands is a list with one or more trans commands, separated by '|'
	// that indicates the allowed commands to be sent by this service
	AllowedCommands string `env:"COMMANDS" envDefault:"transinfo"`
	// Host is the host of the trans Server
	Host string `env:"HOST" envDefault:"localhost"`
	// Port is the port of the trans server
	Port int `env:"PORT" envDefault:"20005"`
	// Timeout wait time before a request times out
	Timeout int `env:"TIMEOUT" envDefault:"15"`
	// RetryAfter wait time between reconnection to the trans server
	RetryAfter int `env:"RETRY" envDefault:"5"`
}

// Config holds all configuration for the service
type Config struct {
	Trans              TransConf          `env:"TRANS_"`
	PrometheusConf     PrometheusConf     `env:"PROMETHEUS_"`
	LoggerConf         LoggerConf         `env:"LOGGER_"`
	Runtime            RuntimeConfig      `env:"APP_"`
	CorsConf           CorsConf           `env:"CORS_"`
	InBrowserCacheConf InBrowserCacheConf `env:"BROWSER_CACHE_"`
}

// LoadFromEnv loads the config data from the environment variables
func LoadFromEnv(data interface{}) {
	load(reflect.ValueOf(data), "", "")
}

// valueFromEnv lookup the best value for a variable on the environment
func valueFromEnv(envTag, envDefault string) string {
	// Maybe it's a secret and <envTag>_FILE points to a file with the value
	// https://rancher.com/docs/rancher/v1.6/en/cattle/secrets/#docker-hub-images
	if fileName, ok := os.LookupEnv(fmt.Sprintf("%s_FILE", envTag)); ok {
		b, err := ioutil.ReadFile(fileName) // nolint: gosec
		if err == nil {
			return string(b)
		}
		fmt.Print(err)
	}
	// The value might be set directly on the environment
	if value, ok := os.LookupEnv(envTag); ok {
		return value
	}
	// Nothing to do, return the default
	return envDefault
}

// load the variable defined in the envTag into Value
func load(conf reflect.Value, envTag, envDefault string) {
	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// Only attempt to set writeable variables
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value := valueFromEnv(envTag, envDefault)
			// Print message if config is missing
			if envTag != "" && value == "" && !strings.HasSuffix(envTag, "_") {
				fmt.Printf("Config for %s missing\n", envTag)
			}
			switch reflectedConf.Kind() {
			case reflect.Struct:
				// Recursively load inner struct fields
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
			// Here for each type we should make a cast of the env variable and then set the value
			case reflect.String:
				reflectedConf.SetString(value)
			case reflect.Int:
				if value, err := strconv.Atoi(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case reflect.Bool:
				if value, err := strconv.ParseBool(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			}
		}
	}
}
