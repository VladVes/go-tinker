package fiber

import "errors"

// ------------------------------------Validation-------------------------------------------------------------
// Пример простой ручной валидации запроса по созданию поста
type CreatePostReq struct {
	UserID int64  `json:"user_id"`
	Text   string `json:"text"`
}

func (req *CreatePostReq) Validate() error {
	if req.UserID < 0 {
		return errors.New("user ID cannot by less than 0")
	}
	if req.Text == "" {
		return errors.New("text is empty")
	}
	if len(req.Text) > 140 {
		return errors.New("text is too long")
	}

	return nil
}
