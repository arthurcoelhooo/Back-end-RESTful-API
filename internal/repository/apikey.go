package repository

import (
	"context"

	"gorm.io/gorm"

	"receitas-api/internal/model"
)

// APIKeyRepository guarda apenas o HASH das chaves, nunca a chave em texto puro.
type APIKeyRepository interface {
	ExisteHash(ctx context.Context, hash string) (bool, error)
	Criar(ctx context.Context, nome, hash string) error
}

type APIKeyDAO struct {
	db *gorm.DB
}

func NewAPIKeyDAO(db *gorm.DB) *APIKeyDAO {
	return &APIKeyDAO{db: db}
}

func (d *APIKeyDAO) ExisteHash(ctx context.Context, hash string) (bool, error) {
	var n int64
	err := d.db.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("key_hash = ?", hash).
		Count(&n).Error
	return n > 0, err
}

func (d *APIKeyDAO) Criar(ctx context.Context, nome, hash string) error {
	return d.db.WithContext(ctx).Create(&model.APIKey{Nome: nome, KeyHash: hash}).Error
}
