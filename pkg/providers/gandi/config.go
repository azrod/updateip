package uip_gandi

import (
	"time"

	"github.com/go-gandi/go-gandi/livedns"
)

type Pgandi struct {
	Secret  PgandiSecret
	Record  PgandiRecord
	clients struct {
		dns *livedns.LiveDNS
	}
	Events chan string
	Loop   time.Ticker
}

type PgandiSecret struct {
	//SharingID is the Organization ID, available from the Organization API
	SharingID string `yaml:"sharing_id"`

	//APIKey is the API Key, available from https://account.gandi.net/en/
	APIKey string `yaml:"api_key"`
}

type PgandiRecord struct {
	// The DNS record to update
	Name string `yaml:"name"`
	// TTL of the record
	TTL int `yaml:"ttl"`
	//Domain name of the record
	Domain string `yaml:"domain"`
}

type record struct {
	Expire       time.Time
	LastValue    string
	LastChangeID string
}
