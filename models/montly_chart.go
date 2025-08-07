package models

type MonthlyCategoryItem struct {
	Category2 string  `json:"category2"`
	Total     float64 `json:"total"`
}

type MonthlyCategoryGroup struct {
	Month      string                `json:"month"`
	Categories []MonthlyCategoryItem `json:"categories"`
}

type ResponseWithMonths struct {
	Months []MonthlyCategoryGroup `json:"months"`
}
