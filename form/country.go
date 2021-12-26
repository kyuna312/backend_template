package form

// CountryParams create body params
type CountryParams struct {
	Name string `json:"name" binding:"required"`
}

// CountryFilterCols sort hiih bolomjtoi column
type CountryFilterCols struct {
	Name string `json:"name"`
}

// CountryFilter sort hiigdej boloh zuils
type CountryFilter struct {
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Sort   SortColumn        `json:"sort"`
	Filter CountryFilterCols `json:"filter"`
}
