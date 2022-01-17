package uip_cloudflare

import (
	"time"

	"github.com/cloudflare/cloudflare-go"
)

type PCloudflare struct {
	Secret  PCloudflareSecret
	Record  PCloudflareRecord
	clients struct {
		cloudflare *cloudflare.API
	}
	Events chan string
	Loop   time.Ticker
}

type PCloudflareSecret struct {
	//APIkey
	APIKey string `yaml:"api_key" env:"CLOUDFLARE_API_KEY"`

	//Email
	Email string `yaml:"email" env:"CLOUDFLARE_EMAIL"`
}

// PCloudflareRecord
type PCloudflareRecord struct {
	// The DNS record to update
	Name string `yaml:"name" env:"CLOUDFLARE_RECORD_NAME"`

	// TTL of the record
	TTL int `yaml:"ttl" env:"CLOUDFLARE_RECORD_TTL"`

	//Domain name of the record
	Domain string `yaml:"domain" env:"CLOUDFLARE_RECORD_DOMAIN"`

	//Zone ID of the record
	ZoneID string `yaml:"zone_id,omitempty" env:"CLOUDFLARE_ZONE_ID"`
}

type record struct {
	Expire       time.Time
	LastValue    string
	LastChangeID string
}
