package tag

type CreateTagDTO struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type UpdateTagDTO struct {
	Name  string `json:"name" db:"name"`
	Color string `json:"color" db:"color"`
}

type GetAllTagsDTO struct {
	Tags []Tag `json:"tags"`
}
