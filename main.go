package main

import (
	"context"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/datasource"
	"github.com/aliarbak/ethereum-connector/destination"
	_ "github.com/aliarbak/ethereum-connector/docs"
	"github.com/aliarbak/ethereum-connector/evm"
	"github.com/aliarbak/ethereum-connector/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"strings"
)

func main() {
	ctx := context.Background()
	config := configs.ReadConfigs()
	configureLogger(config.Logging)

	dataSourceFactory, err := datasource.NewFactory(ctx, config.DataSource)
	if err != nil {
		log.Fatal().Err(err).Msg("datasource factory initialization failed")
	}

	destinationFactory, err := destination.NewFactory(ctx, config.Destination)
	if err != nil {
		log.Fatal().Err(err).Msg("destination factory initialization failed")
	}

	if err = evm.RunListeners(ctx, dataSourceFactory, destinationFactory); err != nil {
		log.Fatal().Err(err).Msg("listeners initialization failed")
	}

	if err = evm.RunContractSyncWorkers(ctx, dataSourceFactory, destinationFactory); err != nil {
		log.Fatal().Err(err).Msg("sync workers initialization failed")
	}

	if err = http.Run(ctx, config.Http, dataSourceFactory, destinationFactory); err != nil {
		log.Fatal().Err(err).Msg("http initialization failed")
	}
}

func configureLogger(config configs.LoggingConfig) {
	logLevel, err := zerolog.ParseLevel(strings.ToLower(config.Level))
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = config.TimeFieldFormat
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFieldName = config.TimestampFieldName
	zerolog.LevelFieldName = config.LevelFieldName
	zerolog.MessageFieldName = config.MessageFieldName
}
