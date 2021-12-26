package form

// PersonCreateParams create body params
type PersonCreateParams struct {
	LastName       string `json:"last_name"`        //
	FirstName      string `json:"first_name"`       //
	StateRegNumber string `json:"state_reg_number"` //rd
	IsActive       bool   `json:"is_active"`        // Идэвхтэй эсэх
	MobileNumber   string `json:"mobile_number"`    // Утас
}

// PersonUpdateParams update body params
type PersonUpdateParams struct {
	LastName       string `json:"last_name"`
	FirstName      string `json:"first_name"`
	StateRegNumber string `json:"state_reg_number"`
	IsActive       bool   `json:"is_active"`     // Идэвхтэй эсэх
	MobileNumber   string `json:"mobile_number"` // Утас
}

// PersonFilterCols sort hiih bolomjtoi column
type PersonFilterCols struct {
	LastName       string `json:"last_name"`
	FirstName      string `json:"first_name"`
	StateRegNumber string `json:"state_reg_number"`
	IsActive       string `json:"is_active"`     // Идэвхтэй эсэх
	MobileNumber   string `json:"mobile_number"` // Утас
}

// PersonFilter sort hiigdej boloh zuils
type PersonFilter struct {
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Sort   SortColumn       `json:"sort"`
	Filter PersonFilterCols `json:"filter"`
}
