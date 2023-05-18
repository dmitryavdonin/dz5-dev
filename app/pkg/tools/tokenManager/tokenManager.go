package tokenManager

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
)

const (
	jwksPath = "./.well-known"
)

var (
	ErrSigingMethod   = "unexpected signing method: %v"
	ErrExpired        = "token is expired"
	ErrClaims         = "error parsing claims"
	ErrEmptyKey       = "empty key"
	ErrTokenIsInvalid = "token is invalid"
)

type TokenManager interface {
	NewJWT(input AuthInfo) (string, error)
	Parse(accessToken string) (result AuthInfo, err error)
	SigningMethod() string
}

type Manager struct {
	signingMethod jwt.SigningMethod
	ttl           time.Duration
	privateKey    *rsa.PrivateKey
	publickKey    *rsa.PublicKey
	keyID         string
}

func NewManager(ttl time.Duration) (*Manager, error) {
	m := &Manager{signingMethod: jwt.SigningMethodES256, ttl: ttl}
	err := m.generateRsaKey()
	return m, err
}

func NewOnlyParserManager(publicKeyUri string) (m *Manager, err error) {
	m = &Manager{signingMethod: jwt.SigningMethodES256}

	resp, err := http.Get(publicKeyUri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	jwkJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	set, err := jwk.Parse([]byte(jwkJSON))
	if err != nil {
		return
	}

	for it := set.Iterate(context.Background()); it.Next(context.Background()); {
		pair := it.Pair()
		key := pair.Value.(jwk.Key)

		var rawkey interface{} // This is the raw key, like *rsa.PrivateKey or *ecdsa.PrivateKey
		if err = key.Raw(&rawkey); err != nil {
			log.Printf("failed to create public key: %s", err)
			return
		}

		// We know this is an RSA Key so...
		rsa, ok := rawkey.(*rsa.PublicKey)
		if !ok {
			panic(fmt.Sprintf("expected ras key, got %T", rawkey))
		}

		m.publickKey = rsa

		break
	}
	return
}

func (m *Manager) generateRsaKey() (err error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return fmt.Errorf("cannot generate RSA key: %s", err.Error())
	}
	m.privateKey = privatekey
	m.publickKey = &privatekey.PublicKey

	err = m.saveJWKs()
	return
}

func (m *Manager) SigningMethod() string {
	return m.signingMethod.Alg()
}

type JWKs struct {
	Keys []map[string]interface{} `json:"keys"`
}

func (m *Manager) saveJWKs() (err error) {
	set, err := jwk.New(m.publickKey)
	if err != nil {
		return
	}

	err = jwk.AssignKeyID(set)
	if err != nil {
		return
	}

	m.keyID = set.KeyID()

	res, err := set.AsMap(context.Background())
	if err != nil {
		return
	}

	jwks := JWKs{
		Keys: []map[string]interface{}{res},
	}

	bytes, err := json.MarshalIndent(jwks, "", "")
	if err != nil {
		return
	}

	file, err := os.Create(jwksPath + "/jwks.json")
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(bytes)

	return
}

func (m *Manager) NewJWT(input AuthInfo) (string, error) {
	if m.privateKey == nil {
		return "", fmt.Errorf("no private key provided")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &tokenClaims{
		input.UserID,
		input.Login,
		input.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	token.Header["kid"] = m.keyID

	return token.SignedString(m.privateKey)
}

func (m *Manager) Parse(accessToken string) (result AuthInfo, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return m.publickKey, nil
	})
	if err != nil {
		return
	}

	if !token.Valid {
		err = errors.New(ErrTokenIsInvalid)
		return
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return result, errors.New(ErrClaims)
	}

	result.UserID = claims.UserID
	result.Login = claims.Login
	result.Role = claims.Role

	return
}
