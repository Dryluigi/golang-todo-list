package models

type Todo struct {
	ID          uint   `json:"id"`
	Todo        string `json:"todo"`
	Description string `json:"description"`
	IsComplete  uint8  `json:"is_complete"`
	UserId      uint   `json:"user_id"`
}
