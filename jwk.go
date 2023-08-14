package adyen_encrypt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
)

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
	Use string `json:"use"`
}

var (
	errSplit = errors.New("invalid key format: missing '|'")
)

func DefaultJWK() *JWK {
	return &JWK{
		Kty: "RSA",
		Kid: "asf-key",
		Alg: "RSA-OAEP",
		Use: "sig",
	}
}

func (jwk *JWK) Marshal() []byte {
	m, _ := json.Marshal(jwk)
	return m
}

func (jwk *JWK) ParseAdyenKey(key string) error {
	parts := strings.Split(key, "|")
	if len(parts) < 2 {
		return errSplit
	}
	decodedExponent := HexDecode(parts[0])
	decodedKey := HexDecode(parts[1])
	encodedExponent := EncodeToBase64(decodedExponent)
	encodedKey := EncodeToBase64(decodedKey)
	jwk.E = encodedExponent
	jwk.N = encodedKey
	return nil
}

func (jwk *JWK) JWKToPem() *rsa.PublicKey {
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil
	}
	nBytes, _ := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil
	}

	rsaPub := &rsa.PublicKey{
		N: big.NewInt(0).SetBytes(nBytes),
		E: int(big.NewInt(0).SetBytes(eBytes).Uint64()),
	}
	return rsaPub
}
