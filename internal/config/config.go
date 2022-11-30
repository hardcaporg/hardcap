package config

import (
    "context"
    "github.com/hardcaporg/hardcap/internal/ctxval"
    "os"

	"github.com/ilyakaznacheev/cleanenv"
)

var config struct {
	App struct {
		HttpListenAddress string `env:"HTTP_LISTEN_ADDRESS" env-default:":8000" env-description:"listen address and port for the REST API"`
        RpcListenAddress string `env:"RPC_LISTEN_ADDRESS" env-default:":8007" env-description:"listen address and port for the RPC API"`
        AdvertisedAddress string `env:"ADV_ADDRESS" env-default:"localhost:8000" env-description:"advertised address in templates"`
	} `env-prefix:"APP_"`
	Logging struct {
		Level    string `env:"LEVEL" env-default:"info" env-description:"logger level (trace, debug, info, warn, error, fatal, panic)"`
		Stdout   bool   `env:"STDOUT" env-default:"true" env-description:"logger standard output"`
		MaxField int    `env:"MAX_FIELD" env-default:"0" env-description:"logger maximum field length (dev only)"`
	} `env-prefix:"LOGGING_"`
	Prometheus struct {
		Port int    `env:"PORT" env-default:"9000" env-description:"prometheus HTTP port"`
		Path string `env:"PATH" env-default:"/metrics" env-description:"prometheus metrics path"`
	} `env-prefix:"PROMETHEUS_"`
}

// Config shortcuts
var (
	Application = &config.App
	Prometheus  = &config.Prometheus
	Logging     = &config.Logging
)

// Initialize loads configuration from provided .env files, the first existing file wins.
func Initialize(configFiles ...string) {
	var loaded bool
	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			// if config file exists, load it (also loads environmental variables)
			err := cleanenv.ReadConfig(configFile, &config)
			if err != nil {
				panic(err)
			}
			loaded = true
		}
	}

	if !loaded {
		// otherwise use only environmental variables instead
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			panic(err)
		}
	}
}

func PrintConfig(ctx context.Context) {
    logger := ctxval.Logger(ctx)
    logger.Trace().Msgf("Config: %+v", config)
}