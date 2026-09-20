package dtos

type RightInfo struct {
	GroupName         string `json:"group_name"`
	GroupKey          string `json:"group_key"`
	GroupTitle        string `json:"group_title"`
	GroupDetails      string `json:"group_details"`
	GroupIcon         string `json:"group_icon"`
	GroupSortOrder    int    `json:"group_sort_order"`
	CategoryName      string `json:"category_name"`
	CategoryKey       string `json:"category_key"`
	CategoryURL       string `json:"category_url"`
	CategoryTitle     string `json:"category_title"`
	CategoryAction    string `json:"category_action"`
	CategoryDetails   string `json:"category_details"`
	CategoryIcon      string `json:"category_icon"`
	CategorySortOrder int    `json:"category_sort_order"`
	RightName         string `json:"right_name"`
	RightURL          string `json:"right_url"`
	RightTitle        string `json:"right_title"`
	RightAction       string `json:"right_action"`
	RightDetails      string `json:"right_details"`
	RightPopupStatus  int    `json:"right_popup_status"`
	RightWidth        int    `json:"right_width"`
	RightIcon         string `json:"right_icon"`
	RightSortOrder    int    `json:"right_sort_order"`
}
