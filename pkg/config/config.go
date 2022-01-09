package config

import (
	"os"
	"path/filepath"
	"strconv"

	uip_aws "github.com/azrod/updateip/pkg/providers/aws"
	uip_cloudflare "github.com/azrod/updateip/pkg/providers/cloudflare"
	uip_ovh "github.com/azrod/updateip/pkg/providers/ovh"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

type CFGMetrics struct {
	// Prefix is the prefix used for all metrics.
	Prefix string `yaml:"prefix"`
	// Enable is a boolean that determines if metrics are enabled.
	Enable bool `yaml:"enable"`
	// Host is the endpoint used for the metrics server.
	Host string `yaml:"host"`
	// Port is the port used for the metrics server.
	Port int `yaml:"port"`
	// Path is the path used for the metrics server.
	Path string `yaml:"path"`
	// Enable logging of metrics.
	Logging bool `yaml:"logging"`
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
		Level    string `yaml:"level"`
		Humanize bool   `yaml:"humanize"`
	}
	Metrics   CFGMetrics `yaml:"metrics"`
	Providers struct {
		AWSAccount struct {
			Enable bool               `yaml:"enable"`
			Secret uip_aws.PawsSecret `yaml:"secret"`
			Record uip_aws.PawsRecord `yaml:"record"`
		} `yaml:"aws"`
		OVHAccount struct {
			Enable bool               `yaml:"enable"`
			Secret uip_ovh.PovhSecret `yaml:"secret"`
			Record uip_ovh.PovhRecord `yaml:"record"`
		} `yaml:"ovh"`
		CLOUDFLAREAccount struct {
			Enable bool                             `yaml:"enable"`
			Secret uip_cloudflare.PCloudflareSecret `yaml:"secret"`
			Record uip_cloudflare.PCloudflareRecord `yaml:"record"`
		} `yaml:"cloudflare"`
	} `yaml:"providers"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config CFG, err error) {

	fp := config.getConfigFileStruct()

	// open yaml file
	f, err := os.Open(fp)
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

	config.readEnvVars()

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

func (c CFG) readEnvVars() {
	var (
		ok  bool
		err error
	)

	// try to read from environment variables
	// env AWS_ACCOUNT_ENABLE
	if _, ok = os.LookupEnv("AWS_ACCOUNT_ENABLE"); ok {
		log.Info().Msg("AWS_ACCOUNT_ENABLE found in environment variables")
		c.Providers.AWSAccount.Enable, err = strconv.ParseBool(os.Getenv("AWS_ACCOUNT_ENABLE"))

		if err != nil {
			log.Error().Err(err).Msg("Failed to parse AWS_ACCOUNT_ENABLE")
		}
	}

	if _, ok = os.LookupEnv("AWS_ACCESS_KEY_ID"); ok {
		log.Info().Msg("Reading AWS_ACCESS_KEY_ID from environment variables")
		c.Providers.AWSAccount.Secret.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}

	if _, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY"); ok {
		log.Info().Msg("Reading AWS_SECRET_ACCESS_KEY from environment variables")
		c.Providers.AWSAccount.Secret.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	if _, ok = os.LookupEnv("AWS_REGION"); ok {
		log.Info().Msg("Reading AWS_REGION from environment variables")
		c.Providers.AWSAccount.Secret.Region = os.Getenv("AWS_REGION")
	}

	if _, ok = os.LookupEnv("AWS_HOSTED_ZONE_ID"); ok {
		log.Info().Msg("Reading AWS_HOSTED_ZONE_ID from environment variables")
		c.Providers.AWSAccount.Record.HostedZoneID = os.Getenv("AWS_HOSTED_ZONE_ID")
	}

	if _, ok = os.LookupEnv("AWS_RECORD_NAME"); ok {
		log.Info().Msg("Reading AWS_RECORD_NAME from environment variables")
		c.Providers.AWSAccount.Record.Name = os.Getenv("AWS_HOSTED_ZONE_NAME")
	}

	// AWS_RECORD_TTL
	if _, ok = os.LookupEnv("AWS_RECORD_TTL"); ok {
		log.Info().Msg("Reading AWS_RECORD_TTL from environment variables")
		ttl, err := strconv.Atoi(os.Getenv("AWS_RECORD_TTL"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("AWS_RECORD_TTL")).Msg("Failed to convert AWS_RECORD_TTL to int")
		} else {
			c.Providers.AWSAccount.Record.TTL = ttl
		}
	}

	// AWS_RECORD_DOMAIN
	if _, ok = os.LookupEnv("AWS_RECORD_DOMAIN"); ok {
		log.Info().Msg("Reading AWS_RECORD_DOMAIN from environment variables")
		c.Providers.AWSAccount.Record.Domain = os.Getenv("AWS_RECORD_DOMAIN")
	}

	// AWS_RECORD_COMMENT
	if _, ok = os.LookupEnv("AWS_RECORD_COMMENT"); ok {
		log.Info().Msg("Reading AWS_RECORD_COMMENT from environment variables")
		c.Providers.AWSAccount.Record.Comment = os.Getenv("AWS_RECORD_COMMENT")
	}

	// OVH_ACCOUNT_ENABLE
	if _, ok = os.LookupEnv("OVH_ACCOUNT_ENABLE"); ok {
		log.Info().Msg("OVH_ACCOUNT_ENABLE found in environment variables")
		c.Providers.OVHAccount.Enable, err = strconv.ParseBool(os.Getenv("OVH_ACCOUNT_ENABLE"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse OVH_ACCOUNT_ENABLE")
		}
	}

	// OVH_APPLICATION_KEY
	if _, ok = os.LookupEnv("OVH_APPLICATION_KEY"); ok {
		log.Info().Msg("Reading OVH_APPLICATION_KEY from environment variables")
		c.Providers.OVHAccount.Secret.ApplicationKey = os.Getenv("OVH_APPLICATION_KEY")
	}

	// OVH_APPLICATION_SECRET
	if _, ok = os.LookupEnv("OVH_APPLICATION_SECRET"); ok {
		log.Info().Msg("Reading OVH_APPLICATION_SECRET from environment variables")
		c.Providers.OVHAccount.Secret.ApplicationSecret = os.Getenv("OVH_APPLICATION_SECRET")
	}

	// OVH_CONSUMER_KEY
	if _, ok = os.LookupEnv("OVH_CONSUMER_KEY"); ok {
		log.Info().Msg("Reading OVH_CONSUMER_KEY from environment variables")
		c.Providers.OVHAccount.Secret.ConsumerKey = os.Getenv("OVH_CONSUMER_KEY")
	}

	// OVH_REGION
	if _, ok = os.LookupEnv("OVH_REGION"); ok {
		log.Info().Msg("Reading OVH_REGION from environment variables")
		c.Providers.OVHAccount.Secret.Region = os.Getenv("OVH_REGION")
	}

	// OVH_RECORD_NAME
	if _, ok = os.LookupEnv("OVH_RECORD_NAME"); ok {
		log.Info().Msg("Reading OVH_RECORD_NAME from environment variables")
		c.Providers.OVHAccount.Record.Name = os.Getenv("OVH_RECORD_NAME")
	}

	// OVH_RECORD_TTL
	if _, ok = os.LookupEnv("OVH_RECORD_TTL"); ok {
		log.Info().Msg("Reading OVH_RECORD_TTL from environment variables")
		ttl, err := strconv.Atoi(os.Getenv("OVH_RECORD_TTL"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("OVH_RECORD_TTL")).Msg("Failed to convert OVH_RECORD_TTL to int")
		} else {
			c.Providers.OVHAccount.Record.TTL = ttl
		}
	}

	// OVH_RECORD_ZONE
	if _, ok = os.LookupEnv("OVH_RECORD_ZONE"); ok {
		log.Info().Msg("Reading OVH_RECORD_ZONE from environment variables")
		c.Providers.OVHAccount.Record.Zone = os.Getenv("OVH_RECORD_ZONE")
	}

	// CLOUDFLARE_ACCOUNT_ENABLE
	if _, ok = os.LookupEnv("CLOUDFLARE_ACCOUNT_ENABLE"); ok {
		log.Info().Msg("CLOUDFLARE_ACCOUNT_ENABLE found in environment variables")
		c.Providers.CLOUDFLAREAccount.Enable, err = strconv.ParseBool(os.Getenv("CLOUDFLARE_ACCOUNT_ENABLE"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse CLOUDFLARE_ACCOUNT_ENABLE")
		}
	}

	// CLOUDFLARE_API_KEY
	if _, ok = os.LookupEnv("CLOUDFLARE_API_KEY"); ok {
		log.Info().Msg("Reading CLOUDFLARE_API_KEY from environment variables")
		c.Providers.CLOUDFLAREAccount.Secret.APIKey = os.Getenv("CLOUDFLARE_API_KEY")
	}

	// CLOUDFLARE_EMAIL
	if _, ok = os.LookupEnv("CLOUDFLARE_EMAIL"); ok {
		log.Info().Msg("Reading CLOUDFLARE_EMAIL from environment variables")
		c.Providers.CLOUDFLAREAccount.Secret.Email = os.Getenv("CLOUDFLARE_EMAIL")
	}

	// CLOUDFLARE_RECORD_NAME
	if _, ok = os.LookupEnv("CLOUDFLARE_RECORD_NAME"); ok {
		log.Info().Msg("Reading CLOUDFLARE_RECORD_NAME from environment variables")
		c.Providers.CLOUDFLAREAccount.Record.Name = os.Getenv("CLOUDFLARE_RECORD_NAME")
	}

	// CLOUDFLARE_RECORD_DOMAIN
	if _, ok = os.LookupEnv("CLOUDFLARE_RECORD_DOMAIN"); ok {
		log.Info().Msg("Reading CLOUDFLARE_RECORD_DOMAIN from environment variables")
		c.Providers.CLOUDFLAREAccount.Record.Domain = os.Getenv("CLOUDFLARE_RECORD_DOMAIN")
	}

	// CLOUDFLARE_RECORD_ZONEID
	if _, ok = os.LookupEnv("CLOUDFLARE_RECORD_ZONEID"); ok {
		log.Info().Msg("Reading CLOUDFLARE_RECORD_ZONEID from environment variables")
		c.Providers.CLOUDFLAREAccount.Record.ZoneID = os.Getenv("CLOUDFLARE_RECORD_ZONEID")
	}

	// CLOUDFLARE_RECORD_TTL
	if _, ok = os.LookupEnv("CLOUDFLARE_RECORD_TTL"); ok {
		log.Info().Msg("Reading CLOUDFLARE_RECORD_TTL from environment variables")
		ttl, err := strconv.Atoi(os.Getenv("CLOUDFLARE_RECORD_TTL"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("CLOUDFLARE_RECORD_TTL")).Msg("Failed to convert CLOUDFLARE_RECORD_TTL to int")
		} else {
			c.Providers.CLOUDFLAREAccount.Record.TTL = ttl
		}
	}

	// LOG_LEVEL
	if _, ok = os.LookupEnv("LOG_LEVEL"); ok {
		log.Info().Msg("Reading LOG_LEVEL from environment variables")
		c.Log.Level = os.Getenv("LOG_LEVEL")
	}

	// LOG_HUMANIZE
	if _, ok = os.LookupEnv("LOG_HUMANIZE"); ok {
		log.Info().Msg("Reading LOG_HUMANIZE from environment variables")
		// convert string to bool
		humanize, err := strconv.ParseBool(os.Getenv("LOG_HUMANIZE"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("LOG_HUMANIZE")).Msg("Failed to convert LOG_HUMANIZE to bool")
		} else {
			c.Log.Humanize = humanize
		}
	}

	// METRICS_ENABLE
	if _, ok = os.LookupEnv("METRICS_ENABLE"); ok {
		log.Info().Msg("Reading METRICS_ENABLE from environment variables")
		// convert string to bool
		enable, err := strconv.ParseBool(os.Getenv("METRICS_ENABLE"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("METRICS_ENABLE")).Msg("Failed to convert METRICS_ENABLE to bool")
		} else {
			c.Metrics.Enable = enable
		}
	}

	// METRICS ADDRESS
	if _, ok = os.LookupEnv("METRICS_HOST"); ok {
		log.Info().Msg("Reading METRICS_ADDRESS from environment variables")
		c.Metrics.Host = os.Getenv("METRICS_HOST")
	}

	// METRICS_PORT
	if _, ok = os.LookupEnv("METRICS_PORT"); ok {
		log.Info().Msg("Reading METRICS_PORT from environment variables")
		port, err := strconv.Atoi(os.Getenv("METRICS_PORT"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("METRICS_PORT")).Msg("Failed to convert METRICS_PORT to int")
		} else {
			c.Metrics.Port = port
		}
	}

	// METRICS_PATH
	if _, ok = os.LookupEnv("METRICS_PATH"); ok {
		log.Info().Msg("Reading METRICS_PATH from environment variables")
		c.Metrics.Path = os.Getenv("METRICS_PATH")
	}

	// METRICS_LOGGING
	if _, ok = os.LookupEnv("METRICS_LOGGING"); ok {
		log.Info().Msg("Reading METRICS_LOGGING from environment variables")
		// convert string to bool
		logging, err := strconv.ParseBool(os.Getenv("METRICS_LOGGING"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("METRICS_LOGGING")).Msg("Failed to convert METRICS_LOGGING to bool")
		} else {
			c.Metrics.Logging = logging
		}
	}

}
