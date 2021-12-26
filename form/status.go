package form

// StatusParams create body params
type StatusParams struct {
	Name        string `json:"name" binding:"required"`
	ColorCode   string `json:"color_code"`
	TypeID      int    `json:"type_id" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// StatusFilterCols sort hiih bolomjtoi column
type StatusFilterCols struct {
	StatusTypeID int    `json:"status_type_id"`
	Name         string `json:"name"`
	ColorCode    string `json:"color_code"`
	Description  string `json:"description"`
	IsActive     string `json:"is_active"`
}

// StatusFilter sort hiigdej boloh zuils
type StatusFilter struct {
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Sort   SortColumn       `json:"sort"`
	Filter StatusFilterCols `json:"filter"`
}
