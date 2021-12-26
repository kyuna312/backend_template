package form

// CustomerTypeCreateParams create body params
type CustomerTypeCreateParams struct {
	Name      string `json:"name" binding:"required"`
	ColorCode string `json:"color_code"`
	IsActive  bool   `json:"is_active"`
}

// CustomerTypeUpdateParams update body params
type CustomerTypeUpdateParams struct {
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	ColorCode string `json:"color_code"`
}

// CustomerTypeFilterCols sort hiih bolomjtoi column
type CustomerTypeFilterCols struct {
	TypeID   int    `json:"type_id"`
	Name     string `json:"name"`
	IsActive string `json:"is_active"`
}

// CustomerTypeFilter sort hiigdej boloh zuils
type CustomerTypeFilter struct {
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Sort   SortColumn             `json:"sort"`
	Filter CustomerTypeFilterCols `json:"filter"`
}
