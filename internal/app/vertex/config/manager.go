package config

import (
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

	jwtSigningKey := &Config{
		Key:   "jwt_signing_key",
		Value: jwtSigningKeyString,
	}

	result := m.db.Where(Config{Key: "jwt_signing_key"}).FirstOrCreate(&jwtSigningKey)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
