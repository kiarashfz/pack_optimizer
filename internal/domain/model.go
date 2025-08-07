package domain

import "context"

type Pack struct {
	ID   uint `gorm:"primaryKey"`
	Size int  `gorm:"uniqueIndex"`
}

type PackRepository interface {
	GetAllPacks(ctx context.Context) ([]Pack, error)
}
