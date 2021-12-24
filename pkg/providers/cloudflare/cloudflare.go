package uip_cloudflare

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/azrod/updateip/pkg/ip"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

var (
	rec = record{
		Expire:    time.Now(),
		LastValue: "",
	}
)

// Setup new cloudflare client
func (d *PCloudflare) NewClient() error {
	var err error

	log.Trace().Msg("Creating new Cloudflare client")
	d.clients.cloudflare, err = cloudflare.New(d.Secret.APIKey, d.Secret.Email)
	if err != nil {
		return err
	}

	if d.Record.ZoneID == "" {
		log.Trace().Msg("Getting zone ID")
		if d.Record.ZoneID, err = d.clients.cloudflare.ZoneIDByName(d.Record.Domain); err != nil {
			return err
		} else {
			log.Trace().Msgf("Zone ID: %s", d.Record.ZoneID)
		}
	}

	d.Events = make(chan string, 100)
	d.Loop = *time.NewTicker(60 * time.Second)
	return err
}

// Update record with new IP
func (d *PCloudflare) UpdateRecord(ip net.IP) error {
	var err error

	ctx := context.Background()

	log.Trace().Msg("Updating record")

	recs, err := d.clients.cloudflare.DNSRecords(ctx, d.Record.ZoneID, cloudflare.DNSRecord{})
	if err != nil {
		return err
	}

	for _, rec := range recs {
		if rec.Name == d.Record.Name {
			log.Trace().Msgf("Updating record %s", rec.Name)
			rec.Content = ip.String()
			if err = d.clients.cloudflare.UpdateDNSRecord(ctx, d.Record.ZoneID, rec.ID, rec); err != nil {

				d.Events <- fmt.Sprintf("Error updating record: %s", err.Error())
				return err
			} else {
				d.Events <- fmt.Sprintf("Updated IP: %s", ip.String())

				d.Loop.Stop()
				d.loopCheckApply(ip.String())
				d.Loop.Reset(60 * time.Second)

			}
			break
		}
	}

	return err
}

// loop check if record is apply

func (d *PCloudflare) loopCheckApply(ip string) bool {

	c := 0

	for {
		c++
		r, err := d.GetRecord()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get record")
			break
		}

		if r == ip {
			d.Events <- "Record is successfully updated"
			return true
		} else {
			d.Events <- fmt.Sprintf("Waiting for record to be updated. Try: %v/10", c)

		}
		time.Sleep(15 * time.Second)

		if c > 10 {
			d.Events <- "Record is not updated. Waiting for next check"
			return false
		}
	}

	return false

}

// Get record
func (d *PCloudflare) GetRecord() (record string, err error) {

	log.Trace().Msg("Getting record")

	if rec.Expire.After(time.Now()) || rec.LastValue == "" {

		ctx := context.Background()

		recs, err := d.clients.cloudflare.DNSRecords(ctx, d.Record.ZoneID, cloudflare.DNSRecord{})
		if err != nil {
			return "", err
		}

		for _, record := range recs {
			if record.Name == d.Record.Name {
				rec.LastValue = record.Content
			}
		}

		// TODO Expire time

		rec.Expire = time.Now().Add(10 * time.Minute)

	}

	return rec.LastValue, nil

}

func (d *PCloudflare) Run() error {

	log.Info().Msg("Starting Cloudflare Module")

	for {
		select {
		case e := <-d.Events:
			log.Info().Msgf("Event: => %s", e)
		case <-d.Loop.C:

			r, err := d.GetRecord()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get record")
				break
			}

			i, err := ip.GetMyExternalIP()
			if err != nil {
				log.Error().Err(err).Msg("Could not get External IP")
				break
			}

			if r != i.String() {
				// go lock()
				log.Info().Str("DNSIP", r).Str("ActualIP", i.String()).Msg("New IP address detected. Update")

				if err = d.UpdateRecord(i); err != nil {
					log.Error().Err(err).Msg("Failed to update dns record")
				}
			}
		}
	}
}
