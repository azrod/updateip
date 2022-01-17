package config

import (
	"os"
	"path/filepath"

	uip_aws "github.com/azrod/updateip/pkg/providers/aws"
	uip_cloudflare "github.com/azrod/updateip/pkg/providers/cloudflare"
	uip_ovh "github.com/azrod/updateip/pkg/providers/ovh"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

type CFGMetrics struct {
	// Prefix is the prefix used for all metrics.
	Prefix string `yaml:"prefix" env:"METRICS_PREFIX"`
	// Enable is a boolean that determines if metrics are enabled.
	Enable bool `yaml:"enable" env:"METRICS_ENABLE"`
	// Host is the endpoint used for the metrics server.
	Host string `yaml:"host" env:"METRICS_HOST"`
	// Port is the port used for the metrics server.
	Port int `yaml:"port" env:"METRICS_PORT"`
	// Path is the path used for the metrics server.
	Path string `yaml:"path" env:"METRICS_PATH"`
	// Enable logging of metrics.
	Logging bool `yaml:"logging" env:"METRICS_LOGGING"`
}

type CFG struct {
	Log struct {
		// trace is the value used for the trace level field.
		// debug is the value used for the debug level field.
		// info is the value used for the info level field.
		// warn is the value used for the warn level field.
		// error is the value used for the error level field.
		// fatal is the value used for the fatal level field.
		// LevelPanicValue is the value used for the panic level field.
		Level    string `yaml:"level" env:"LOG_LEVEL"`
		Humanize bool   `yaml:"humanize" env:"LOG_HUMANIZE"`
	}
	Metrics   CFGMetrics `yaml:"metrics" env:"METRICS"`
	Providers struct {
		AWSAccount struct {
			Enable bool               `yaml:"enable" env:"AWS_ACCOUNT_ENABLE"`
			Secret uip_aws.PawsSecret `yaml:"secret"`
			Record uip_aws.PawsRecord `yaml:"record"`
		} `yaml:"aws"`
		OVHAccount struct {
			Enable bool               `yaml:"enable" env:"OVH_ACCOUNT_ENABLE"`
			Secret uip_ovh.PovhSecret `yaml:"secret"`
			Record uip_ovh.PovhRecord `yaml:"record"`
		} `yaml:"ovh"`
		CLOUDFLAREAccount struct {
			Enable bool                             `yaml:"enable" env:"CLOUDFLARE_ACCOUNT_ENABLE"`
			Secret uip_cloudflare.PCloudflareSecret `yaml:"secret"`
			Record uip_cloudflare.PCloudflareRecord `yaml:"record"`
		} `yaml:"cloudflare"`
	} `yaml:"providers"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config CFG, err error) {

	fp := config.getConfigFileStruct()

	// open yaml file
	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		log.Error().Err(err).Str("path", fp).Msg("Failed to open config file")
	} else {
		defer func() {
			if err := f.Close(); err != nil {
				log.Error().Err(err).Msg("Error closing file")
			}
		}()

		// read yaml file
		err = yaml.NewDecoder(f).Decode(&config)
		if err != nil {
			log.Error().Err(err).Str("path", fp).Msg("Failed to decode config file")
		}
	}

	if err := env.Parse(&config); err != nil {
		log.Error().Err(err).Msg("Failed to parse environment variables")
	}

	return
}

func (c CFG) getConfigFileStruct() (filePath string) {

	var (
		dir  = "/config"
		file = "config.yaml"
		ok   bool
	)

	if _, ok = os.LookupEnv("PATH_CONFIG_DIRECTORY"); ok {
		dir = os.Getenv("PATH_CONFIG_DIRECTORY")
		log.Info().Msgf("Using config directory: %s", dir)
	}

	if _, ok = os.LookupEnv("PATH_CONFIG_FILE"); ok {
		file = os.Getenv("PATH_CONFIG_FILE")
		log.Info().Msgf("Using config file: %s", file)
	}

	return filepath.Clean(dir + "/" + file)
}
