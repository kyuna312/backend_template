package form

// ItemParams body params
type ItemParams struct {
	TargetCount      uint   `json:"target_count" binding:"required"` // Агуулахад байх шаардлагатай тоо
	BarCode          uint   `json:"barcode" binding:"required"`      // Бар код**
	Name             string `json:"name" binding:"required"`         // Барааны нэр
	InterName        string `json:"inter_name"`                      // Олон улсын нэршил
	MedDose          string `json:"med_dose"`
	ClassificationID uint   `json:"classification_id" binding:"required"`
	SpecificationID  uint   `json:"specification_id" binding:"required"`
	TermIssuanceID   uint   `json:"term_issuance_id" binding:"required"`
	ManuName         string `json:"manu_name" binding:"required"` // Үйлдвэрлэгчийн нэр
	CustomerID       uint   `json:"customer_id"`
	CountryID        uint   `json:"country_id" binding:"required"`
	Packaging        string `json:"packaging" binding:"required"` // Савлагаа
	MeasureID        uint   `json:"measure_id" binding:"required"`
	TypeID           uint   `json:"type_id" binding:"required"`
	Description      string `json:"description"` // Тайлбар
	IsActive         bool   `json:"is_active"`   // Идэвхитэй эсэх
	IsImport         bool   `json:"is_import"`   // Импортын бараа эсэх
}

// FilterItem sort hiih bolomjtoi column
type FilterItem struct {
	BarCode          string `json:"barcode"` // Бар код**
	Code             string `json:"code"`
	Name             string `json:"name"`       // Барааны нэр
	InterName        string `json:"inter_name"` // Олон улсын нэршил
	MedDose          string `json:"med_dose"`
	ClassificationID int    `json:"classification_id"`
	SpecificationID  int    `json:"specification_id"`
	ManuName         string `json:"manu_name" ` // Үйлдвэрлэгчийн нэр
	CustomerID       int    `json:"customer_id"`
	CountryID        int    `json:"country_id" `
	Packaging        string `json:"packaging" ` // Савлагаа
	MeasureID        int    `json:"measure_id" `
	TypeID           int    `json:"type_id"`
	Description      string `json:"description"` // Тайлбар
	IsImport         string `json:"is_import"`   // Импортын бараа эсэх
	IsActive         string `json:"is_active"`
	CreatedDate      string `json:"created_date"`
}

// SortItems sort hiigdej boloh zuils
type SortItems struct {
	Page   int        `json:"page"`
	Size   int        `json:"size"`
	Sort   SortColumn `json:"sort"`
	Filter FilterItem `json:"filter"`
}
