package form

// MeasureCreateParams create body params
type MeasureCreateParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Code        string `json:"code"`
	IsActive    bool   `json:"is_active"`
}

// MeasureUpdateParams update body params
type MeasureUpdateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	IsActive    bool   `json:"is_active"`
}

// MeasureFilterCols sort hiih bolomjtoi column
type MeasureFilterCols struct {
	TypeID      int    `json:"type_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    string `json:"is_active"`
}

// MeasureFilter sort hiigdej boloh zuils
type MeasureFilter struct {
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Sort   SortColumn        `json:"sort"`
	Filter MeasureFilterCols `json:"filter"`
}
