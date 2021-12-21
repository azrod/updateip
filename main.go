package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/azrod/updateip/config"
	uip_aws "github.com/azrod/updateip/pkg/providers/aws"
	uip_ovh "github.com/azrod/updateip/pkg/providers/ovh"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Load config
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if c.Log.Humanize {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	log.Info().Msg("Starting UpdateIP")

	// Parse loglevel
	if l, err := zerolog.ParseLevel(c.Log.Level); err != nil {
		log.Fatal().Err(err).Msg("cannot parse log level")
	} else {
		zerolog.SetGlobalLevel(l)
		log.Info().Msgf("Log level set to %s", l.String())
	}

	if c.AWSAccount.Enable {

		Paws := uip_aws.Paws{
			Record: c.AWSAccount.Record,
			Secret: c.AWSAccount.Secret,
		}

		if err := Paws.NewClient(); err != nil {
			log.Fatal().Err(err).Interface("config", Paws).Msg("Failed to setup AWS client")
		}

		if err := Paws.Run(); err != nil {
			log.Error().Err(err).Interface("config", Paws).Msg("Error on module AWS")
		}

	}

	if c.OVHAccount.Enable {

		Povh := uip_ovh.Povh{
			Record: c.OVHAccount.Record,
			Secret: c.OVHAccount.Secret,
		}

		if err := Povh.NewClient(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup OVH client")
		}

		if err := Povh.Run(); err != nil {
			log.Error().Err(err).Msg("Error on module OVH")
		}

	}

	if c.OVHAccount.Enable {

		Povh := uip_ovh.Povh{
			Record: c.OVHAccount.Record,
			Secret: c.OVHAccount.Secret,
		}

		if err := Povh.NewClient(); err != nil {
			log.Fatal().Err(err).Msg("Failed to setup OVH client")
		}

		if err := Povh.Run(); err != nil {
			log.Error().Err(err).Msg("Error on module OVH")
		}

	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

LOOP:
	for {
		select {
		case sig := <-sigs:
			log.Info().Msg(sig.String())
			break LOOP
		}
	}
}
