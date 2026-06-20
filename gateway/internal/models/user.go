package models

type NonceRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required,hexadecimal,len=42"`
}

type VerifyRequest struct {
	Message   string `json:"message" validate:"required"`
	Signature string `json:"signature" validate:"required"`
	UserType  string  `json:"user_type"`
}


type NonceResponse struct {
	Nonce    string `json:"nonce"`
	Message  string `json:"message"` 
}

type UserResponse struct {
	ID            string `json:"id"`
	WalletAddress string `json:"wallet_address"`
	EnsName       string `json:"ens_name,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	Role          string `json:"role"`
	IsOnboarded   bool   `json:"is_onboarded"`
}
