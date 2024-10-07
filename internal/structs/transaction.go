package structs

type TransactionCompletePayload struct {
	TransactionType                string                         `json:"transaction_type" validate:"required"`
	BatchID                        string                         `json:"batch_id" validate:"required"`
	TransactionID                  string                         `json:"transaction_id" validate:"required"`
	Medium                         string                         `json:"medium" validate:"required"`
	Amount                         string                         `json:"amount" validate:"required"`
	InitiatedPersona               string                         `json:"initiated_persona" validate:"required"`
	SenderWalletNumber             string                         `json:"sender_wallet_number" validate:"required"`
	SenderWalletID                 int64                          `json:"sender_wallet_id" validate:"required"`
	ReceiverWalletNumber           string                         `json:"receiver_wallet_number" validate:"required"`
	ReceiverWalletID               int64                          `json:"receiver_wallet_id"`
	TransactionTypeSpecificDetails TransactionTypeSpecificDetails `json:"transaction_type_specific_details" validate:"required"`
}
type TransactionTypeSpecificDetails struct {
	Mno                string `json:"mno" validate:"required"`
	Mid                string `json:"mid" validate:"required"`
	ParentWalletNumber string `json:"parent_wallet_number"`
	MerchantTerminal   string `json:"merchant_terminal"`
}
