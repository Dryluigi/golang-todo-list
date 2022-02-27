package requests

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
