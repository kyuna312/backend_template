package form

// OrderBookDtlParams update body params
type OrderBookDtlParams struct {
	OrderBookID uint                  `json:"order_book_id"` // Захиалгын ID
	Items       []DetailOrderDtlParam `json:"items"`
}

// DetailOrderDtlParam ...
type DetailOrderDtlParam struct {
	ID                uint    `json:"id"`                  //
	ItemID            uint    `json:"item_id"`             // Барааны ID
	OrderQty          uint    `json:"order_qty"`           // Захиалгын тоо
	UnitPrice         float64 `json:"unit_price"`          // Нэгжийн үнэ
	PercentDiscount   float64 `json:"percent_discount"`    // Хөнгөлөлт/хувь/
	UnitDiscount      float64 `json:"unit_discount"`       // Нэгжийн хөнгөлсөн дүн
	ItemTotalDiscount float64 `json:"item_total_discount"` // Тухайн барааны нийт хөнгөлсөн дүн
	PercentVat        float64 `json:"percent_vat"`         // НӨАТ-н хувь
	UnitVat           float64 `json:"unit_vat"`            // Нэгж НӨАТ
	ItemTotalVat      float64 `json:"item_total_vat"`      // Тухайн барааны нийт НӨАТ
	ItemTotalAmount   float64 `json:"item_total_amount"`   // Тухайн барааны нийт дүн
	IsRemoved         bool    `json:"is_removed"`          // Буцаагдсан эсэх
}

// OrderBookDtlFilterCols sort hiih bolomjtoi column
type OrderBookDtlFilterCols struct {
	OrderBookID       int  `json:"order_book_id"`       // Захиалгын ID
	ItemID            int  `json:"item_id"`             // Барааны ID
	OrderQty          int  `json:"order_qty"`           // Захиалгын тоо
	UnitPrice         int  `json:"unit_price"`          // Нэгжийн үнэ
	PercentDiscount   int  `json:"percent_discount"`    // Хөнгөлөлт/хувь/
	UnitDiscount      int  `json:"unit_discount"`       // Нэгжийн хөнгөлсөн дүн
	ItemTotalDiscount int  `json:"item_total_discount"` // Тухайн барааны нийт хөнгөлсөн дүн
	PercentVat        int  `json:"percent_vat"`         // НӨАТ-н хувь
	UnitVat           int  `json:"unit_vat"`            // Нэгж НӨАТ
	ItemTotalVat      int  `json:"item_total_vat"`      // Тухайн барааны нийт НӨАТ
	ItemTotalAmount   int  `json:"item_total_amount"`   // Тухайн барааны нийт дүн
	IsRemoved         bool `json:"is_removed"`          // Буцаагдсан эсэх
	IsVat             bool `json:"is_vat"`              // Татвартай эсэх
	IsDiscount        bool `json:"is_discount"`         // Хөнгөлөлттэй эсэх
}

// OrderBookDtlFilter sort hiigdej boloh zuils
type OrderBookDtlFilter struct {
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Sort   SortColumn             `json:"sort"`
	Filter OrderBookDtlFilterCols `json:"filter"`
}
