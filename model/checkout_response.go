package model

type CheckoutResponse struct {
	Items []CheckoutItem
	Total int64
}
