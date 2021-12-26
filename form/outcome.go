package form

import "time"

// OutcomeParams update body params
type OutcomeParams struct {
	OrderBookID   uint                    `json:"order_book_id" binding:"required"`   //
	Description   string                  `json:"description"`                        // Тайлбар
	WarehouseID   uint                    `json:"warehouse_id" binding:"required"`    // Агуулах ID
	CustomerID    uint                    `json:"customer_id" binding:"required"`     // Харилцагчийн ID
	PaymentTypeID uint                    `json:"payment_type_id" binding:"required"` // Төлбөрийн төрлийн ID
	ValuteValue   float64                 `json:"valute_value"`                       // Валют утга
	ValuteID      uint                    `json:"valute_id"`                          // Валют ID
	IsVat         bool                    `json:"is_vat"`                             // Татвартай эсэх
	VatPercent    float64                 `json:"vat_percent"`                        // НӨАТ-н хувь
	IsDiscount    bool                    `json:"is_discount"`                        // Хөнгөлөлттэй эсэх
	IsDelivery    bool                    `json:"is_delivery"`                        // Хөнгөлөлттэй эсэх
	DeliveryDate  time.Time               `json:"delivery_date"`
	Items         []DetailOutcomeDtlParam `json:"items"`
}

// OutcomeStatusChangeParams update body params
type OutcomeStatusChangeParams struct {
	OrderBookID uint   `json:"order_book_id" binding:"required"`
	StatusID    uint   `json:"status_id" binding:"required"` //
	Description string `json:"description"`                  // Тайлбар
	// OutcomeID   uint   `json:"outcome_id" binding:"required"` //
	// Items       []OutcomeItems `json:"items"`
}

// OutcomeFindItem search
type OutcomeFindItem struct {
	ItemName    string `json:"item_name"` //
	Barcode     string `json:"barcode"`   // barcode
	WarehouseID int    `json:"warehouse_id" binding:"required"`
	CustomerID  int    `json:"customer_id"`
	PriceTypeID int    `json:"price_type_id"`
}

// DetailOutcomeDtlParam ...
type DetailOutcomeDtlParam struct {
	ID                 uint      `json:"id"`                               // OrderBookDtlID
	WareHouseID        uint      `json:"warehouse_id" binding:"required"`  // Агуулах ID
	ItemID             uint      `json:"item_id" binding:"required"`       // Барааны ID
	OutcomeQty         uint      `json:"income_qty" binding:"required"`    // орлого тоо
	UnitPrice          float64   `json:"unit_price" binding:"required"`    // Нэгжийн үнэ
	PercentDiscount    uint      `json:"percent_discount"`                 // Хөнгөлөлт/хувь/
	UnitDiscount       uint      `json:"unit_discount"`                    // Нэгжийн хөнгөлсөн дүн
	ItemTotalDiscount  uint      `json:"item_total_discount"`              // Тухайн барааны нийт хөнгөлсөн дүн
	UnitVat            float64   `json:"unit_vat"`                         // Нэгж НӨАТ
	ItemTotalVat       float64   `json:"item_total_vat"`                   // Тухайн барааны нийт НӨАТ
	ItemTotalAmount    float64   `json:"item_total_amount"`                // Тухайн барааны нийт дүн
	ExpireDate         time.Time `json:"expire_date" binding:"required"`   //
	SerialNumber       string    `json:"serial_number" binding:"required"` //
	TransportationCost float64   `json:"transportation_cost"`              //
	CustomTax          float64   `json:"custom_tax"`                       //
	CompostionCost     float64   `json:"compostion_cost"`                  //
	Price              float64   `json:"price"`                            //
}

// OutcomeFilterCols sort hiih bolomjtoi column
type OutcomeFilterCols struct {
	Name           string `json:"name"`              //
	Description    string `json:"description"`       // Тайлбар
	OrderNumber    string `json:"order_number"`      // Захиалгын дугаар
	CustomerID     int    `json:"customer_id" `      // Харилцагчийн ID
	OrderTypeID    int    `json:"order_type_id"`     // Захиалгын төрөл /in, out/
	OrderSubTypeID int    `json:"order_sub_type_id"` // Захиалгын дэд төрөл /Дотоод, Импортын, Солилцооны/
	OrderStatusID  int    `json:"order_status_id"`   // Захиалгын төлөв
	PaymentTypeID  int    `json:"payment_type_id"`   // Төлбөрийн төрлийн ID
	Discount       int    `json:"discount"`          // Хөнгөлөлтийн дүн
	Total          int    `json:"total"`             // Захиалгын нийт дүн
	IsVat          bool   `json:"is_vat"`            // Татвартай эсэх
	IsPaid         bool   `json:"is_paid"`           // Төлбөр хийгдсэн эсэх
	IsDiscount     bool   `json:"is_discount"`       // Хөнгөлөлттэй эсэх
	ValuteID       int    `json:"valute_id"`         // Валют ID
	ValuteValue    int    `json:"valute_value"`      // Валют утга
}

// OutcomeFilter sort hiigdej boloh zuils
type OutcomeFilter struct {
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Sort   SortColumn        `json:"sort"`
	Filter OutcomeFilterCols `json:"filter"`
}
