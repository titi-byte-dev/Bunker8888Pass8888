// Package api cliente HTTP (Bearer + mTLS) para a API AegisPass.
package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Client fala com a API via mTLS.
type Client struct {
	mtls *http.Client
	base string
}

// NewMTLS cria cliente com certificado de dispositivo registado.
func NewMTLS(baseURL, certFile, keyFile, caFile string) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("CA PEM inválido")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            pool,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: os.Getenv("AEGIS_MTLS_INSECURE") == "1", // só dev
		},
	}
	return &Client{
		mtls: &http.Client{Transport: tr},
		base: baseURL,
	}, nil
}

// VaultItem metadados + blob opcional.
type VaultItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Blob      string `json:"blob"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListVaultItems devolve metadados dos itens (mTLS).
func (c *Client) ListVaultItems(itemType string) ([]VaultItem, error) {
	path := c.base + "/api/cli/vault/items"
	if itemType != "" {
		path += "?type=" + itemType
	}
	res, err := c.mtls.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return decodeItems(res)
}

// GetVaultItem devolve um item com blob cifrado (mTLS).
func (c *Client) GetVaultItem(id string) (*VaultItem, error) {
	res, err := c.mtls.Get(c.base + "/api/cli/vault/items/" + id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("GET item (%d): %s", res.StatusCode, string(b))
	}
	var out VaultItem
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func decodeItems(res *http.Response) ([]VaultItem, error) {
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("list (%d): %s", res.StatusCode, string(b))
	}
	var wrap struct {
		Items []VaultItem `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&wrap); err != nil {
		return nil, err
	}
	return wrap.Items, nil
}

// RegisterDeviceResponse resposta POST /api/cli/devices.
type RegisterDeviceResponse struct {
	CertPEM string `json:"cert_pem"`
	CAPEM   string `json:"ca_pem"`
}

// RegisterDevice envia CSR com token Bearer e recebe certificado assinado.
func RegisterDevice(baseURL, token, name, csrPEM string) (RegisterDeviceResponse, error) {
	body, _ := json.Marshal(map[string]string{"name": name, "csr_pem": csrPEM})
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/cli/devices", bytes.NewReader(body))
	if err != nil {
		return RegisterDeviceResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RegisterDeviceResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		return RegisterDeviceResponse{}, fmt.Errorf("registo dispositivo (%d): %s", res.StatusCode, string(b))
	}
	var out RegisterDeviceResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return RegisterDeviceResponse{}, err
	}
	return out, nil
}
