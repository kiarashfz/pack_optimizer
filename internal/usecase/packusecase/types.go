package packusecase

type Pack struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

type CalculatePacksOutput struct {
	TotalItems     int    `json:"total_items"`     // Total items that fit in the packs
	RemainingItems int    `json:"remaining_items"` // Number of empty spaces in packs
	TotalPacks     int    `json:"total_packs"`     // Total number of packs used
	Packs          []Pack `json:"packs"`           // Calculated packs with their sizes and counts
}
