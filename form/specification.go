package form

// SpecificationCreateParams create body params
type SpecificationCreateParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// SpecificationUpdateParams update body params
type SpecificationUpdateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// SpecFilterCols sort hiih bolomjtoi column
type SpecFilterCols struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    string `json:"is_active"`
}

// SpecFilter sort hiigdej boloh zuils
type SpecFilter struct {
	Page   int             `json:"page"`
	Size   int             `json:"size"`
	Sort   SortColumn      `json:"sort"`
	Filter ClassFilterCols `json:"filter"`
}
