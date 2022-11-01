package model

type Tag struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name" binding:"required"`
	Color string `json:"color" db:"color" binding:"required"`
}
