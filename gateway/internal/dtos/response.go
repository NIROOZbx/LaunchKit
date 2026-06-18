package dtos

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
