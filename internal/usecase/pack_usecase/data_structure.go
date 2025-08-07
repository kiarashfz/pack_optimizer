package pack_usecase

// A node represents a state: total items, total packs, pack count
type Node struct {
	totalItems int   // total items shipped so far
	totalPacks int   // number of packs used
	packCount  []int // count of each pack size used
	index      int   // index for heap interface
}

// Priority: first by totalItems, then by totalPacks
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].totalItems != pq[j].totalItems {
		return pq[i].totalItems < pq[j].totalItems
	}
	return pq[i].totalPacks < pq[j].totalPacks
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	node := x.(*Node)
	node.index = len(*pq)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.index = -1
	*pq = old[:n-1]
	return node
}
