package form

// DistrictParams create body params
type DistrictParams struct {
	Name   string `json:"name" binding:"required"`
	CityID int    `json:"city_id" binding:"required"`
}

// DistrictFilterCols sort hiih bolomjtoi column
type DistrictFilterCols struct {
	Name   string `json:"name"`
	CityID int    `json:"city_id"`
}

// DistrictFilter sort hiigdej boloh zuils
type DistrictFilter struct {
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Sort   SortColumn         `json:"sort"`
	Filter DistrictFilterCols `json:"filter"`
}
