package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/calvinlauyh/cosmosutils"
	"github.com/google/uuid"
	"github.com/matvoy/cparser/internal/models"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

const (
	decodedMsgSendType = "/cosmos.bank.v1beta1.MsgSend"
	apiMsgSendType     = "cosmos-sdk/MsgSend"
)

// map to general Block model from /block and /tx api responses
func mapToDomainWithTxResponse(b *BlockResponse, t *TxResponse) *models.Block {
	if b == nil || t == nil {
		return nil
	}
	res := &models.Block{
		Hash:            b.BlockID.Hash,
		ProposerAddress: b.Block.Header.ProposerAddress,
		CreatedDate:     b.Block.Header.Time,
	}
	res.BlockNumber, _ = strconv.ParseUint(b.Block.Header.Height, 10, 64)
	res.Txs = mapTxToDomain(t)
	return res
}

// map to general Block model with parsing base64 encoded txs from /block response
func mapToDomainWithTxDecoding(b *BlockResponse) *models.Block {
	if b == nil {
		return nil
	}
	res := &models.Block{
		Hash:            b.BlockID.Hash,
		ProposerAddress: b.Block.Header.ProposerAddress,
		CreatedDate:     b.Block.Header.Time,
	}
	res.BlockNumber, _ = strconv.ParseUint(b.Block.Header.Height, 10, 64)
	res.Txs = parseEncodedTxs(res.BlockNumber, b.Block.Data.Txs)
	return res
}

// map /txs response to general Tx
func mapTxToDomain(in *TxResponse) []*models.Tx {
	if in == nil {
		return nil
	}
	out := make([]*models.Tx, 0, len(in.Txs))
	for _, item := range in.Txs {
		tmp := &models.Tx{
			Hash:        item.TxHash,
			Status:      "Success", // How should I check the status?? Looks like every completed tx is Success
			CreatedDate: item.Timestamp,
			Transfers:   []*models.Transfer{},
		}
		tmp.BlockNumber, _ = strconv.ParseUint(item.Height, 10, 64)
		if len(item.Tx.Value.Fee.Amount) > 0 {
			tmp.Fee, _ = strconv.ParseUint(item.Tx.Value.Fee.Amount[0].Amount, 10, 64)
			tmp.FeeCurrency = item.Tx.Value.Fee.Amount[0].Denom
			tmp.FeePayerAddress = "allah" // how to find the first signer address?
		}
		for _, message := range item.Tx.Value.Msg {
			if message.Type != apiMsgSendType {
				continue
			}
			var tmpMsg MsgSend
			if err := json.Unmarshal(message.Value, &tmpMsg); err != nil {
				log.Error().Msg(err.Error())
				continue
			}
			tmpTransfer := &models.Transfer{
				ID:          uuid.New().String(),
				TxHash:      item.TxHash,
				FromAddress: tmpMsg.FromAddress,
				ToAddress:   tmpMsg.ToAddress,
			}
			if len(tmpMsg.Amount) > 0 {
				tmpTransfer.Amount, _ = strconv.ParseUint(tmpMsg.Amount[0].Amount, 10, 64)
				tmpTransfer.Currency = tmpMsg.Amount[0].Denom
			}
			tmp.Transfers = append(tmp.Transfers, tmpTransfer)
		}
		out = append(out, tmp)
	}
	return out
}

// parsing encoded txs
func parseEncodedTxs(height uint64, txs []string) []*models.Tx {
	if height == 0 || len(txs) == 0 {
		return nil
	}

	dec := cosmosutils.DefaultDecoder
	result := make([]*models.Tx, 0, len(txs))

	for _, txBase64 := range txs {
		txDecoded, err := dec.DecodeBase64(txBase64)
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}
		// Get Tx Hash.
		// That's WRONG, but I haven't found the right way for this type of response.
		// When I'm decoding from base64 before result is also wrong
		hash := sha256.Sum256([]byte(txBase64))
		hashStr := hex.EncodeToString(hash[:])

		txDecodedBytes, _ := txDecoded.MarshalToJSON()
		var cosmosTx *cosmosutils.CosmosTx
		err = json.Unmarshal(txDecodedBytes, &cosmosTx)
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}
		tx := mapCosmosTxToDomain(height, hashStr, cosmosTx)

		result = append(result, tx)
	}

	return result
}

// map to general Tx model
func mapCosmosTxToDomain(height uint64, hash string, in *cosmosutils.CosmosTx) *models.Tx {
	if in == nil || height == 0 || len(hash) == 0 {
		return nil
	}
	out := &models.Tx{
		Hash:        hash, // Wrong
		BlockNumber: height,
		Status:      "Success",  // How should I check the status?? Looks like every completed tx is Success
		CreatedDate: time.Now(), // Wrong again, response doesn't have any datetimes
		Transfers:   []*models.Transfer{},
	}
	if len(in.AuthInfo.Fee.Amount) > 0 {
		out.Fee, _ = strconv.ParseUint(in.AuthInfo.Fee.Amount[0].Amount, 10, 64)
		out.FeeCurrency = in.AuthInfo.Fee.Amount[0].Denom
		out.FeePayerAddress = in.AuthInfo.Fee.Payer // if empty how to find the first signer address?
	}
	if len(in.Body.Messages) > 0 {
		for _, message := range in.Body.Messages {
			if v, ok := message["@type"].(string); ok && v == decodedMsgSendType {
				var tmp Message
				mapstructure.Decode(message, &tmp)
				transfer := &models.Transfer{
					ID:          uuid.New().String(),
					TxHash:      hash,
					FromAddress: tmp.FromAddress,
					ToAddress:   tmp.ToAddress,
				}
				if len(tmp.Amount) > 0 {
					transfer.Amount, _ = strconv.ParseUint(tmp.Amount[0].Amount, 10, 64)
					transfer.Currency = tmp.Amount[0].Denom
				}
				out.Transfers = append(out.Transfers, transfer)
			}
		}
	}
	return out
}
