package form

// ClassificationCreateParams create body params
type ClassificationCreateParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// ClassificationUpdateParams update body params
type ClassificationUpdateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// ClassFilterCols sort hiih bolomjtoi column
type ClassFilterCols struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    string `json:"is_active"`
}

// ClassFilter sort hiigdej boloh zuils
type ClassFilter struct {
	Page   int             `json:"page"`
	Size   int             `json:"size"`
	Sort   SortColumn      `json:"sort"`
	Filter ClassFilterCols `json:"filter"`
}
