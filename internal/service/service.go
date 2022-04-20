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

var txsGlob = []string{
	"ClYKUgobL2Nvc21vcy5nb3YudjFiZXRhMS5Nc2dWb3RlEjMIQRItY29zbW9zMXljNmczcHF1YzVucTd6cWRsNnVsNGR4bnhmaDc3NHd1dWVxMmRnGAESABJnClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiED3QPf5hB9k70tomhgqxWAD53Hq/Kidgc0ZMjywHlEgQUSBAoCCH8YFhITCg0KBXVhdG9tEgQ2MjUwEJChDxpA5HqQHjv+592qZET4/IznVRH8GH2rFXQpA6eAEIlgr1Q9sV8XGnsmLqRv7nPAqvyyt5sEC9nA6FG047GacIXQiA==",
	"CpMBCpABChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnAKLWNvc21vczF0NXUwamZnM2xqc2pyaDJtOWU0N2Q0bnkyaGVhN2VlaHhyemRnZBItY29zbW9zMWthajNha3ZqbHh6NHBsMDVrd3Bmc2EwZmc0Z3EybHJkNHg0Nm5mGhAKBXVhdG9tEgcxNDk5NTMxEmkKUgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQLqGodLEwWelzPDXpDLKlhlpGtk7ZDSVVR8B8pSOSf8YhIECgIIfxi8yCMSEwoNCgV1YXRvbRIEMjUwMBCImAUaQGMiQt+hGEPcfpJVuQsv6gC1IQP56RQ/QAU+7lwYgiPiCVns+y2dbjX9DvQKIScTPo1bT3hVgVwB59cCj1BqhlY=",
	"ClYKUgobL2Nvc21vcy5nb3YudjFiZXRhMS5Nc2dWb3RlEjMIQBItY29zbW9zMXE4MjczejNtaGFsZ2twYzdnNHY3N241emN0aDRmNTVsdmVzcmN1GAESABJnClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECy4iEMuiLJit5ro5DiGDdNU6fsYdWVKFh2faDHNXulOoSBAoCCH8YDxITCg0KBXVhdG9tEgQ2MjUwEJChDxpAQd/G3XMC9jOJlg3jcURLISGlseVD/FBMo+edTXFOx81E0KY9L6PBkaEnv8mV25COp+48gbNPxfXrMRA4HSzjYA==",
	"CpUBCpABChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnAKLWNvc21vczFwdGRzdnFoMjc5YXhqY242c3ZlYWZzcmhweDQwMHdwanY2bTN2dxItY29zbW9zMXJ6NDg0anZjd2hzNTNtejBhbTNqaHJrMzB0ejdsdnprd2FuZGh3GhAKBXVhdG9tEgcxNDAwMDAwEgASaApRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohA6gkjOQ+Z+i+NzGOZTeqwajN4R/un2wme516+tbUXVrTEgQKAgh/GKsBEhMKDQoFdWF0b20SBDIwMDAQgPEEGkDQ4/7fgL+jSDJNJdqawsGe1IE0b3xhq7svKPU+r1+wi3E+115T+bz/8yBdqobs0KU/DHZyXcDqEH2SVy+JgE8r",
	"CqIBCp0BCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRJ2Ci1jb3Ntb3Mxa3V0enNkZ3pnZG1ndnk0enluNDNuaHBtNGgzbXE1cWtoNjVybTcSNGNvc21vc3ZhbG9wZXIxNTZncWY5ODM3dTdkNGM0Njc4eXQzcmw0bHM5YzV2dXVyc3JyemYaDwoFdWF0b20SBjIyMDAwMBIAEmcKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQPwmhoykhK0GPWuY0LCuoizdVMgdMz05caIHaMKAccaNBIECgIIfxgTEhMKDQoFdWF0b20SBDYyNTAQkKEPGkDq7BWoQl6ikDjpX8RMWCESR7d3RzdZ6WjrClRRO7PqmQBB4bkzsLQG/UU/OZm8jBI8I51or/wiYMFw6Lddwjeq",
	"CrkBCp8BCiMvY29zbW9zLnN0YWtpbmcudjFiZXRhMS5Nc2dEZWxlZ2F0ZRJ4Ci1jb3Ntb3MxMHE5N3Rnczl6d3c5dzl3N3JhcTJtc3lhM3Y0YzNmZmZsdzZ3bXoSNGNvc21vc3ZhbG9wZXIxdGZsazMwbXE1dmdxamRseTkya2toaHEzcmFldjJobno2ZWV0ZTMaEQoFdWF0b20SCDIwNjU5MTE2EhVEZWxlZ2F0ZWQgd2l0aCBFeG9kdXMSZQpOCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAwPz3mGcwP1jfNo00jS6u1is9zCq+M2U4j9X8eNogT5TEgQKAggBEhMKDQoFdWF0b20SBDUwNzUQ+LEMGkBHLWlkRId/e7Co5ls4YWjcmuC1xOCdlYwi8FMJOnvLY3kAVyM/ITiFdDngVSXJsHT/Gm6HMCZkjaY9U2l/IOor",
	"CqMBCqABCjcvY29zbW9zLmRpc3RyaWJ1dGlvbi52MWJldGExLk1zZ1dpdGhkcmF3RGVsZWdhdG9yUmV3YXJkEmUKLWNvc21vczFwNGYyZmtmazU3NW14bHZlbjB5aGc3dm41YXJuZHowNms4NHY5chI0Y29zbW9zdmFsb3BlcjF0ZmxrMzBtcTV2Z3FqZGx5OTJra2hocTNyYWV2MmhuejZlZXRlMxJnClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDzu6S6x4OSiSE0CeC9J/hyErIHNRwYHeYYL9u/wG54jISBAoCCAEYNxITCg0KBXVhdG9tEgQzNTAwEODFCBpAiNQ8nFuzbzL7CsozNi701qgX6OIrapO9zMhfMFqBs/QzqYBRu8r814cbujv1XxQBRMNMJgo5JE62ruZ1HrI9bQ==",
}

var dd = &models.Block{
	BlockNumber:     222,
	Hash:            "123123123",
	ProposerAddress: "434343434",
	CreatedDate:     time.Now(),
	Txs: []*models.Tx{
		{
			Hash:            "dddddd",
			BlockNumber:     222,
			Status:          "323",
			Fee:             41344,
			FeeCurrency:     "rettr",
			FeePayerAddress: "tttttt",
			CreatedDate:     time.Now(),
			Transfers: []*models.Transfer{
				{
					ID:          "dddddd",
					TxHash:      "dddddd",
					FromAddress: "eeee",
					ToAddress:   "ddddd",
					Amount:      22,
					Currency:    "saasas",
				},
				{
					ID:          "dddddd",
					TxHash:      "dddddd",
					FromAddress: "eeee",
					ToAddress:   "ddddd",
					Amount:      22,
					Currency:    "saasas",
				},
			},
		},
		{
			Hash:            "dddddd2",
			BlockNumber:     222,
			Status:          "323",
			Fee:             41344,
			FeeCurrency:     "rettr",
			FeePayerAddress: "tttttt",
			CreatedDate:     time.Now(),
			Transfers: []*models.Transfer{
				{
					ID:          "dddddd",
					TxHash:      "dddddd2",
					FromAddress: "jjjjjjj",
					ToAddress:   "kkkkkkkkk",
					Amount:      333,
					Currency:    "llllll",
				},
			},
		},
	},
}
