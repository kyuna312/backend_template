package form

// GiftTypeParams create body params
type GiftTypeParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// GiftTypeFilterCols sort hiih bolomjtoi column
type GiftTypeFilterCols struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// GiftTypeFilter sort hiigdej boloh zuils
type GiftTypeFilter struct {
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Sort   SortColumn     `json:"sort"`
	Filter CityFilterCols `json:"filter"`
}
