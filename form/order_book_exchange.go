package form

import "time"

// OrderBookExchangeParams update body params
type OrderBookExchangeParams struct {
	PriceTypeID      int                              `json:"price_type_id" binding:"required"`
	CustomerID       int                              `json:"customer_id" binding:"required"` // Харилцагчийн ID
	Description      string                           `json:"description"`                    // Тайлбар
	IsVat            bool                             `json:"is_vat"`                         // Татвартай эсэх\
	DeliveryPersonID int                              `json:"delivery_person_id"`
	WarehouseID      int                              `json:"warehouse_id"`
	OutItems         []OrderBookExchangeOutItemParams `json:"out_items"`
	InItems          []OrderBookExchangeInItemParams  `json:"in_items"`
}

// OrderBookExchangeOutItemParams ...
type OrderBookExchangeOutItemParams struct {
	OrderDtlID      int `json:"order_dtl_id"`
	WarehouseItemID int `json:"warehouse_item_id"`
	OrderQty        int `json:"order_qty"`
}

// OrderBookExchangeInItemParams ...
type OrderBookExchangeInItemParams struct {
	ItemID    int     `json:"item_id"`
	OrderQty  int     `json:"order_qty"`  // Захиалгын тоо
	UnitPrice float64 `json:"unit_price"` // Нэгжийн үнэ
}

// OrderBookExchangeFilterCols ..
type OrderBookExchangeFilterCols struct {
	CustomerID  int       `json:"customer_id"`
	OrderNumber string    `json:"order_number"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// OrderBookExchangeFilter sort hiigdej boloh zuils
type OrderBookExchangeFilter struct {
	Page   int                         `json:"page"`
	Size   int                         `json:"size"`
	Sort   SortColumn                  `json:"sort"`
	Filter OrderBookExchangeFilterCols `json:"filter"`
}
