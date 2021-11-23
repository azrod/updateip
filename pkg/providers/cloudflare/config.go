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
	APIKey string `yaml:"api_key"`

	//Email
	Email string `yaml:"email"`
}

// PCloudflareRecord
type PCloudflareRecord struct {
	// The DNS record to update
	Name string `yaml:"name"`

	// TTL of the record
	TTL int `yaml:"ttl"`

	//Domain name of the record
	Domain string `yaml:"domain"`

	//Zone ID of the record
	ZoneID string `yaml:"zone_id,omitempty"`
}

type record struct {
	Expire       time.Time
	LastValue    string
	LastChangeID string
}
