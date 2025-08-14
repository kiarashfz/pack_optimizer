package usecasetest

import (
	"context"
	"errors"
	"fmt"
	"pack_optimizer/internal/domain"
	"pack_optimizer/internal/usecase/packusecase"
	"testing"

	crand "crypto/rand"
	"math/big"

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

func TestCalculatePacks_StaticCases(t *testing.T) {
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
			err:      errors.New("order quantity must be greater than 0"),
		},
		{
			name:     "Negative order -> should return error",
			orderQty: -100,
			err:      errors.New("order quantity must be greater than 0"),
		},
		{
			name:     "Order 250 -> 1 x 250",
			orderQty: 250,
			expected: packusecase.CalculatePacksOutput{
				TotalItems: 250, TotalPacks: 1, Packs: []packusecase.Pack{{Size: 250, Count: 1}},
			},
		},
		{
			name:     "Order 251 -> 1 x 500",
			orderQty: 251,
			expected: packusecase.CalculatePacksOutput{
				TotalItems: 500, TotalPacks: 1, Packs: []packusecase.Pack{{Size: 500, Count: 1}},
			},
		},
		{
			name:     "Order 5000 -> exact large pack",
			orderQty: 5000,
			expected: packusecase.CalculatePacksOutput{
				TotalItems: 5000, TotalPacks: 1, Packs: []packusecase.Pack{{Size: 5000, Count: 1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.CalculatePacks(context.Background(), tt.orderQty)

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.TotalPacks, result.TotalPacks)
				assert.Equal(t, tt.expected.TotalItems, result.TotalItems)
				assert.ElementsMatch(t, tt.expected.Packs, result.Packs)
			}
		})
	}
}

func TestCalculatePacks_DynamicGeneratedCases(t *testing.T) {
	n, err := crand.Int(crand.Reader, big.NewInt(999999))
	if err != nil {
		t.Fatalf("failed to generate random number: %v", err)
		return
	}

	packSizes := []domain.Pack{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
	}
	minimumPackSize := packSizes[0].Size
	uc := packusecase.NewPackUseCase(&dynamicMockRepo{packs: packSizes})

	for i := 0; i < 20; i++ { // generate 20 random test cases
		orderQty := int(n.Int64()) + 1
		t.Run(fmt.Sprintf("Random_Order_%d", orderQty), func(t *testing.T) {
			result, err := uc.CalculatePacks(context.Background(), orderQty)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, result.TotalItems, orderQty, "total items must cover the order quantity")
			assert.Greater(t, result.TotalPacks, 0, "should have at least one pack")
			assert.NotEmpty(t, result.Packs, "pack list should not be empty")
			assert.Equal(t, orderQty, result.TotalItems-result.RemainingItems, "total items should equal order quantity plus remaining items")
			packsSum := 0
			for _, pack := range result.Packs {
				packsSum += pack.Size * pack.Count
			}
			assert.Equal(t, result.TotalItems, packsSum, "sum of pack sizes should equal total items")
			assert.Equal(t, packsSum-result.RemainingItems, orderQty, "remaining items should equal total items minus order quantity")
			assert.Less(t, result.RemainingItems, minimumPackSize, "remaining items should be less than the smallest pack size")
		})
	}
}

func TestCalculatePacks_DynamicRandomPackSizes(t *testing.T) {
	// Helper to get a crypto-safe random int between min and max (inclusive)
	randInt := func(min, max int64) int {
		nBig, err := crand.Int(crand.Reader, big.NewInt(max-min+1))
		if err != nil {
			t.Fatalf("failed to generate random number: %v", err)
		}
		return int(nBig.Int64() + min)
	}

	// Generate between 3 and 8 pack sizes in ascending order
	numPacks := randInt(3, 8)
	var packSizes []domain.Pack
	currentSize := randInt(50, 200) // start small

	for i := 0; i < numPacks; i++ {
		// Increase size by at least 1 and up to 10,000 per step
		increment := randInt(1, 2000)
		currentSize += increment
		packSizes = append(packSizes, domain.Pack{Size: currentSize})
	}

	minimumPackSize := packSizes[0].Size
	uc := packusecase.NewPackUseCase(&dynamicMockRepo{packs: packSizes})

	fmt.Println("Generated Pack Sizes:", packSizes)
	// Generate between 10 and 30 random orders for this pack set
	numOrders := randInt(5, 10)
	for i := 0; i < numOrders; i++ {
		orderQty := randInt(1, 50000) // random order quantity
		t.Run(fmt.Sprintf("Run_Random_Pack_Size_With_Order_%d", orderQty), func(t *testing.T) {
			result, err := uc.CalculatePacks(context.Background(), orderQty)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, result.TotalItems, orderQty, "total items must cover the order quantity")
			assert.Greater(t, result.TotalPacks, 0, "should have at least one pack")
			assert.NotEmpty(t, result.Packs, "pack list should not be empty")
			assert.Equal(t, orderQty, result.TotalItems-result.RemainingItems, "total items should equal order quantity plus remaining items")

			packsSum := 0
			for _, pack := range result.Packs {
				packsSum += pack.Size * pack.Count
			}
			assert.Equal(t, result.TotalItems, packsSum, "sum of pack sizes should equal total items")
			assert.Equal(t, packsSum-result.RemainingItems, orderQty, "remaining items should equal total items minus order quantity")
			assert.Less(t, result.RemainingItems, minimumPackSize, "remaining items should be less than the smallest pack size")
		})
	}

}

func TestCalculatePacks_RepoError(t *testing.T) {
	uc := packusecase.NewPackUseCase(&errorMockRepo{})

	result, err := uc.CalculatePacks(context.Background(), 10)
	assert.Error(t, err)
	assert.Equal(t, "use case failed to get packs: database connection failed", err.Error())
	assert.Equal(t, packusecase.CalculatePacksOutput{}, result)
}
