package models

type Vacancy struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" binding:"required" db:"name"`
	Salary     int    `json:"salary" binding:"required" db:"salary"`
	Experience string `json:"experience,omitempty" binding:"required" db:"experience"`
	City       string `json:"city,omitempty" binding:"required" db:"city"`
}
