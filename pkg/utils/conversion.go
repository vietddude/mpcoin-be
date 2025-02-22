package utils

import (
	"encoding/asn1"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type ECDSASignature struct {
	R, S *big.Int
}

var (
	secp256k1N     = crypto.S256().Params().N
	secp256k1halfN = new(big.Int).Rsh(secp256k1N, 1)
)

func ConvertDERToEthSignature(derSig, txHash []byte, expectedPubKeyHex string) ([]byte, error) {
	var sig ECDSASignature
	if _, err := asn1.Unmarshal(derSig, &sig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DER signature: %v", err)
	}

	// Normalize S (Ethereum yêu cầu low S values)
	if sig.S.Cmp(secp256k1halfN) > 0 {
		sig.S = new(big.Int).Sub(secp256k1N, sig.S)
	}

	// Pad R và S về 32 byte
	rPadded := padTo32(sig.R.Bytes())
	sPadded := padTo32(sig.S.Bytes())
	sigBytes := append(rPadded, sPadded...)

	// Thử recovery ID 0 và 1
	for v := byte(0); v < 2; v++ {
		candidate := append(sigBytes, v)
		recoveredPub, err := crypto.SigToPub(txHash, candidate)
		if err != nil {
			continue
		}
		if strings.EqualFold(crypto.PubkeyToAddress(*recoveredPub).Hex(), expectedPubKeyHex) {
			return candidate, nil
		}
	}

	return nil, fmt.Errorf("failed to find valid recovery id")
}

func padTo32(b []byte) []byte {
	p := make([]byte, 32)
	copy(p[32-len(b):], b)
	return p
}
