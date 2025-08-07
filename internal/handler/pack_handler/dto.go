package pack_handler

type CalculatePacksReq struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
