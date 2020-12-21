package model

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math/big"
)

type (
	JWKS struct {
		Keys []*JSONWebKeys `json:"keys"`
	}
)

type JSONWebKeys struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

func NewJsonToJWKS(buf []byte) (*JSONWebKeys, error) {
	jwks := &JSONWebKeys{}
	if err := json.Unmarshal(buf, jwks); err != nil {
		return nil, err
	}
	return jwks, nil
}

func (jwk *JSONWebKeys) GetPublicKey() *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}

func (jwk *JSONWebKeys) ToJson() string {
	buf, _ := json.Marshal(jwk)
	return string(buf)
}
