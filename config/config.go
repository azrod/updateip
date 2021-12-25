package config

import (
	"os"
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
	Metrics    CFGMetrics `yaml:"metrics"`
	AWSAccount struct {
		Enable bool               `yaml:"enable"`
		Secret uip_aws.PawsSecret `yaml:"secret"`
		Record uip_aws.PawsRecord `yaml:"record"`
	} `yaml:"aws_account"`
	OVHAccount struct {
		Enable bool               `yaml:"enable"`
		Secret uip_ovh.PovhSecret `yaml:"secret"`
		Record uip_ovh.PovhRecord `yaml:"record"`
	} `yaml:"ovh_account"`
	CLOUDFLAREAccount struct {
		Enable bool                             `yaml:"enable"`
		Secret uip_cloudflare.PCloudflareSecret `yaml:"secret"`
		Record uip_cloudflare.PCloudflareRecord `yaml:"record"`
	} `yaml:"cloudflare_account"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config CFG) {

	dir := os.Getenv("PATH_CONFIG_DIRECTORY")
	file := os.Getenv("PATH_CONFIG_FILE")

	if dir == "" {
		dir = "/config"
	} else {
		log.Info().Msgf("Using config directory: %s", dir)
	}

	if file == "" {
		file = "config.yaml"
	} else {
		log.Info().Msgf("Using config file: %s", file)
	}

	// open yaml file
	f, err := os.Open(dir + "/" + file)
	if err != nil {
		log.Error().Err(err).Str("path", dir+"/"+file).Msg("Failed to open config file")
	} else {
		defer f.Close()

		// read yaml file
		err = yaml.NewDecoder(f).Decode(&config)
		if err != nil {
			log.Error().Err(err).Str("path", dir+"/"+file).Msg("Failed to decode config file")
		}
	}

	var ok bool

	// try to read from environment variables
	if _, ok = os.LookupEnv("AWS_ACCESS_KEY_ID"); ok {
		log.Info().Msg("Reading AWS_ACCESS_KEY_ID from environment variables")
		config.AWSAccount.Secret.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}

	if _, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY"); ok {
		log.Info().Msg("Reading AWS_SECRET_ACCESS_KEY from environment variables")
		config.AWSAccount.Secret.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	if _, ok = os.LookupEnv("AWS_REGION"); ok {
		log.Info().Msg("Reading AWS_REGION from environment variables")
		config.AWSAccount.Secret.Region = os.Getenv("AWS_REGION")
	}

	if _, ok = os.LookupEnv("AWS_HOSTED_ZONE_ID"); ok {
		log.Info().Msg("Reading AWS_HOSTED_ZONE_ID from environment variables")
		config.AWSAccount.Record.HostedZoneID = os.Getenv("AWS_HOSTED_ZONE_ID")
	}

	if _, ok = os.LookupEnv("AWS_RECORD_NAME"); ok {
		log.Info().Msg("Reading AWS_RECORD_NAME from environment variables")
		config.AWSAccount.Record.Name = os.Getenv("AWS_HOSTED_ZONE_NAME")
	}

	// AWS_RECORD_TTL
	if _, ok = os.LookupEnv("AWS_RECORD_TTL"); ok {
		log.Info().Msg("Reading AWS_RECORD_TTL from environment variables")
		ttl, err := strconv.Atoi(os.Getenv("AWS_RECORD_TTL"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("AWS_RECORD_TTL")).Msg("Failed to convert AWS_RECORD_TTL to int")
		} else {
			config.AWSAccount.Record.TTL = ttl
		}
	}

	// AWS_RECORD_DOMAIN
	if _, ok = os.LookupEnv("AWS_RECORD_DOMAIN"); ok {
		log.Info().Msg("Reading AWS_RECORD_DOMAIN from environment variables")
		config.AWSAccount.Record.Domain = os.Getenv("AWS_RECORD_DOMAIN")
	}

	// AWS_RECORD_COMMENT
	if _, ok = os.LookupEnv("AWS_RECORD_COMMENT"); ok {
		log.Info().Msg("Reading AWS_RECORD_COMMENT from environment variables")
		config.AWSAccount.Record.Comment = os.Getenv("AWS_RECORD_COMMENT")
	}

	// OVH_APPLICATION_KEY
	if _, ok = os.LookupEnv("OVH_APPLICATION_KEY"); ok {
		log.Info().Msg("Reading OVH_APPLICATION_KEY from environment variables")
		config.OVHAccount.Secret.ApplicationKey = os.Getenv("OVH_APPLICATION_KEY")
	}

	// OVH_APPLICATION_SECRET
	if _, ok = os.LookupEnv("OVH_APPLICATION_SECRET"); ok {
		log.Info().Msg("Reading OVH_APPLICATION_SECRET from environment variables")
		config.OVHAccount.Secret.ApplicationSecret = os.Getenv("OVH_APPLICATION_SECRET")
	}

	// OVH_CONSUMER_KEY
	if _, ok = os.LookupEnv("OVH_CONSUMER_KEY"); ok {
		log.Info().Msg("Reading OVH_CONSUMER_KEY from environment variables")
		config.OVHAccount.Secret.ConsumerKey = os.Getenv("OVH_CONSUMER_KEY")
	}

	// OVH_RECORD_NAME
	if _, ok = os.LookupEnv("OVH_RECORD_NAME"); ok {
		log.Info().Msg("Reading OVH_RECORD_NAME from environment variables")
		config.OVHAccount.Record.Name = os.Getenv("OVH_RECORD_NAME")
	}

	// OVH_RECORD_TTL
	if _, ok = os.LookupEnv("OVH_RECORD_TTL"); ok {
		log.Info().Msg("Reading OVH_RECORD_TTL from environment variables")
		ttl, err := strconv.Atoi(os.Getenv("OVH_RECORD_TTL"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("OVH_RECORD_TTL")).Msg("Failed to convert OVH_RECORD_TTL to int")
		} else {
			config.OVHAccount.Record.TTL = ttl
		}
	}

	// OVH_RECORD_ZONE
	if _, ok = os.LookupEnv("OVH_RECORD_ZONE"); ok {
		log.Info().Msg("Reading OVH_RECORD_ZONE from environment variables")
		config.OVHAccount.Record.Zone = os.Getenv("OVH_RECORD_ZONE")
	}

	// CLOUDFLARE_API_KEY
	if _, ok = os.LookupEnv("CLOUDFLARE_API_KEY"); ok {
		log.Info().Msg("Reading CLOUDFLARE_API_KEY from environment variables")
		config.CLOUDFLAREAccount.Secret.APIKey = os.Getenv("CLOUDFLARE_API_KEY")
	}

	// CLOUDFLARE_EMAIL
	if _, ok = os.LookupEnv("CLOUDFLARE_EMAIL"); ok {
		log.Info().Msg("Reading CLOUDFLARE_EMAIL from environment variables")
		config.CLOUDFLAREAccount.Secret.Email = os.Getenv("CLOUDFLARE_EMAIL")
	}

	// CLOUDFLARE_RECORD_NAME
	if _, ok = os.LookupEnv("CLOUDFLARE_RECORD_NAME"); ok {
		log.Info().Msg("Reading CLOUDFLARE_RECORD_NAME from environment variables")
		config.CLOUDFLAREAccount.Record.Name = os.Getenv("CLOUDFLARE_RECORD_NAME")
	}

	// CLOUDFLARE_RECORD_TTL
	if _, ok = os.LookupEnv("CLOUDFLARE_RECORD_TTL"); ok {
		log.Info().Msg("Reading CLOUDFLARE_RECORD_TTL from environment variables")
		ttl, err := strconv.Atoi(os.Getenv("CLOUDFLARE_RECORD_TTL"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("CLOUDFLARE_RECORD_TTL")).Msg("Failed to convert CLOUDFLARE_RECORD_TTL to int")
		} else {
			config.CLOUDFLAREAccount.Record.TTL = ttl
		}
	}

	// LOG_LEVEL
	if _, ok = os.LookupEnv("LOG_LEVEL"); ok {
		log.Info().Msg("Reading LOG_LEVEL from environment variables")
		config.Log.Level = os.Getenv("LOG_LEVEL")
	}

	// LOG_HUMANIZE
	if _, ok = os.LookupEnv("LOG_HUMANIZE"); ok {
		log.Info().Msg("Reading LOG_HUMANIZE from environment variables")
		// convert string to bool
		humanize, err := strconv.ParseBool(os.Getenv("LOG_HUMANIZE"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("LOG_HUMANIZE")).Msg("Failed to convert LOG_HUMANIZE to bool")
		} else {
			config.Log.Humanize = humanize
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
			config.Metrics.Enable = enable
		}
	}

	// METRICS_PORT
	if _, ok = os.LookupEnv("METRICS_PORT"); ok {
		log.Info().Msg("Reading METRICS_PORT from environment variables")
		port, err := strconv.Atoi(os.Getenv("METRICS_PORT"))
		if err != nil {
			log.Error().Err(err).Str("value", os.Getenv("METRICS_PORT")).Msg("Failed to convert METRICS_PORT to int")
		} else {
			config.Metrics.Port = port
		}
	}

	// METRICS_PATH
	if _, ok = os.LookupEnv("METRICS_PATH"); ok {
		log.Info().Msg("Reading METRICS_PATH from environment variables")
		config.Metrics.Path = os.Getenv("METRICS_PATH")
	}

	return config

}
