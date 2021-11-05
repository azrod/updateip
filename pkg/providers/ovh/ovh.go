package uip_ovh

import (
	"fmt"
	"net"
	"time"

	"github.com/azrod/updateip/pkg/ip"
	"github.com/jpillora/go-tld"
	"github.com/ovh/go-ovh/ovh"
	"github.com/rs/zerolog/log"
)

// New Client
func (d *Povh) NewClient() (err error) {
	d.clients.ovh, err = ovh.NewClient(
		d.Secret.Region,
		d.Secret.ApplicationKey,
		d.Secret.ApplicationSecret,
		d.Secret.ConsumerKey,
	)

	u, _ := tld.Parse(d.Record.Name)
	d.Record.Zone = u.Domain
	d.Record.Name = u.Subdomain

	d.Events = make(chan string, 100)
	d.Loop = *time.NewTicker(5 * time.Second)

	return
}

// update record
func (d *Povh) UpdateRecord(recordID int, ip net.IP) (err error) {
	var respData txtRecordResponse
	reqData := txtRecordRequest{FieldType: "A", Target: ip.String(), TTL: 300}

	// update record
	if err = d.clients.ovh.Put(
		"/domain/zone/"+d.Record.Zone+"/record/"+fmt.Sprint(recordID),
		reqData,
		respData,
	); err != nil {
		return
	}

	if err = d.RefreshZoneRecords(); err != nil {
		return
	}

	return
}

// get record ID
func (d *Povh) GetRecordID() (recordID int, err error) {

	var records []int
	err = d.clients.ovh.Get(
		"/domain/zone/"+d.Record.Zone+"/record?fieldType=TXT&subDomain="+d.Record.Name,
		records,
	)

	if err != nil {
		return
	}

	if len(records) == 0 {
		err = fmt.Errorf("no record found")
		return
	}

	for _, r := range records {
		var respData txtRecordResponse
		err = d.clients.ovh.Get(
			"/domain/zone/"+d.Record.Zone+"/record/"+fmt.Sprint(r),
			respData,
		)

		if err != nil {
			return
		}

		if respData.SubDomain == d.Record.Name {
			recordID = r
			return
		}
	}

	return
}

// get record with record ID
func (d *Povh) GetRecord(recordID int) (record string, err error) {

	var respData txtRecordResponse
	err = d.clients.ovh.Get(
		"/domain/zone/"+d.Record.Zone+"/record/"+fmt.Sprint(recordID),
		respData,
	)

	if err != nil {
		return
	}

	record = respData.Target
	return
}

// GetchangeStatus returns the status of the change
func (d *Povh) GetChangeStatus() (status bool, err error) {
	var respData StatusResponse
	err = d.clients.ovh.Get("/domain/zone/"+d.Record.Zone+"/status", respData)
	return respData.IsDeployed, err
}

// Refresh Zone Records
func (d *Povh) RefreshZoneRecords() (err error) {
	var respData []int
	err = d.clients.ovh.Post(
		"/domain/zone/"+d.Record.Zone+"/refresh",
		nil,
		nil,
	)

	if err != nil {
		return
	} else {
		d.Events <- "Zone " + d.Record.Zone + " Refreshed"
	}

	for _, r := range respData {
		var respData txtRecordResponse
		err = d.clients.ovh.Get(
			"/domain/zone/"+d.Record.Zone+"/record/"+fmt.Sprint(r),
			respData,
		)

		if err != nil {
			return
		}

	}

	return
}

// Run
func (d *Povh) Run() (err error) {
	log.Info().Msg("Starting OVH Module")
	for {
		select {
		case e := <-d.Events:
			log.Info().Msgf("Event: %s", e)
		case <-d.Loop.C:

			if ok, err := d.GetChangeStatus(); ok && err == nil {
				recordID, err := d.GetRecordID()
				if err != nil {
					return err
				}

				r, err := d.GetRecord(recordID)
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
					if err = d.UpdateRecord(recordID, i); err != nil {
						log.Error().Err(err).Msg("Failed to update dns record")
					}
				}

			} else if err != nil {
				log.Error().Err(err).Msg("Failed to get change status")
			}

		}
	}
}
