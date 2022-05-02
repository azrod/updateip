package uip_aws

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/azrod/updateip/pkg/ip"
	"github.com/rs/zerolog/log"
)

var (
	rec = record{
		Expire:    time.Now(),
		LastValue: "",
	}
)

func (d *Paws) NewClient() error {
	defer timeTrackS(time.Now(), "aws_NewClient")

	var err error
	d.clients.aws, err = session.NewSession(&aws.Config{
		Region:      aws.String(d.Secret.Region),
		Credentials: credentials.NewStaticCredentials(d.Secret.AccessKeyID, d.Secret.SecretAccessKey, ""),
	})

	d.clients.route53 = route53.New(d.clients.aws)

	if d.Record.HostedZoneID == "" {
		if d.Record.HostedZoneID, err = d.getHostedZoneID(); err != nil {
			return err
		}
	}

	d.Events = make(chan string, 100)
	d.Loop = *time.NewTicker(60 * time.Second)
	return err
}

func (d *Paws) UpdateRecord(ip net.IP) error {
	defer timeTrackS(time.Now(), "aws_UpdateRecord")

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(d.Record.Name),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip.String()),
							},
						},
						TTL:  aws.Int64(int64(d.Record.TTL)),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String(d.Record.Comment),
		},
		HostedZoneId: aws.String(d.Record.HostedZoneID), // Required
	}

	resp, err := d.clients.route53.ChangeResourceRecordSets(input)
	if err != nil {
		return err
	}

	if *resp.ChangeInfo.Id != "" {
		rec.LastChangeID = *resp.ChangeInfo.Id
		d.Events <- "Update Record (ChangeID : " + rec.LastChangeID + ") : IN PROGRESS"
	} else {
		return errors.New("no ChangeID")
	}

	return nil
}

func (d *Paws) getHostedZoneID() (HostedZoneID string, err error) {
	defer timeTrackS(time.Now(), "aws_getHostedZoneID")

	listParams := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(d.Record.Domain), // Required
	}
	req, resp := d.clients.route53.ListHostedZonesByNameRequest(listParams)
	err = req.Send()
	if err != nil {
		return "", err
	}

	if len(resp.HostedZones) != 0 {
		HostedZoneID = *resp.HostedZones[0].Id
		// remove the /hostedzone/ path if it's there
		if strings.HasPrefix(HostedZoneID, "/hostedzone/") {
			HostedZoneID = strings.TrimPrefix(HostedZoneID, "/hostedzone/")
		}
		return HostedZoneID, nil
	} else {
		return "", fmt.Errorf("failed to get hostedzoneID")
	}

}

func (d *Paws) GetRecord() (record string, err error) {
	defer timeTrackS(time.Now(), "aws_GetRecord")

	if time.Now().After(rec.Expire) || rec.LastValue == "" {

		log.Trace().Str("LastValue", rec.LastValue).Time("Now", time.Now()).Time("Expire", rec.Expire).Msg("GetRecord is expired or empty")

		var resp *route53.ListResourceRecordSetsOutput

		listParams := &route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(d.Record.HostedZoneID), // Required
			StartRecordName: aws.String(d.Record.Name),
		}
		if resp, err = d.clients.route53.ListResourceRecordSets(listParams); err != nil {
			return "", err
		}

		rec.Expire = time.Now().Add(10 * time.Minute)
		rec.LastValue = *resp.ResourceRecordSets[0].ResourceRecords[0].Value

	} else {
		log.Trace().Str("LastValue", rec.LastValue).Time("Now", time.Now()).Time("Expire", rec.Expire).Msg("GetRecord is not expired")
	}

	return rec.LastValue, nil

}

// get change status route53
func (d *Paws) GetChangeStatus() (terminated bool, err error) {
	defer timeTrackS(time.Now(), "aws_GetChangeStatus")

	if rec.LastChangeID != "" {

		input := &route53.GetChangeInput{
			Id: aws.String(rec.LastChangeID), // Required
		}
		req, resp := d.clients.route53.GetChangeRequest(input)
		err = req.Send()
		if err != nil {
			return false, err
		}

		if *resp.ChangeInfo.Status == "INSYNC" {
			d.Events <- "Update Record (ChangeID : " + rec.LastChangeID + ") : Done"
			rec.LastChangeID = ""
			return true, nil
		} else if *resp.ChangeInfo.Status == "PENDING" {
			d.Events <- "Update Record (ChangeID : " + rec.LastChangeID + ") : Pending"
			return false, nil
		} else {
			return false, nil
		}
	}

	return true, nil

}
func (d *Paws) Run() error {

	log.Info().Msg("Starting AWS Route53 Module")

	providerStatus.Set(1)

	for {
		select {
		case e := <-d.Events:
			eventReceive.Inc()
			log.Info().Msgf("Event => %s", e)
		case <-d.Loop.C:

			if ok, err := d.GetChangeStatus(); ok && err == nil {
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
					log.Info().Str("Record", d.Record.Name).Str("DNSIP", r).Str("ActualIP", i.String()).Msg("New IP address detected. Update")
					countUpdate.Inc()
					if err = d.UpdateRecord(i); err != nil {
						log.Error().Err(err).Msg("Failed to update dns record")
					}
				}

			} else if err != nil {
				log.Error().Err(err).Msg("Failed to get change status")
			}

		}
	}
}
