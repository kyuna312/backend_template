package form

import "time"

// WareHouseParams update body params
type WareHouseParams struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	IsActive           bool   `json:"is_active"`
	AddressDescription string `json:"address_description"`
	CityID             uint   `json:"city_id"`
}

// ChangePriceWareHouseItemParams update body params
type ChangePriceWareHouseItemParams struct {
	WarehouseItemID int     `json:"warehouse_item_id"`
	NewPrice        float64 `json:"new_price"`
	NewPriceType    int     `json:"new_price_type"`
	NewBarterPrice  float64 `json:"barter_price"`
}

// WareHouseFilterCols sort hiih bolomjtoi column
type WareHouseFilterCols struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	IsActive           string `json:"is_active"`
	AddressDescription string `json:"address_description"`
	CityID             uint   `json:"city_id"`
}

// WareHouseFilter sort hiigdej boloh zuils
type WareHouseFilter struct {
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Sort   SortColumn          `json:"sort"`
	Filter WareHouseFilterCols `json:"filter"`
}

// WarehouseItemFilter sort hiigdej boloh zuils
type WarehouseItemFilter struct {
	Page   int                    `json:"page"`
	Size   int                    `json:"size"`
	Sort   SortColumn             `json:"sort"`
	Filter WarehoseItemFilterCols `json:"filter"`
}

// WarehousePriceHistory ...
type WarehousePriceHistory struct {
	WarehouseItemID int `json:"warehouse_item_id"`
	PriceTypeID     int `json:"price_type_id"`
}

// WarehoseItemFilterCols ...
type WarehoseItemFilterCols struct {
	OnlyChild       bool      `json:"only_child"`
	Name            string    `json:"name"` //
	WarehouseID     int       `json:"warehouse_id"`
	CountryID       int       `json:"country_id"`
	MedDose         string    `json:"med_dose"`          // Тун хэмжээ
	SerialNumber    string    `json:"serial_number"`     //
	Packaging       string    `json:"packaging"`         // Савлагаа - No
	InterName       string    `json:"inter_name"`        // Олон улсын нэршил
	ItemCode        string    `json:"item_code"`         // Барааны код
	ManuName        string    `json:"manu_name"`         // Үйлдвэрлэгчийн нэр
	StartExpireDate time.Time `json:"start_expire_date"` // Дуусах огноо
	EndExpireDate   time.Time `json:"end_expire_date"`   // Дуусах огноо
	MeasureID       int       `json:"measure_id"`        //
	CustomerID      int       `json:"customer_id"`       //
	PriceTypeID     int       `json:"price_type_id"`
}
