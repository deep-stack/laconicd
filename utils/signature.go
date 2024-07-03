package utils

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/statechannels/go-nitro/crypto"
)

func DecodeEthereumAddress(message []byte, sig string) (string, error) {
	if len(sig) > 2 && sig[:2] == "0x" {
		sig = sig[2:]
	}

	signature := crypto.SplitSignature(common.Hex2Bytes(sig))
	ethereumAddress, err := crypto.RecoverEthereumMessageSigner(message, signature)

	return ethereumAddress.String(), err
}
