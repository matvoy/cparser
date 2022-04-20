package http

import (
	"encoding/json"
	"time"
)

type BlockResponse struct {
	BlockID struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
	Block struct {
		Header struct {
			Height          string    `json:"height"`
			Time            time.Time `json:"time"`
			ProposerAddress string    `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
	} `json:"block"`
}

type TxResponse struct {
	TotalCount string `json:"total_count"`
	Txs        []struct {
		Height    string    `json:"height"`
		TxHash    string    `json:"txhash"`
		Timestamp time.Time `json:"timestamp"`
		Tx        struct {
			Value struct {
				Msg []struct {
					Type  string          `json:"type"` // cosmos-sdk/MsgSend
					Value json.RawMessage `json:"value"`
				} `json:"msg"`
				Fee struct {
					Amount []Amount `json:"amount"`
				} `json:"fee"`
			} `json:"value"`
		} `json:"tx"`
	} `json:"txs"`
}

type MsgSend struct {
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Amount      []Amount `json:"amount"`
}

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Message struct {
	Type        string   `json:"@type" mapstructure:"@type"`
	FromAddress string   `json:"from_address" mapstructure:"from_address"`
	ToAddress   string   `json:"to_address" mapstructure:"to_address"`
	Amount      []Amount `json:"amount" mapstructure:"amount"`
}
