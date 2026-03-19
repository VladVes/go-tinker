package schemas

type (
	CreateEmployeeRequest struct {
		Email string `json:"email"`
		Role string `json:"role"`
	}

	CreateEmployeeResponse struct {
		ID string `json:"id"`
	}
)