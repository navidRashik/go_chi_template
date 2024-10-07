package structs

type MerchantStatusUpdatePayload struct {
	ActiveMerchant   []string `json:"active_merchant" validate:"required"`
	InactiveMerchant []string `json:"inactive_merchant" validate:"required"`
}
