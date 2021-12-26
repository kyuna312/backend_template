package form

// PaymentTypeParams create body params
type PaymentTypeParams struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	IsActive         bool   `json:"is_active"`
	PaymentDay       int    `json:"payment_day"`       // Төлбөр төлөх хугацаа
	PrepaidPercent   int    `json:"prepaid_percent"`   // Урьдчилж төлөх хувь
	PaymentCondition string `json:"payment_condition"` // Төлбөр төлөх нөхцөл
}

// PaymentTypeFilterCols sort hiih bolomjtoi column
type PaymentTypeFilterCols struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsActive         string `json:"is_active"`
	PaymentDay       int    `json:"payment_day"`       // Төлбөр төлөх хугацаа
	PrepaidPercent   int    `json:"prepaid_percent"`   // Урьдчилж төлөх хувь
	PaymentCondition string `json:"payment_condition"` // Төлбөр төлөх нөхцөл
}

// PaymentTypeFilter sort hiigdej boloh zuils
type PaymentTypeFilter struct {
	Page   int                   `json:"page"`
	Size   int                   `json:"size"`
	Sort   SortColumn            `json:"sort"`
	Filter PaymentTypeFilterCols `json:"filter"`
}
