package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/matvoy/cparser/internal/config"
	"github.com/matvoy/cparser/internal/service"
)

func main() {
	cfg, err := config.InitConfig()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if err != nil {
		log.Fatal().Err(err)
	}
	service.NewApp(cfg).Run()
}
