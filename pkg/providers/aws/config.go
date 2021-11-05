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
	AccessKeyID string `yaml:"access_key_id" `
	//Secret access key
	SecretAccessKey string `yaml:"secret_access_key"`
	// Region
	Region string `yaml:"region"`
}

type PawsRecord struct {
	// The DNS record to update
	Name string `yaml:"name"`
	// TTL of the record
	TTL int `yaml:"ttl"`
	//Domain name of the record
	Domain string `yaml:"domain"`
	//Hosted zone ID of the record
	HostedZoneID string `yaml:"hosted_zone_id,omitempty"`
	//Comment for updated record
	Comment string `yaml:"comment"`
}

type record struct {
	Expire       time.Time
	LastValue    string
	LastChangeID string
}
