// Package sqlrepo provides database operations for pack-related data.
package sqlrepo

import (
	"context"
	"pack_optimizer/internal/domain"

	"gorm.io/gorm"
)

// PackRepo is a repository that provides methods to interact with the Pack data in the database.
type PackRepo struct {
	db *gorm.DB // db is the GORM database connection.
}

// NewPackRepo creates a new instance of PackRepo.
func NewPackRepo(db *gorm.DB) domain.PackRepository {
	return &PackRepo{db: db}
}

// GetAllPacks retrieves all packs from the database, ordered by size in ascending order.
func (r *PackRepo) GetAllPacks(ctx context.Context) ([]domain.Pack, error) {
	var packs []domain.Pack
	err := r.db.WithContext(ctx).Select("size").Order("size ASC").Find(&packs).Error
	return packs, err
}
