package form

// StreetParams create body params
type StreetParams struct {
	Name       string `json:"name" binding:"required"`
	DistrictID int    `json:"district_id" binding:"required"`
}

// StreetFilterCols sort hiih bolomjtoi column
type StreetFilterCols struct {
	Name       string `json:"name"`
	DistrictID int    `json:"district_id"`
}

// StreetFilter sort hiigdej boloh zuils
type StreetFilter struct {
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Sort   SortColumn       `json:"sort"`
	Filter StreetFilterCols `json:"filter"`
}
