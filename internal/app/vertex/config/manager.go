package config

import (
	"errors"
	"github.com/lucidnet/lucidnet/internal/pkg/secret"
	"gorm.io/gorm"
)

type Manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) Init() error {
	jwtSigningKeyString, err := secret.GenerateSecret(64)

	if err != nil {
		return err
	}

	err = m.Set("jwt_signing_key", jwtSigningKeyString)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Set(key string, value string) error {
	config := &Config{
		Key:   key,
		Value: value,
	}

	result := m.db.Where(Config{Key: key}).FirstOrCreate(config)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *Manager) Get(key string, defaultValue string) (string, error) {
	var config Config
	txn := m.db.Where(Config{Key: key}).First(&config)

	if txn.Error != nil {
		if errors.Is(txn.Error, gorm.ErrRecordNotFound) {
			return defaultValue, nil
		}

		return "", txn.Error
	}

	return config.Value, nil
}

func (m *Manager) GetJwtSigningKey(defaultValue string) ([]byte, error) {
	jwtSigningKeyString, err := m.Get("jwt_signing_key", defaultValue)

	if err != nil {
		return nil, err
	}

	return []byte(jwtSigningKeyString), nil
}

func (m *Manager) SetVertexEndpoint(vertexEndpoint string) error {
	return m.Set("vertex_endpoint", vertexEndpoint)
}

func (m *Manager) GetVertexEndpoint() (string, error) {
	return m.Get("vertex_endpoint", "")
}
