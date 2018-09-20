package models

type UserRole string

const (
	Viewer UserRole = "viewer"
	Editor          = "editor"
)

type User struct {
	ID       int      `json:"id" db:"id"`
	Role     UserRole `json:"role" db:"role"`
	Login    string   `json:"login" binding:"required" db:"login"`
	Password string   `json:"password,omitempty" binding:"required" db:"password"`
}

type ChangeUserRole struct {
	Login   string   `json:"login" binding:"required" db:"login"`
	NewRole UserRole `json:"role" binding:"required" db:"role"`
}
