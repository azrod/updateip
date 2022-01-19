package uip_aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type Paws struct {
	Secret  PawsSecret
	Record  PawsRecord
	clients struct {
		aws     *session.Session
		route53 *route53.Route53
	}
	Events chan string
	Loop   time.Ticker
}

type PawsSecret struct {
	//Access key ID
	AccessKeyID string `yaml:"access_key_id" env:"AWS_ACCESS_KEY_ID"`
	//Secret access key
	SecretAccessKey string `yaml:"secret_access_key" env:"AWS_SECRET_ACCESS_KEY"`
	// Region
	Region string `yaml:"region" env:"AWS_REGION"`
}

type PawsRecord struct {
	// The DNS record to update
	Name string `yaml:"name" env:"AWS_RECORD_NAME"`
	// TTL of the record
	TTL int `yaml:"ttl" env:"AWS_RECORD_TTL"`
	//Domain name of the record
	Domain string `yaml:"domain" env:"AWS_RECORD_DOMAIN"`
	//Hosted zone ID of the record
	HostedZoneID string `yaml:"hosted_zone_id,omitempty" env:"AWS_HOSTED_ZONE_ID"`
	//Comment for updated record
	Comment string `yaml:"comment" env:"AWS_RECORD_COMMENT"`
}

type record struct {
	Expire       time.Time
	LastValue    string
	LastChangeID string
}
