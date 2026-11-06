package dtos

type RightCategoryDTO struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	Action     string `json:"action"`
	Icon       string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

type RightGroupTreeDTO struct {
	Id         uint64             `json:"id"`
	Name       string             `json:"name"`
	Key        string             `json:"key"`
	Title      string             `json:"title"`
	Icon       string             `json:"icon"`
	SortOrder int                `json:"sort_order"`
	Categories []RightCategoryDTO `json:"categories"`
}

type GridRightItem struct {
	Name        string `json:"name"`
	ActionName  string `json:"action_name"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Action      string `json:"action"`
	PopupStatus int    `json:"popup_status"`
	Width       int    `json:"width"`
	Icon        string `json:"icon"`
	SortOrder  int    `json:"sort_order"`
	ShortName    string `json:"short_name"`
}
