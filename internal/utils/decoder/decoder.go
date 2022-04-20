package decoder

import (
	"encoding/base64"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
)

var errBadInterfaceCast = errors.New("failed to cast interface to types.Tx")

type Decoder struct {
	interfacesRegistry codectypes.InterfaceRegistry
}

// RegisterInterfaces register decoding interface to the decoder by using the provided interface
// registry.
func (decoder *Decoder) RegisterInterfaces(registry func(registry codectypes.InterfaceRegistry)) *Decoder {
	registry(decoder.interfacesRegistry)

	return decoder
}

// Decode decodes the base64-encoded transaction bytes to Tx
func (decoder *Decoder) DecodeBase64(base64Tx string) (Tx, error) {
	txBytes, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		return Tx{}, err
	}

	return decoder.Decode(txBytes)
}

// Decode decodes the transaction bytes to Tx
func (decoder *Decoder) Decode(txBytes []byte) (Tx, error) {
	marshaler := codec.NewProtoCodec(decoder.interfacesRegistry)

	tx, err := DefaultTxDecoder(marshaler)(txBytes)
	if err != nil {
		return Tx{}, err
	}
	res, ok := tx.(cosmostypes.Tx)
	if !ok {
		return Tx{}, errBadInterfaceCast
	}
	return Tx{
		res,
		marshaler,
	}, nil
}

// NewDecoder creates a new decoder
func NewDecoder() *Decoder {
	interfaceRegistry := codectypes.NewInterfaceRegistry()

	return &Decoder{
		interfaceRegistry,
	}
}

// DefaultDecoder is a decoder with all Cosmos builtin modules interfaces registered
var DefaultDecoder = NewDecoder().RegisterInterfaces(RegisterDefaultInterfaces)
