// Package auth deriva Master Key e Auth Hash (espelha frontend e backend).
package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/crypto"
)

// ClientKDF parâmetros guardados na BD por utilizador.
type ClientKDF struct {
	Salt    []byte
	Time    int
	Memory  int
	Threads int
}

// LoginResult contém a Master Key (32 bytes) e token de sessão.
type LoginResult struct {
	MasterKey []byte
	Token     string
}

// FetchKDF obtém salt/parâmetros para um email.
func FetchKDF(baseURL, email string) (ClientKDF, error) {
	u := fmt.Sprintf("%s/api/auth/kdf?email=%s", baseURL, url.QueryEscape(email))
	res, err := http.Get(u)
	if err != nil {
		return ClientKDF{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ClientKDF{}, fmt.Errorf("utilizador não encontrado (%d)", res.StatusCode)
	}
	var j struct {
		Salt    string `json:"salt"`
		Time    int    `json:"time"`
		Memory  int    `json:"memory"`
		Threads int    `json:"threads"`
	}
	if err := json.NewDecoder(res.Body).Decode(&j); err != nil {
		return ClientKDF{}, err
	}
	salt, err := base64.StdEncoding.DecodeString(j.Salt)
	if err != nil {
		return ClientKDF{}, err
	}
	return ClientKDF{Salt: salt, Time: j.Time, Memory: j.Memory, Threads: j.Threads}, nil
}

// Login deriva chaves e obtém token Bearer.
func Login(baseURL, email, password string) (LoginResult, error) {
	kdf, err := FetchKDF(baseURL, email)
	if err != nil {
		return LoginResult{}, err
	}
	p := crypto.KDFParams{
		TimeCost:  uint32(kdf.Time),
		MemoryKiB: uint32(kdf.Memory),
		Threads:   uint8(kdf.Threads),
		KeyLen:    crypto.KeySize,
	}
	mk := crypto.DeriveMasterKey([]byte(password), kdf.Salt, p)
	authHash := crypto.DeriveAuthHash(mk, []byte(password), p)

	body, _ := json.Marshal(map[string]string{
		"email":     email,
		"auth_hash": base64.StdEncoding.EncodeToString(authHash),
	})
	res, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return LoginResult{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return LoginResult{}, fmt.Errorf("login falhou (%d): %s", res.StatusCode, string(b))
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{MasterKey: mk, Token: out.Token}, nil
}

// MasterKeyFromPassword re-deriva a Master Key (para run após device registado).
func MasterKeyFromPassword(email, password, baseURL string) ([]byte, error) {
	kdf, err := FetchKDF(baseURL, email)
	if err != nil {
		return nil, err
	}
	p := crypto.KDFParams{
		TimeCost:  uint32(kdf.Time),
		MemoryKiB: uint32(kdf.Memory),
		Threads:   uint8(kdf.Threads),
		KeyLen:    crypto.KeySize,
	}
	return crypto.DeriveMasterKey([]byte(password), kdf.Salt, p), nil
}
