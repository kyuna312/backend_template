package form

import "time"

// IncomeParams update body params
type IncomeParams struct {
	OrderBookID        uint                   `json:"order_book_id" binding:"required"`   //
	WareHouseID        uint                   `json:"warehouse_id" binding:"required"`    // Агуулах ID
	PaymentTypeID      uint                   `json:"payment_type_id" binding:"required"` // Төлбөрийн төрлийн ID
	CompostionCost     float64                `json:"compostion_cost"`                    // Бүрдүүлэлтийн зардал
	TransportationCost float64                `json:"transportation_cost"`                // Тээвэрлэлтийн зардал
	CustomTaxPercent   float64                `json:"custom_tax_percent"`                 // Гааль 5%
	Description        string                 `json:"description"`                        // Тайлбар
	ValuteValue        float64                `json:"valute_value"`                       // Валют утга
	ValuteID           uint                   `json:"valute_id"`                          // Валют ID
	IsVat              bool                   `json:"is_vat"`                             // Татвартай эсэх
	IsDiscount         bool                   `json:"is_discount"`                        // Хямдарлтай эсэх
	Items              []DetailIncomeDtlParam `json:"items"`
	PercentVat         float64                `json:"percent_vat"`
	// CustomerID         uint                   `json:"customer_id" binding:"required"`     // Харилцагчийн ID
	// IsDiscount         bool                   `json:"is_discount"`                        // Хөнгөлөлттэй эсэх
}

// DetailIncomeDtlParam ...
type DetailIncomeDtlParam struct {
	ID              uint      `json:"id"`                             // OrderBookDtlID
	ItemID          uint      `json:"item_id" binding:"required"`     // Барааны ID
	OrderQty        uint      `json:"order_qty" binding:"required"`   // орлого тоо
	UnitPrice       float64   `json:"unit_price" binding:"required"`  // Нэгжийн үнэ
	PercentDiscount uint      `json:"percent_discount"`               // Хөнгөлөлт/хувь/
	ExpireDate      time.Time `json:"expire_date" binding:"required"` //
	// PercentVat      float64   `json:"percent_vat" binding:"required"`   // НӨАТ-н хувь
	SerialNumber string `json:"serial_number" binding:"required"` //
}

// UpdateIncomeStatus ...
type UpdateIncomeStatus struct {
	IncomeID    int    `json:"income_id" binding:"required"`   //
	StatusID    int    `json:"status_id" binding:"required"`   //
	Description string `json:"description" binding:"required"` //
}

// IncomeStatusChangeParams update body params
type IncomeStatusChangeParams struct {
	IncomeID    uint            `json:"income_id" binding:"required"` //
	StatusID    uint            `json:"status_id" binding:"required"` //
	Description string          `json:"description"`                  // Тайлбар
	Items       []CheckDtlParam `json:"items"`
}

// CheckDtlParam ...
type CheckDtlParam struct {
	InComeItemDtlID uint      `json:"in_come_item_dtl_id" `             //
	ItemID          uint      `json:"item_id"`                          //
	ExpireDate      time.Time `json:"expire_date" binding:"required"`   //
	SerialNumber    string    `json:"serial_number" binding:"required"` //
	Price           float64   `json:"price"`                            //
}

// IncomeFilterCols sort hiih bolomjtoi column
type IncomeFilterCols struct {
	Name                string    `json:"name"`              //
	Description         string    `json:"description"`       // Тайлбар
	OrderNumber         string    `json:"order_number"`      // Захиалгын дугаар
	CustomerID          int       `json:"customer_id" `      // Харилцагчийн ID
	OrderTypeID         int       `json:"order_type_id"`     // Захиалгын төрөл /in, out/
	OrderSubTypeID      int       `json:"order_sub_type_id"` // Захиалгын дэд төрөл /Дотоод, Импортын, Солилцооны/
	OrderStatusID       int       `json:"order_status_id"`   // Захиалгын төлөв
	PaymentTypeID       int       `json:"payment_type_id"`   // Төлбөрийн төрлийн ID
	Discount            int       `json:"discount"`          // Хөнгөлөлтийн дүн
	Total               int       `json:"total"`             // Захиалгын нийт дүн
	IsVat               bool      `json:"is_vat"`            // Татвартай эсэх
	IsPaid              bool      `json:"is_paid"`           // Төлбөр хийгдсэн эсэх
	IsDiscount          bool      `json:"is_discount"`       // Хөнгөлөлттэй эсэх
	ValuteID            int       `json:"valute_id"`         // Валют ID
	ValuteValue         int       `json:"valute_value"`      // Валют утга
	ExternalStartDate   time.Time `json:"external_start_date"`
	ExternalEndDate     time.Time `json:"external_end_date"`
	ExternalWarehouseID int       `json:"external_warehouse_id"`
}

// FindItem ...
type FindItem struct {
	ItemName string `json:"item_name"`
	Barcode  string `json:"barcode"`
}

// IncomeFilter sort hiigdej boloh zuils
type IncomeFilter struct {
	Page   int              `json:"page"`
	Size   int              `json:"size"`
	Sort   SortColumn       `json:"sort"`
	Filter IncomeFilterCols `json:"filter"`
}
