package core

type User struct {
	ID          int    `json:"id,omitempty"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Class       string `json:"class"`
	LastLogin   string `json:"last_login,omitempty"`
	DisplayName string `json:"display_name"`
	Hash        string `json:"-"`
}
