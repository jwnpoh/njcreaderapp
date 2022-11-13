package core

type User struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	LastLogin string `json:"lastLogin"`
	Name      string `json:"name"`
	ID        int    `json:"-"`
	Hash      string `json:"-"`
}
