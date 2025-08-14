// Package domain defines the core data structures and interfaces for the application.
package domain

import "context"

type Pack struct {
	ID   uint `gorm:"primaryKey"`
	Size int  `gorm:"uniqueIndex"`
}

type PackRepository interface {
	GetAllPacks(ctx context.Context) ([]Pack, error)
}
