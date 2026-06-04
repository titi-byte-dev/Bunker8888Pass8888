// Package vault decifra blobs do cofre no cliente (Zero-Knowledge).
package vault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/crypto"
)

const payloadVersion = 1

type payloadEnvelope struct {
	V    int             `json:"v"`
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

// LoginItem espelha o tipo frontend (campo password para injecção).
type LoginItem struct {
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url,omitempty"`
}

// OpenLogin decifra um blob base64 e extrai um item login.
func OpenLogin(masterKey, blobB64 []byte) (*LoginItem, error) {
	blob, err := base64.StdEncoding.DecodeString(string(blobB64))
	if err != nil {
		return nil, fmt.Errorf("blob base64 inválido: %w", err)
	}
	plain, err := crypto.Decrypt(masterKey, blob, nil)
	if err != nil {
		return nil, fmt.Errorf("decifragem falhou: %w", err)
	}
	var env payloadEnvelope
	if err := json.Unmarshal(plain, &env); err != nil {
		return nil, err
	}
	if env.V != payloadVersion {
		return nil, fmt.Errorf("versão não suportada: %d", env.V)
	}
	if env.Kind != "login" {
		return nil, fmt.Errorf("tipo %q não suportado pela CLI (use login)", env.Kind)
	}
	var item LoginItem
	if err := json.Unmarshal(env.Data, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// FieldValue devolve um campo do item login (password, username, title, url).
func FieldValue(item *LoginItem, field string) (string, error) {
	switch field {
	case "password":
		return item.Password, nil
	case "username":
		return item.Username, nil
	case "title":
		return item.Title, nil
	case "url":
		return item.URL, nil
	default:
		return "", fmt.Errorf("campo desconhecido: %s (use password, username, title, url)", field)
	}
}
