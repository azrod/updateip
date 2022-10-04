package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/azrod/updateip/pkg/config"
	"github.com/azrod/updateip/pkg/metrics"
	uip_aws "github.com/azrod/updateip/pkg/providers/aws"
	uip_cloudflare "github.com/azrod/updateip/pkg/providers/cloudflare"
	uip_ovh "github.com/azrod/updateip/pkg/providers/ovh"
	"github.com/azrod/zr"
	"github.com/azrod/zr/pkg/format"
	hr "github.com/azrod/zr/pkg/hotreload"
	"github.com/azrod/zr/pkg/level"
)

var (
	m           *metrics.Metrics
	Paws        uip_aws.Paws
	Povh        uip_ovh.Povh
	PCloudflare uip_cloudflare.PCloudflare
)

func main() {

	// Load config
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	logLevel, err := level.ParseLogLevel(c.Log.Level)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse log level")
	}

	var logFormat format.LogFormat
	if c.Log.Humanize {
		logFormat = format.LogFormatHuman
	} else {
		logFormat = format.LogFormatJson
	}

	zr.Setup(
		zr.Level(logLevel),
		zr.Format(logFormat),
		zr.WithCustomHotReload(
			hr.WithNoHotReload(),
		),
	)

	log.Info().Msg("Starting UpdateIP")

	if c.Providers.AWSAccount.Enable {

		Paws = uip_aws.Paws{
			Record: c.Providers.AWSAccount.Record,
			Secret: c.Providers.AWSAccount.Secret,
		}

		if err := Paws.NewClient(); err != nil {
			log.Fatal().Err(err).Interface("config", Paws).Msg("Failed to setup AWS client")
		}

		go func() {
			if err := Paws.Run(); err != nil {
				log.Error().Err(err).Interface("config", Paws).Msg("Error on module AWS")
			}
		}()

	}

	if c.Providers.OVHAccount.Enable {

		Povh = uip_ovh.Povh{
			Record: c.Providers.OVHAccount.Record,
			Secret: c.Providers.OVHAccount.Secret,
		}

		if err := Povh.NewClient(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup OVH client")
		}
		go func() {
			if err := Povh.Run(); err != nil {
				log.Error().Err(err).Msg("Error on module OVH")
			}
		}()

	}

	if c.Providers.CLOUDFLAREAccount.Enable {

		PCloudflare = uip_cloudflare.PCloudflare{
			Record: c.Providers.CLOUDFLAREAccount.Record,
			Secret: c.Providers.CLOUDFLAREAccount.Secret,
		}

		if err := PCloudflare.NewClient(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup Cloudflare client")
		}

		go func() {
			if err := PCloudflare.Run(); err != nil {
				log.Error().Err(err).Msg("Error on module Cloudflare")
			}
		}()

	}

	if c.Metrics.Enable {
		log.Info().Msg("Starting Metrics Server")
		m = metrics.Init(c.Metrics)

		if c.Providers.AWSAccount.Enable {
			m.RegisterPkg(Paws.RegistryMetrics())
		}

		if c.Providers.CLOUDFLAREAccount.Enable {
			m.RegisterPkg(PCloudflare.RegistryMetrics())
		}

		if c.Providers.OVHAccount.Enable {
			m.RegisterPkg(PCloudflare.RegistryMetrics())
		}

		m.Run()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	func() {
		<-sigs
		log.Info().Msg("Shutting down")
		os.Exit(1)
	}()

}
