package form

// CityParams create body params
type CityParams struct {
	Name      string `json:"name" binding:"required"`
	CountryID uint   `json:"country_id" binding:"required"`
}

// CityFilterCols sort hiih bolomjtoi column
type CityFilterCols struct {
	Name      string `json:"name"`
	CountryID int    `json:"country_id"`
}

// CityFilter sort hiigdej boloh zuils
type CityFilter struct {
	Page   int            `json:"page"`
	Size   int            `json:"size"`
	Sort   SortColumn     `json:"sort"`
	Filter CityFilterCols `json:"filter"`
}
