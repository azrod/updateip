package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/azrod/updateip/pkg/config"
	"github.com/azrod/updateip/pkg/metrics"
	uip_aws "github.com/azrod/updateip/pkg/providers/aws"
	uip_cloudflare "github.com/azrod/updateip/pkg/providers/cloudflare"
	uip_ovh "github.com/azrod/updateip/pkg/providers/ovh"
)

var (
	m           *metrics.Metrics
	Paws        uip_aws.Paws
	Povh        uip_ovh.Povh
	PCloudflare uip_cloudflare.PCloudflare
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Load config
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	if c.Log.Humanize {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	// Parse loglevel
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if l, err := zerolog.ParseLevel(c.Log.Level); err != nil {
		log.Fatal().Err(err).Msg("cannot parse log level")
	} else {
		zerolog.SetGlobalLevel(l)
		log.Info().Msgf("Log level set to %s", l.String())
	}

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
