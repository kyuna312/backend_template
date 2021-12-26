package form

// MarketingItemsParams create body params
type MarketingItemsParams struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price"`
	Count       int     `json:"count"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
}

// MarketingItemsFilterCols sort hiih bolomjtoi column
type MarketingItemsFilterCols struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// MarketingItemsFilter sort hiigdej boloh zuils
type MarketingItemsFilter struct {
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Sort   SortColumn     `json:"sort"`
	Filter CityFilterCols `json:"filter"`
}
