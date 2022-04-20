package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresDSN    string
	ClickHouseDSN1 string
	ClickHouseDSN2 string
	ClickHouseDSN3 string
	ClickHouseDSN4 string
	GrpcURL        string
	RestURL        string
	BlockNumber    uint64
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	c := new(Config)

	c.PostgresDSN = viper.GetString("PG_DSN")
	c.ClickHouseDSN1 = viper.GetString("CH_DSN_1")
	c.ClickHouseDSN2 = viper.GetString("CH_DSN_2")
	c.ClickHouseDSN3 = viper.GetString("CH_DSN_3")
	c.ClickHouseDSN4 = viper.GetString("CH_DSN_4")
	c.GrpcURL = viper.GetString("GRPC_URL")
	c.RestURL = viper.GetString("REST_URL")
	c.BlockNumber = viper.GetUint64("BLOCK_NUMBER")

	return c, nil
}
