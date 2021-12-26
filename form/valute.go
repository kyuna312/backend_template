package form

// ValuteCreateParams create body params
type ValuteCreateParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Symbol      string `json:"symbol" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

// ValuteUpdateParams update body params
type ValuteUpdateParams struct {
	Name        string `json:"name" binding:"required"`
	Symbol      string `json:"symbol" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// ValuteFilterCols sort hiih bolomjtoi column
type ValuteFilterCols struct {
	TypeID      int    `json:"type_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    string `json:"is_active"`
}

// ValuteFilter sort hiigdej boloh zuils
type ValuteFilter struct {
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Sort   SortColumn       `json:"sort"`
	Filter ValuteFilterCols `json:"filter"`
}
