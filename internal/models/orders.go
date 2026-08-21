package models

type UpdateOrderRequest struct {
	OrderName        string  `json:"orderName" binding:"required"`
	OrderID          int     `json:"orderId" binding:"required"`
	ValidateStatus   string  `json:"validateStatus"`
	Amount           float64 `json:"amount"`
	LogisticValidate bool    `json:"logisticValidate"`
	PaymentMethodID  int     `json:"paymentMethodId"`
	RequiresChange   bool    `json:"requiresChange"`
}

type ChangePaidProcessRequest struct {
	OrderID int     `json:"orderId" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
}

// CashReturnData representa los datos necesarios para la devolución en efectivo
type CashReturnData struct {
	Bank    string `json:"bank"`
	Phone   string `json:"phone"`
	DNI     string `json:"dni"`
	DNIType string `json:"dniType"`
}

type BCVTasaUSDResponse struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}
