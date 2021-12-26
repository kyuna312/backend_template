package form

// PositionCreateParams create body params
type PositionCreateParams struct {
	Name        string `json:"name" binding:"required"`
	TypeID      int    `json:"type_id" binding:"required"`
	Description string `json:"description"`
	Code        string `json:"code"`
	IsActive    bool   `json:"is_active"`
}

// PositionUpdateParams update body params
type PositionUpdateParams struct {
	Name        string `json:"name"`
	TypeID      int    `json:"type_id"`
	Description string `json:"description"`
	Code        string `json:"code"`
	IsActive    bool   `json:"is_active"`
}

// PositionFilterCols sort hiih bolomjtoi column
type PositionFilterCols struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsActive       string `json:"is_active"`
	PositionTypeID int    `json:"position_type_id"`
}

// PositionFilter sort hiigdej boloh zuils
type PositionFilter struct {
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Sort   SortColumn         `json:"sort"`
	Filter PositionFilterCols `json:"filter"`
}
