package form

import "time"

// MarketingCreateOneTimeParams create body params
type MarketingCreateOneTimeParams struct {
	Name                 string                    `json:"name"`
	Description          string                    `json:"description"`
	MarketingOrderTypeID int                       `json:"marketing_order_type" binding:"required"`
	PaymentTypeIDs       []int                     `json:"payment_type_ids" binding:"required"`
	StartDate            time.Time                 `json:"start_date" binding:"required"`
	EndDate              time.Time                 `json:"end_date" binding:"required"`
	IsAllCustomer        bool                      `json:"is_all_customer"`
	CustomerDistrictIDs  []int                     `json:"customer_district_ids"`
	CustomerIDs          []int                     `json:"customer_ids"`
	SalesAmount          float64                   `json:"sales_amount" binding:"required"` // Нэг удаагийн борлуулалтын дүн
	IsTutamd             bool                      `json:"is_tutamd"`                       // Бараа тутамд
	IsMarketingItem      bool                      `json:"is_marketing_item"`               // Борлуулалтын бараа юу эсвэл энгийн бараа юу
	IsPercent            bool                      `json:"is_percent"`                      // Хувиар хөнглөлт байхуу
	PercentValue         float64                   `json:"percent_value"`
	Items                []*MarketingItemDtlParams `json:"items"`
}

// MarketingCreateTimeIntervalParams create body params
type MarketingCreateTimeIntervalParams struct {
	Name                 string                   `json:"name"`
	Description          string                   `json:"description"`
	MarketingOrderTypeID int                      `json:"marketing_order_type" binding:"required"`
	PaymentTypeIDs       []int                    `json:"payment_type_ids" binding:"required"`
	StartDate            time.Time                `json:"start_date" binding:"required"`
	EndDate              time.Time                `json:"end_date" binding:"required"`
	IsAllCustomer        bool                     `json:"is_all_customer"`
	CustomerIDs          []int                    `json:"customer_ids"`
	CustomerDistrictIDs  []int                    `json:"customer_district_ids"`
	SalesAmount          float64                  `json:"sales_amount" binding:"required"`
	IsTutamd             bool                     `json:"is_tutamd"`         // Бараа тутамд
	IsMarketingItem      bool                     `json:"is_marketing_item"` // Борлуулалтын бараа юу эсвэл энгийн бараа юу
	Items                []MarketingItemDtlParams `json:"items"`
}

// MarketingCreateItemDiscountParams create body params
type MarketingCreateItemDiscountParams struct {
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	MarketingOrderTypeID int       `json:"marketing_order_type" binding:"required"`
	PaymentTypeIDs       []int     `json:"payment_type_ids" binding:"required"`
	StartDate            time.Time `json:"start_date" binding:"required"`
	EndDate              time.Time `json:"end_date" binding:"required"`
	IsAllCustomer        bool      `json:"is_all_customer"`
	CustomerIDs          []int     `json:"customer_ids"`
	CustomerDistrictIDs  []int     `json:"customer_district_ids"`
	IsAllItems           bool      `json:"is_all_items"`                          // Бүх item дээр зоохуу
	WarehouseItemIDs     []int     `json:"warehouse_item_ids" binding:"required"` // Агуулахаас хайсан бараа
	DiscountPercent      float64   `json:"discount_percent" binding:"required"`   // Хөнглөлт хувь
}

// MarketingCreateItemRewardParams create body params
type MarketingCreateItemRewardParams struct {
	Name                 string                   `json:"name"`
	Description          string                   `json:"description"`
	MarketingOrderTypeID int                      `json:"marketing_order_type" binding:"required"`
	PaymentTypeIDs       []int                    `json:"payment_type_ids" binding:"required"`
	StartDate            time.Time                `json:"start_date" binding:"required"`
	EndDate              time.Time                `json:"end_date" binding:"required"`
	IsAllCustomer        bool                     `json:"is_all_customer"`
	CustomerIDs          []int                    `json:"customer_ids"`
	CustomerDistrictIDs  []int                    `json:"customer_district_ids"`
	IsAllItems           bool                     `json:"is_all_items"` // Бүх item дээр зоохуу
	IsTutamd             bool                     `json:"is_tutamd"`
	PurchasedItems       []*ItemRewadDtlParams    `json:"purchased_items" binding:"required"` // Агуулахаас хайсан бараа
	IsMarketingItem      bool                     `json:"is_marketing_item"`                  // Борлуулалтын бараа юу эсвэл энгийн бараа юу
	Items                []MarketingItemDtlParams `json:"items"`
}

// ItemRewadDtlParams ...
type ItemRewadDtlParams struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

// MarketingItemDtlParams ...
type MarketingItemDtlParams struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

// MarketingCreateItemCountDiscountParams create body params
type MarketingCreateItemCountDiscountParams struct {
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	MarketingOrderTypeID int       `json:"marketing_order_type" binding:"required"`
	PaymentTypeIDs       []int     `json:"payment_type_ids" binding:"required"`
	StartDate            time.Time `json:"start_date" binding:"required"`
	EndDate              time.Time `json:"end_date" binding:"required"`
	IsAllCustomer        bool      `json:"is_all_customer"`
	CustomerIDs          []int     `json:"customer_ids"`
	CustomerDistrictIDs  []int     `json:"customer_district_ids"`
	MinCount             int       `json:"min_count" binding:"required"`
	IsAllItems           bool      `json:"is_all_items"`                          // Бүх item дээр зоохуу
	WarehouseItemIDs     []int     `json:"warehouse_item_ids" binding:"required"` // Агуулахаас хайсан бараа
	DiscountPercent      float64   `json:"discount_percent" binding:"required"`   // Хөнглөлт хувь
}

// MarketingFilterCols sort hiih bolomjtoi column
type MarketingFilterCols struct {
	Name                   string    `json:"name"`
	IsActive               bool      `json:"is_active"`
	MarketingOrderTypeID   int       `json:"marketing_order_type_id"`
	MarketingTypeID        int       `json:"marketing_type_id"`
	ExternalStartDate      time.Time `json:"external_start_date"`
	ExternalEndDate        time.Time `json:"external_end_date"`
	ExternalCustomerIDs    []int     `json:"external_customer_ids"`
	ExtermalPaymentTypeIDs []int     `json:"external_payment_type_ids"`
}

// MarketingFilter sort hiigdej boloh zuils
type MarketingFilter struct {
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Sort   SortColumn          `json:"sort"`
	Filter MarketingFilterCols `json:"filter"`
}
