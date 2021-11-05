package uip_ovh

import (
	"time"

	"github.com/ovh/go-ovh/ovh"
)

type Povh struct {
	Secret  PovhSecret
	Record  PovhRecord
	clients struct {
		ovh *ovh.Client
	}
	Events chan string
	Loop   time.Ticker
}

type PovhSecret struct {
	// Application key to access the API
	ApplicationKey string `yaml:"application_key"`
	// Application Secret to access the API
	ApplicationSecret string `yaml:"application_secret"`
	// Region
	Region string `yaml:"region"`
	//Consumer key
	ConsumerKey string `yaml:"consumer_key"`
}

type PovhRecord struct {
	// The DNS record to update
	Name string `yaml:"name"`
	// TTL of the record
	TTL int `yaml:"ttl"`
	//Zone of the record
	Zone string `yaml:"zone"`
}

// txtRecordRequest represents the request body to DO's API to make a TXT record
type txtRecordRequest struct {
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
}

// txtRecordResponse represents a response from DO's API after making a TXT record
type txtRecordResponse struct {
	ID        int    `json:"id"`
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl"`
	Zone      string `json:"zone"`
}

type StatusResponse struct {
	Errors     []string `json:"errors"`     // Error list
	IsDeployed bool     `json:"isDeployed"` // True if the zone has successfully been deployed
	Warnings   []string `json:"warnings"`   // Warning list
}
