package model

type Tag struct {
	ID    int    `json:"id" db:"id"`
	Label string `json:"label" db:"label" binding:"required"`
	Color string `json:"color" db:"color" binding:"required"`
}
