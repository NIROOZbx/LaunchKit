package dtos

type NonceRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required,hexadecimal,len=42"`
}

type VerifyRequest struct {
	Message   string `json:"message" validate:"required"`
	Signature string `json:"signature" validate:"required"`
	UserType  string  `json:"user_type"`
}