package usecasetest

import (
	"context"
	"errors"
	"pack_optimizer/internal/domain"
	"pack_optimizer/internal/usecase/packusecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

// dynamicMockRepo allows us to configure the packs returned by GetAllPacks.
type dynamicMockRepo struct {
	packs []domain.Pack
}

func (m *dynamicMockRepo) GetAllPacks(_ context.Context) ([]domain.Pack, error) {
	return m.packs, nil
}

// A mock repository that always returns an error.
type errorMockRepo struct{}

func (m *errorMockRepo) GetAllPacks(_ context.Context) ([]domain.Pack, error) {
	return nil, errors.New("database connection failed")
}

func TestCalculatePacks(t *testing.T) {
	uc := packusecase.NewPackUseCase(&dynamicMockRepo{packs: []domain.Pack{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
	}})

	tests := []struct {
		name     string
		orderQty int
		expected packusecase.CalculatePacksOutput
		err      error
	}{
		{
			name:     "Order 0 -> should return error",
			orderQty: 0,
			expected: packusecase.CalculatePacksOutput{},
			err:      errors.New("order quantity must be greater than 0"),
		},
		{
			name:     "Negative order -> should return error",
			orderQty: -100,
			expected: packusecase.CalculatePacksOutput{},
			err:      errors.New("order quantity must be greater than 0"),
		},
		{
			name:     "Order 1 -> 1 x 250",
			orderQty: 1,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     250,
				RemainingItems: 249,
				TotalPacks:     1,
				Packs:          []packusecase.Pack{{Size: 250, Count: 1}},
			},
			err: nil,
		},
		{
			name:     "Order 250 -> 1 x 250",
			orderQty: 250,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     250,
				RemainingItems: 0,
				TotalPacks:     1,
				Packs:          []packusecase.Pack{{Size: 250, Count: 1}},
			},
			err: nil,
		},
		{
			name:     "Order 251 -> 1 x 500",
			orderQty: 251,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     500,
				RemainingItems: 249,
				TotalPacks:     1,
				Packs:          []packusecase.Pack{{Size: 500, Count: 1}},
			},
			err: nil,
		},
		{
			name:     "Order 501 -> 1 x 500 + 1 x 250",
			orderQty: 501,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     750,
				RemainingItems: 249,
				TotalPacks:     2,
				Packs:          []packusecase.Pack{{Size: 250, Count: 1}, {Size: 500, Count: 1}},
			},
			err: nil,
		},
		{
			name:     "Order 12001 -> 2 x 5000 + 1 x 2000 + 1 x 250",
			orderQty: 12001,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     12250,
				RemainingItems: 249,
				TotalPacks:     4,
				Packs: []packusecase.Pack{
					{Size: 250, Count: 1},
					{Size: 2000, Count: 1},
					{Size: 5000, Count: 2},
				},
			},
			err: nil,
		},
		{
			name:     "Order exactly a large pack -> 1 x 5000",
			orderQty: 5000,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     5000,
				RemainingItems: 0,
				TotalPacks:     1,
				Packs:          []packusecase.Pack{{Size: 5000, Count: 1}},
			},
			err: nil,
		},
		{
			name:     "Order not divisible by any pack size",
			orderQty: 1234,
			expected: packusecase.CalculatePacksOutput{
				TotalItems:     1250,
				RemainingItems: 16,
				TotalPacks:     2,
				Packs:          []packusecase.Pack{{Size: 250, Count: 1}, {Size: 1000, Count: 1}},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.CalculatePacks(context.Background(), tt.orderQty)

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expected.Packs, result.Packs)
				assert.Equal(t, tt.expected.TotalItems, result.TotalItems)
				assert.Equal(t, tt.expected.RemainingItems, result.RemainingItems)
				assert.Equal(t, tt.expected.TotalPacks, result.TotalPacks)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCalculatePacks_RepoError(t *testing.T) {
	uc := packusecase.NewPackUseCase(&errorMockRepo{})

	// Test case for repository error
	result, err := uc.CalculatePacks(context.Background(), 10)
	assert.Error(t, err)
	assert.Equal(t, "use case failed to get packs: database connection failed", err.Error())
	assert.Equal(t, packusecase.CalculatePacksOutput{}, result)
}
