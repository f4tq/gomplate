package vault

import (
	"bytes"
	"encoding/json"
	"log"

	vaultapi "github.com/hashicorp/vault/api"
)

// logFatal is defined so log.Fatal calls can be overridden for testing
var logFatal = log.Fatal

// Vault -
type Vault struct {
	client *vaultapi.Client
}

// New -
func New() *Vault {
	vaultConfig := vaultapi.DefaultConfig()

	err := vaultConfig.ReadEnvironment()
	if err != nil {
		logFatal("Vault setup failed", err)
	}

	client, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		logFatal("Vault setup failed", err)
	}

	return &Vault{client}
}

// Login -
func (v *Vault) Login() {
	v.client.SetToken(v.GetToken())
}

// Logout -
func (v *Vault) Logout() {
}

// Read - returns the value of a given path. If no value is found at the given
// path, returns empty slice.
func (v *Vault) Read(path string) ([]byte, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(secret.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *Vault) Write(path string, data map[string]interface{}) ([]byte, error) {
	secret, err := v.client.Logical().Write(path, data)
	if secret == nil {
		return []byte{}, err
	}
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(secret.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
