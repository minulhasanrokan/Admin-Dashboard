package dtos

type RightCategoryDTO struct {
	Id         uint64   `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	Action     string `json:"action"`
	Icon       string `json:"icon"`
	ShortOrder int    `json:"short_order"`
}

type RightGroupTreeDTO struct {
	Id         uint64               `json:"id"`
	Name       string             `json:"name"`
	Key        string             `json:"key"`
	Title      string             `json:"title"`
	Icon       string             `json:"icon"`
	ShortOrder int                `json:"short_order"`
	Categories []RightCategoryDTO `json:"categories"`
}
