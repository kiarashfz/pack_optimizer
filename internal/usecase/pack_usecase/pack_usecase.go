package pack_usecase

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"pack_optimizer/internal/domain"
	"sort"
)

// PackUseCase is a use case that provides methods to calculate the optimal pack combinations.
type PackUseCase struct {
	packRepo domain.PackRepository // packRepo is the repository interface for accessing pack data.
}

// NewPackUseCase creates a new instance of PackUseCase.
// Parameters:
//   - packRepo: An implementation of the domain.PackRepository interface.
//
// Returns:
//   - A pointer to a new PackUseCase instance.
func NewPackUseCase(packRepo domain.PackRepository) *PackUseCase {
	return &PackUseCase{packRepo: packRepo}
}

// CalculatePacks calculates the optimal combination of packs to fulfill an order quantity.
// Parameters:
//   - ctx: The context for managing request-scoped values, deadlines, and cancellation signals.
//   - orderQty: The quantity of items to fulfill in the order.
//
// Returns:
//   - A CalculatePacksOutput struct containing the details of the calculated packs.
//   - An error if the operation fails or the input is invalid.
func (uc *PackUseCase) CalculatePacks(ctx context.Context, orderQty int) (CalculatePacksOutput, error) {
	if orderQty <= 0 {
		return CalculatePacksOutput{}, errors.New("order quantity must be greater than 0")
	}

	packs, err := uc.packRepo.GetAllPacks(ctx)
	if err != nil {
		return CalculatePacksOutput{}, err
	}

	packSizes := make([]int, 0, len(packs))
	for _, p := range packs {
		packSizes = append(packSizes, p.Size)
	}

	res, err := findBestPackCombination(orderQty, packSizes)
	if err != nil {
		return CalculatePacksOutput{}, err
	}

	var packDetails []Pack
	for i, count := range res.packCount {
		if count > 0 {
			packDetails = append(packDetails, Pack{
				Size:  packSizes[i],
				Count: count,
			})
		}
	}
	output := CalculatePacksOutput{
		TotalItems:     res.totalItems,
		RemainingItems: res.totalItems - orderQty,
		TotalPacks:     res.totalPacks,
		Packs:          packDetails,
	}

	return output, nil
}

// findBestPackCombination finds the optimal combination of packs to fulfill the order quantity.
// Parameters:
//   - order: The quantity of items to fulfill in the order.
//   - packSizes: A slice of integers representing the available pack sizes.
//
// Returns:
//   - A pointer to a Node struct representing the optimal combination of packs.
func findBestPackCombination(order int, packSizes []int) (*Node, error) {
	sort.Ints(packSizes) // Sort pack sizes in ascending order for easier index mapping.
	maxPack := packSizes[len(packSizes)-1]

	// Visited map to avoid revisiting the same total with the same pack count.
	visited := make(map[string]bool)

	initial := &Node{
		totalItems: 0,
		totalPacks: 0,
		packCount:  make([]int, len(packSizes)),
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, initial)

	for pq.Len() > 0 {
		curr, ok := heap.Pop(pq).(*Node)
		if !ok {
			return nil, errors.New("failed to pop from priority queue")
		}

		// If the current state satisfies the order, it's optimal.
		if curr.totalItems >= order {
			return curr, nil
		}

		for i, size := range packSizes {
			nextTotal := curr.totalItems + size
			key := fmt.Sprintf("%d:%d", nextTotal, i) // Unique key for visited map.

			if nextTotal > order+maxPack || visited[key] {
				continue
			}
			visited[key] = true

			newPackCount := append([]int(nil), curr.packCount...)
			newPackCount[i]++

			next := &Node{
				totalItems: nextTotal,
				totalPacks: curr.totalPacks + 1,
				packCount:  newPackCount,
			}

			heap.Push(pq, next)
		}
	}

	return nil, nil // Return nil if no combination is found
}
