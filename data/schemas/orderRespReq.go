package schemas

type (
	CreateOrderRequest struct {
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}

	CreateOrderResponse struct {
		ID string `json:"id"`
	}

	GetOrderResponse struct {
		ID         string  `json:"id"`
		UserID     int64   `json:"user_id"`
		ProductIDs []int64 `json:"product_ids"`
	}
)