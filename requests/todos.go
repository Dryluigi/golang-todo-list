package requests

type TodoCreate struct {
	Todo        string `json:"todo"`
	Description string `json:"description"`
	IsComplete  bool   `json:"is_complete"`
}

type TodoUpdate struct {
	Todo        string `json:"todo"`
	Description string `json:"description"`
	IsComplete  bool   `json:"is_complete"`
}
