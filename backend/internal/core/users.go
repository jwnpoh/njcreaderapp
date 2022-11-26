package core

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Class       string `json:"class"`
	LastLogin   string `json:"last_login"`
	DisplayName string `json:"display_name"`
	Hash        string `json:"-"`
}
