package form

import "time"

// SystemUserParams create body params
type SystemUserParams struct {
	IsActive  bool      `json:"is_active"`  // Идэвхтэй эсэх
	Username  string    `json:"username"`   // Нэр
	StartDate time.Time `json:"start_date"` // Эхлэх огноо
	EndDate   time.Time `json:"end_date"`   // Дуусах огноо
	Password  string    `json:"password"`   // Нууц үг
	PersondID uint      `json:"persond_id"` //
}

// SystemUserFilterCols sort hiih bolomjtoi column
type SystemUserFilterCols struct {
	LastName       string `json:"last_name"`
	FirstName      string `json:"first_name"`
	StateRegNumber string `json:"state_reg_number"`
	IsActive       bool   `json:"is_active"`     // Идэвхтэй эсэх
	MobileNumber   string `json:"mobile_number"` // Утас
	Username       string `json:"username"`      // Нэр
}

// SystemUserFilter sort hiigdej boloh zuils
type SystemUserFilter struct {
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Sort   SortColumn           `json:"sort"`
	Filter SystemUserFilterCols `json:"filter"`
}
