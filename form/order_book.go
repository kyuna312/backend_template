package form

import (
	"time"

	"gitlab.com/fibocloud/medtech/gin/databases"
)

// OrderBookInParams update body params
type OrderBookInParams struct {
	Name              string                `json:"name" binding:"required"`              //
	Description       string                `json:"description"`                          // Тайлбар
	CustomerID        int                   `json:"customer_id" binding:"required"`       // Харилцагчийн ID
	OrderSubTypeID    int                   `json:"order_sub_type_id" binding:"required"` // Захиалгын дэд төрөл / Дотоод, Импортын, Солилцооны /
	PaymentTypeID     int                   `json:"payment_type_id" binding:"required"`   // Төлбөрийн төрлийн ID
	DueDate           time.Time             `json:"due_date"`                             // Дуусах хугацаа
	ValuteValue       float64               `json:"valute_value"`                         // Валют утга
	StatusDescription string                `json:"status_description"`                   // Төлөв шилжих үеийн тайлбар
	ValuteID          int                   `json:"valute_id"`                            // Валют ID
	IsVat             bool                  `json:"is_vat"`                               // Татвартай эсэх
	IsDiscount        bool                  `json:"is_discount"`                          // Хөнгөлөлттэй эсэх
	Items             []DetailOrderDtlParam `json:"items"`
}

// OrderBookOutParams update body params
type OrderBookOutParams struct {
	CustomerID           int                                        `json:"customer_id" binding:"required"`       // Харилцагчийн ID
	OrderSubTypeID       int                                        `json:"order_sub_type_id" binding:"required"` // Захиалгын дэд төрөл / Борлуулалтын энгийн, Тендерийн, Байршуулалт /
	PaymentTypeID        int                                        `json:"payment_type_id" binding:"required"`   // Төлбөрийн төрлийн ID
	PaymentMethodID      int                                        `json:"payment_method_id"`                    // Төлбөрийн хэлбэр ID
	WarehouseID          int                                        `json:"warehouse_id" binding:"required"`
	DeliveryPersonID     int                                        `json:"delivery_person_id" binding:"required"` //
	Description          string                                     `json:"description"`                           // Тайлбар
	IsVat                bool                                       `json:"is_vat"`                                // Татвартай эсэх
	IsDiscount           bool                                       `json:"is_discount"`                           // Хямдаралтай эсэх
	PercentVat           float64                                    `json:"percent_vat"`                           //
	Items                []DetailOrderBookOutParams                 `json:"items"`
	RewardWarehouseItems []*databases.MedMarketingOutWarehouseItems `json:"reward_warehouse_items"`
	RewardMarketingItems []*databases.MedMarketingOutRewardItems    `json:"reward_marketing_items"`
}

// DetailOrderBookOutParams ...
type DetailOrderBookOutParams struct {
	WarehouseItemID uint    `json:"warehouse_item_id"`              // WareHouseItemID
	ItemID          uint    `json:"item_id"`                        // Барааны ID
	PercentDiscount float64 `json:"percent_discount"`               // Хөнгөлөлт/хувь/
	OutcomeQty      uint    `json:"outcome_qty" binding:"required"` // Зарлага тоо
	Price           float64 `json:"price"`
}

// RewardWarehouseItems ...
type RewardWarehouseItems struct {
	Item  *databases.MedWarehouseItem `json:"item"`
	Count int                         `json:"count"`
}

// RewardMarketingItems ...
type RewardMarketingItems struct {
	Item  *databases.MedMarketingItem
	Count int
}

// OrderBookFilterCols sort hiih bolomjtoi column
type OrderBookFilterCols struct {
	Name                string    `json:"name"`               //
	Code                string    `json:"code"`               // код
	Description         string    `json:"description"`        // Тайлбар
	OrderNumber         string    `json:"order_number"`       // Захиалгын дугаар
	OrderDate           time.Time `json:"order_date"`         // Захиалгын огноо
	CustomerID          int       `json:"customer_id" `       // Харилцагчийн ID
	OrderTypeID         int       `json:"order_type_id"`      // Захиалгын төрөл /орлого, зарлага/
	OrderSubTypeID      int       `json:"order_sub_type_id"`  // Захиалгын дэд төрөл /Дотоод, Импортын, Солилцооны/
	StatusID            int       `json:"status_id"`          // Захиалгын төлөв
	PaymentTypeID       int       `json:"payment_type_id"`    // Төлбөрийн төрлийн ID
	DeliveryPersonID    int       `json:"delivery_person_id"` // Хүргэлт хийсэн ажилтны ID
	ReturnTypeID        int       `json:"return_type_id"`     // Буцаалтын төрлийн ID
	Discount            int       `json:"discount"`           // Хөнгөлөлтийн дүн
	Vat                 int       `json:"vat"`                // Татаварын дүн
	Total               int       `json:"total"`              // Захиалгын нийт дүн
	IsRemoved           bool      `json:"is_removed"`         // Буцаагдсан эсэх
	IsVat               bool      `json:"is_vat"`             // Татвартай эсэх
	IsPaid              bool      `json:"is_paid"`            // Төлбөр хийгдсэн эсэх
	IsDiscount          bool      `json:"is_discount"`        // Хөнгөлөлттэй эсэх
	IsDelivery          bool      `json:"is_delivery"`        // Хүргэгдсэн эсэх
	DeliveryDate        time.Time `json:"delivery_date"`      // Хүргэлт хийгдсэн огноо
	DueDate             time.Time `json:"due_date"`           // Дуусах хугацаа
	ValuteID            int       `json:"valute_id"`          // Валют ID
	ValuteValue         int       `json:"valute_value"`       // Валют утга
	ExternalStartDate   time.Time `json:"external_start_date"`
	ExternalEndDate     time.Time `json:"external_end_date"`
	ExternalWarehouseID int       `json:"external_warehouse_id"`
}

// OrderBookFilter sort hiigdej boloh zuils
type OrderBookFilter struct {
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Sort   SortColumn          `json:"sort"`
	Filter OrderBookFilterCols `json:"filter"`
}

// ChangeStatus ...
type ChangeStatus struct {
	OrderBookID int    `json:"order_book_id"`
	StatusID    int    `json:"status_id"`
	Description string `json:"Description"`
}
