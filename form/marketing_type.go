package form

// MarketingTypeParams create body params
type MarketingTypeParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// MarketingTypeFilterCols sort hiih bolomjtoi column
type MarketingTypeFilterCols struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// MarketingTypeFilter sort hiigdej boloh zuils
type MarketingTypeFilter struct {
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Sort   SortColumn     `json:"sort"`
	Filter CityFilterCols `json:"filter"`
}
