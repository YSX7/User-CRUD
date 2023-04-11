package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"time"

	"vuegolang/graph/model"
)

type Token struct {
	Value string
	rsa.PrivateKey
	Expired time.Time
}

type Claims struct {
	Id      int
	Expired int64
}

// New маршализирует клаймсы, шифрует их и возвращает их вместе с приватным ключом и сроком годности
func New(id int) (Token, error) {
	t := Token{}
	expireUnix := time.Now().AddDate(1, 0, 0)
	claims := Claims{
		Id:      id,
		Expired: expireUnix.Unix(),
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	claimsBytes, errClaims := json.Marshal(claims)
	if errClaims != nil {
		return t, fmt.Errorf("error encoding claims: %w", err)
	}

	encryptedBytes, err := EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&privateKey.PublicKey,
		claimsBytes,
		nil,
	)
	if err != nil {
		panic(err)
	}

	tokenStr := base64.StdEncoding.EncodeToString(encryptedBytes)

	t = Token{
		Value:      tokenStr,
		PrivateKey: *privateKey,
		Expired:    expireUnix,
	}

	return t, nil
}

func DecryptBytes(code string, key rsa.PrivateKey) (Claims, error) {
	var claims Claims
	b, errDecode := base64.StdEncoding.DecodeString(code)
	if errDecode != nil {
		return Claims{}, errDecode
	}

	decrypted, errDecrypt := DecryptOAEP(sha256.New(), rand.Reader, &key, b, nil)
	if errDecrypt != nil {
		return Claims{}, errDecrypt
	}

	errUnmarshal := json.Unmarshal(decrypted, &claims)
	if errUnmarshal != nil {
		return Claims{}, errUnmarshal
	}
	return claims, nil
}

func EncryptOAEP(hash hash.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func DecryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

// ToJsonString превращает слайс байтов с строку, запихивает в model.AuthInfo, маршализирует и кодирует в base64
func (t *Token) ToJsonString() (string, error) {
	tokenStr := t.Value

	authInfo := &model.AuthInfo{Token: &tokenStr}
	b, errJson := json.Marshal(authInfo)
	if errJson != nil {
		return "", errJson
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func FromJsonString(str string) (string, error) {
	if str == "" {
		return "", errors.New("empty string")
	}

	b, errDecode := base64.StdEncoding.DecodeString(str)
	authInfo := new(model.AuthInfo)
	if errDecode != nil {
		return "", errDecode
	}

	if err := json.Unmarshal(b, authInfo); err != nil {
		return "", err
	}

	return *authInfo.Token, nil
}
