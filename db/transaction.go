package db

import (
	"context"

	"gorm.io/gorm"
)

type Transaction struct {
	MySQL *gorm.DB
}

func NewTransaction() *Transaction {
	return &Transaction{MySQL: MySQL}
}

func (t *Transaction) Begin(ctx context.Context) *gorm.DB {
	return t.MySQL.Begin().WithContext(ctx)
}

func (t *Transaction) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (t *Transaction) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}
