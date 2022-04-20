package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/matvoy/cparser/internal/client/http"
	"github.com/matvoy/cparser/internal/config"
	"github.com/matvoy/cparser/internal/models"
	"github.com/matvoy/cparser/internal/storage/ch"
	"github.com/matvoy/cparser/internal/storage/pg"
)

type CosmosClient interface {
	GetBlockData(height uint64) (*models.Block, error)
}

type GrpcCosmosClient interface {
	CosmosClient
	Close()
}

type Repository interface {
	InsertBlockData(b *models.Block) error
	SelectTransfers() ([]*models.TransferView, error)
	Close()
	Cleanup() error
}

type App struct {
	pg          Repository
	ch          Repository
	rest        CosmosClient
	grpc        GrpcCosmosClient
	blockNumber uint64
}

func NewApp(cfg *config.Config) *App {
	pg, err := pg.NewRepository(cfg.PostgresDSN)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	ch, err := ch.NewRepository(cfg.ClickHouseDSN1, cfg.ClickHouseDSN2, cfg.ClickHouseDSN3, cfg.ClickHouseDSN4)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	httpCli, err := http.NewClient(cfg.RestURL) // https://api.cosmos.network
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return &App{
		blockNumber: cfg.BlockNumber,
		pg:          pg,
		ch:          ch,
		rest:        httpCli,
		grpc:        nil,
	}
}

func (a *App) Run() {
	log.Info().Msg("start")

	defer a.pg.Close()
	defer a.ch.Close()

	if err := a.ch.Cleanup(); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if err := a.pg.Cleanup(); err != nil {
		log.Error().Msg(err.Error())
		return
	}

	start := time.Now()
	res, err := a.rest.GetBlockData(a.blockNumber) // 9989379
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.pg.InsertBlockData(res); err != nil {
			log.Error().Msg(err.Error())
			return
		}
		transfersPG, err := a.pg.SelectTransfers()
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		s, _ := json.MarshalIndent(transfersPG, "", "\t")
		log.Info().RawJSON("result", s).Msg("print from postgres transfer_view")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.ch.InsertBlockData(res); err != nil {
			log.Error().Msg(err.Error())
			return
		}
		transfersCH, err := a.ch.SelectTransfers()
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		s, _ := json.MarshalIndent(transfersCH, "", "\t")
		log.Info().RawJSON("result", s).Msg("print from clickhouse transfer_view")
	}()

	wg.Wait()
	log.Info().Dur("duration", time.Since(start)).Msg("finish")
}
