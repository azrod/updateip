package uip_gandi

import (
	"fmt"
	"net"
	"time"

	"github.com/go-gandi/go-gandi"
	"github.com/go-gandi/go-gandi/livedns"
)

func (d *Pgandi) NewClient() error {

	if d.Secret.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if d.Secret.SharingID == "" {
		return fmt.Errorf("SharingID is required")
	}

	d.clients.dns = gandi.NewLiveDNSClient(d.Secret.APIKey, gandi.Config{SharingID: d.Secret.SharingID, Debug: true, DryRun: false})

	d.Events = make(chan string, 100)
	d.Loop = *time.NewTicker(5 * time.Second)

	return nil

}

func (d *Pgandi) UpdateRecord(ip net.IP) error {

	if d.clients.dns == nil {
		return fmt.Errorf("client is not initialized")
	}

	if d.Record.Domain == "" {
		return fmt.Errorf("domain is required")
	}

	if d.Record.Name == "" {
		return fmt.Errorf("RecordName is required")
	}

	if d.Record.TTL == 0 {
		return fmt.Errorf("TTL is required")
	}

	if ip == nil {
		return fmt.Errorf("IP is required")
	}

	// Get current IP
	currentIP, err := d.clients.dns.GetDomainRecordsByName(d.Record.Domain, d.Record.Name)
	if err != nil {
		return err
	}

	if currentIP[0].RrsetValues[0] != ip.String() {

		a := livedns.DomainRecord{
			RrsetName:   d.Record.Name,
			RrsetTTL:    d.Record.TTL,
			RrsetValues: []string{ip.String()},
		}

		recs := []livedns.DomainRecord{a}

		// Update record
		resp, err := d.clients.dns.UpdateDomainRecordsByName(d.Record.Domain, d.Record.Name, recs)
		if resp.Status != "success" || err != nil {
			d.Events <- fmt.Sprintf("Error updating record: %s", resp.Message)
			return fmt.Errorf("error updating record: %s", resp.Message)
		}

		d.Events <- fmt.Sprintf("Updated IP: %s", ip.String())

	} else {

		d.Events <- fmt.Sprintf("IP is already up to date: %s", ip.String())

	}

	return nil

}
