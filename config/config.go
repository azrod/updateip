package config

import (
	"os"
	"strconv"

	uip_aws "github.com/azrod/updateip/pkg/providers/aws"
	uip_cloudflare "github.com/azrod/updateip/pkg/providers/cloudflare"
	uip_ovh "github.com/azrod/updateip/pkg/providers/ovh"

	"gopkg.in/yaml.v2"
)

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
	Metrics struct {
		// Prefix is the prefix used for all metrics.
		Prefix string `yaml:"prefix"`
		// Enable is a boolean that determines if metrics are enabled.
		Enable bool `yaml:"enable"`
		// Host is the endpoint used for the metrics server.
		Host string `yaml:"host"`
		// Port is the port used for the metrics server.
		Port int `yaml:"port"`
	}
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
func LoadConfig() (config CFG, err error) {

	// open yaml file
	f, err := os.Open("config.yaml")

	// if file not found, try to read from environment variables
	if err != nil {

		// AWS Account
		config.AWSAccount.Secret.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")         //Access key ID
		config.AWSAccount.Secret.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY") //Secret access key
		config.AWSAccount.Secret.Region = os.Getenv("AWS_REGION")                     // Region
		config.AWSAccount.Record.Name = os.Getenv("AWS_RECORD_NAME")                  // The DNS record to update
		config.AWSAccount.Record.TTL, err = strconv.Atoi(os.Getenv("AWS_RECORD_TTL")) // TTL of the record
		if err != nil {
			return config, err
		}
		config.AWSAccount.Record.Domain = os.Getenv("AWS_RECORD_DOMAIN")               //Domain name of the record
		config.AWSAccount.Record.HostedZoneID = os.Getenv("AWS_RECORD_HOSTED_ZONE_ID") // optional
		config.AWSAccount.Record.Comment = os.Getenv("AWS_RECORD_COMMENT")             // optional

		// OVH ACCOUNT
		config.OVHAccount.Secret.ApplicationKey = os.Getenv("OVH_APPLICATION_KEY")       // Application key
		config.OVHAccount.Secret.ApplicationSecret = os.Getenv("OVH_APPLICATION_SECRET") // Application secret
		config.OVHAccount.Secret.ConsumerKey = os.Getenv("OVH_CONSUMER_KEY")             // Consumer key
		config.OVHAccount.Record.Name = os.Getenv("OVH_RECORD_NAME")                     // The DNS record to update
		config.OVHAccount.Record.TTL, err = strconv.Atoi(os.Getenv("OVH_RECORD_TTL"))    // TTL of the record
		if err != nil {
			return config, err
		}
		config.OVHAccount.Record.Zone = os.Getenv("OVH_RECORD_ZONE") //Zone of the record

		config.Log.Level = os.Getenv("LOG_LEVEL")                               // trace, debug, info, warn, error, fatal, panic
		config.Log.Humanize, err = strconv.ParseBool(os.Getenv("LOG_HUMANIZE")) // humanize log output
		if err != nil {
			return config, err
		}
		config.Metrics.Enable, err = strconv.ParseBool(os.Getenv("METRICS_ENABLE")) // enable metrics
		if err != nil {
			return config, err
		}
		config.Metrics.Host = os.Getenv("METRICS_HOST")                    // metrics host
		config.Metrics.Port, err = strconv.Atoi(os.Getenv("METRICS_PORT")) // metrics port
		if err != nil {
			return config, err
		}
		config.Metrics.Prefix = os.Getenv("METRICS_PREFIX") // metrics prefix

		return config, nil
	}

	// read yaml file
	err = yaml.NewDecoder(f).Decode(&config)
	return config, err

}
