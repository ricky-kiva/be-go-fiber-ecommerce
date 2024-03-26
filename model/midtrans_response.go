package model

type MidtransResponse struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}
