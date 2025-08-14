// Package packhandler defines request and response structures for pack-related operations.
package packhandler

type CalculatePacksReq struct {
	Quantity int `json:"quantity" validate:"required,min=1,max=99999999"` // Quantity must be between 1 and 99,999,999
}
