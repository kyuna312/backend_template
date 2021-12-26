package form

import (
	"time"

	"gitlab.com/fibocloud/medtech/gin/databases"
)

// CustomerCreateParams create body params
type CustomerCreateParams struct {
	Name                 string                           `form:"name" binding:"required"`
	Description          string                           `form:"description"`
	ParentID             int                              `form:"parent_id"`
	CompanyRD            string                           `form:"company_rd"`
	PaymentTypeID        int                              `form:"payment_type_id"`
	MaximumReceivables   float64                          `form:"maximum_receivables"`
	OneTimePurchaseLimit float64                          `form:"one_time_purchase_limit"`
	Types                []*databases.MedCustomerType     `form:"customer_types" binding:"required"`
	ClassificationID     int                              `form:"classification_id" binding:"required"`
	CityID               int                              `form:"city_id" binding:"required"`
	DistrictID           int                              `form:"district_id" binding:"required"`
	Addresses            []*databases.MedCustomerAddress  `form:"addresses"`
	Contacts             []*databases.MedCustomerContacts `form:"contacts"`
}

// UpdateCustomerStatus ...
type UpdateCustomerStatus struct {
	CustomerID  int    `json:"customer_id" binding:"required"` //
	StatusID    int    `json:"status_id" binding:"required"`   //
	Description string `json:"description" binding:"required"` //
}

// GetPriceDetailParam ...
type GetPriceDetailParam struct {
	CustomerID int `json:"customer_id" binding:"required"` //
	ItemID     int `json:"item_id" binding:"required"`     //
}

// CustomerLoginPermissionParam ...
type CustomerLoginPermissionParam struct {
	CustomerID int    `json:"customer_id" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// CustomerUpdatePriceParams create body params
type CustomerUpdatePriceParams struct {
	CustomerIDs []int           `json:"customer_ids" binding:"required"`
	StartDate   time.Time       `json:"start_date"  binding:"required"`
	EndDate     time.Time       `json:"end_date"  binding:"required"`
	IsPercent   bool            `json:"is_percent"`
	Items       []*CustomerItem `json:"items"`
}

// CustomerItem ...
type CustomerItem struct {
	ItemID          int     `json:"item_id"`
	DiscountPrice   float64 `json:"discount_price"`
	DiscountPercent float64 `json:"discount_percent"`
}

// CustomerUpdateParams update body params
type CustomerUpdateParams struct {
	Name             string                       `json:"name"`
	Description      string                       `json:"description"`
	IsActive         bool                         `json:"is_active"`
	Types            []*databases.MedCustomerType `json:"customer_types" binding:"required"`
	ClassificationID int                          `json:"classification_id" binding:"required"`
	StatusID         int                          `json:"status_id" binding:"required"`
	CityID           int                          `json:"city_id" binding:"required"`
	DistrictID       int                          `json:"district_id" binding:"required"`
}

// CustomerFilterCols sort hiih bolomjtoi column
type CustomerFilterCols struct {
	Name                      string `json:"name"`
	Description               string `json:"description"`
	IsActive                  string `json:"is_active"`
	ClassificationID          int    `json:"classification_id"`
	StatusID                  int    `json:"status_id"`
	PaymentTypeID             int    `json:"payment_type_id"`
	CountryID                 int    `json:"country_id"`
	CityID                    int    `json:"city_id"`
	DistrictID                int    `json:"district_id"`
	CreatedUserID             int    `json:"created_user_id"`
	ModifiedUserID            int    `json:"modified_user_id"`
	ExternalRegistryNumber    string `json:"external_registry_number"`
	ExternalCustomerTypeID    int    `json:"external_customer_type_id"`
	ExternalContactPositionID int    `json:"external_contact_position_id"`
	ExternalAddressTypeID     int    `json:"external_address_type_id"`
	ExternalIsPercent         bool   `json:"external_is_percent"`
	ExternalContactPhone      string `json:"external_contact_phone"`
}

// CustomerFilter sort hiigdej boloh zuils
type CustomerFilter struct {
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Sort   SortColumn         `json:"sort"`
	Filter CustomerFilterCols `json:"filter"`
}
