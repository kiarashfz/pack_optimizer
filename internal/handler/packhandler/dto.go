// Package packhandler defines request and response structures for pack-related operations.
package packhandler

type CalculatePacksReq struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
