package core

type User struct {
	ID        int    `json:"-"`
	Email     string `json:"email"`
	Hash      string `json:"password"`
	Role      string `json:"role"`
	LastLogin string `json:"lastLogin"`
}
