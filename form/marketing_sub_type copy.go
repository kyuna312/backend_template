package form

// SubTypeParams create body params
type SubTypeParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// SubTypeFilterCols sort hiih bolomjtoi column
type SubTypeFilterCols struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// SubTypeFilter sort hiigdej boloh zuils
type SubTypeFilter struct {
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Sort   SortColumn     `json:"sort"`
	Filter CityFilterCols `json:"filter"`
}
