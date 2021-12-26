package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gitlab.com/fibocloud/medtech/gin/constracts"
	"gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
	utils "gitlab.com/fibocloud/medtech/gin/utils"
)

// AllRefs ...
type AllRefs struct {
	Classification []databases.MedCustomerClassification `json:"customer_classifications"`
	Types          []databases.MedCustomerType           `json:"types"`
	Parents        []databases.MedCustomer               `json:"parents"`
	Cities         []databases.RefCity                   `json:"cities"`
	PaymentTypes   []databases.MedPaymentType            `json:"payment_types"`
}

// CustomerController struct
type CustomerController struct {
	shared.BaseController
}

// ListCustomers ...
type ListCustomers struct {
	Total int64                   `json:"total"`
	List  []databases.MedCustomer `json:"list"`
}

// Init Controller
func (co CustomerController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)                                  // List
	router.GET("get/:id", co.Get)                                  // Show
	router.GET("/find/company/:rdCode", co.FindComapnyRD)          // Show
	router.POST("", co.Create)                                     // Create
	router.POST("/permission", co.CreatePermission)                // Erh uusgeh
	router.GET("/status/list", co.ListStatuses)                    // Status list
	router.POST("/status/change", co.ChangeStatus)                 // Status Change
	router.POST("/item/price", co.UpdateItemPrice)                 // Erh uusgeh
	router.POST("/list/price", co.ListPrice)                       // ListPrice
	router.POST("/get/price", co.GetPriceDetail)                   // ListPrice
	router.GET("/refs", co.Refs)                                   // Refs
	router.PUT("/:id", co.Update)                                  // Update
	router.GET("/country/:id", co.CustomerCountry)                 // CustomerCountry
	router.GET("/status/history/:id", co.HistoryStatuses)          // Status history
	router.DELETE("", co.Delete)                                   // Delete
	router.GET("/list/active", co.ListActive)                      // ListActive
	router.GET("/change/active/:id", co.ChangeActive)              // ChangeActive
	router.POST("/warehouse/item/list", co.WarehouseItemPriceList) // WarehouseItemPriceList
}

// HistoryStatuses customer
// @Summary HistoryStatuses customer
// @Description HistoryStatuses customer
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedStatusLog}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/status/list/{id} [get]
func (co CustomerController) HistoryStatuses(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var logs []databases.MedStatusLog
	co.DB.Where("record_id = ?", c.Param("id")).Where("hdr_table_name = ?", "med_customers").Preload("Status").Find(&logs)
	co.SetBody(logs)

	return
}

// FindComapnyRD customer
// @Summary FindComapnyRD customer
// @Description FindComapnyRD customer
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/find/company/{rdCode} [get]
func (co CustomerController) FindComapnyRD(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	type EbarintCompany struct {
		VatpayerRegisteredDate string `json:"vatpayer_registered_date"`
		LastReceiptDate        string `json:"last_receipt_date"`
		ReceiptFound           string `json:"receipt_found"`
		Name                   string `json:"name"`
		Citypayer              string `json:"citypayer"`
		Vatpayer               string `json:"vatpayer"`
		Found                  string `json:"found"`
	}

	type Response struct {
		IsRegistered   bool                    `json:"is_registered"`
		ChildCustomers []databases.MedCustomer `json:"child_customers"`
		Ebarint        interface{}             `json:"company"`
	}

	var response Response
	var ebarimtCompany EbarintCompany

	url := "http://info.ebarimt.mn/rest/merchant/info?regno=" + c.Param("rdCode")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	json.Unmarshal([]byte(string(body)), &ebarimtCompany)

	var customer databases.MedCustomer
	result := co.DB.Where("company_registry_number = ?", c.Param("rdCode")).First(&customer)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			response.IsRegistered = false
		}
	} else {
		var childCustomers []databases.MedCustomer
		co.DB.Where("parent_id = ?", customer.Base.ID).Preload("Classification").Find(&childCustomers)
		response.IsRegistered = true
		response.ChildCustomers = childCustomers
	}

	if ebarimtCompany.Name != "" {
		response.Ebarint = ebarimtCompany
	}
	co.SetBody(response)
	return
}

// ListPrice customer
// @Summary ListPrice customer
// @Description ListPrice customer
// @Tags City
// @Accept json
// @Produce json
// @Param filter body form.CustomerFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/list/price/ [post]
func (co CustomerController) ListPrice(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CustomerFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	type Child struct {
		ID           int                      `json:"id"`
		Item         *databases.MedItem       `json:"item"`             //
		Customer     *databases.MedCustomer   `json:"customer"`         // Үүсгэсэн хэрэглэгч
		PriceType    *databases.MedPriceType  `json:"price_type"`       // Зарах үнэ, солилцооны үнэ
		SalesPrice   float64                  `json:"sales_price"`      //
		IsPercent    bool                     `json:"is_percent"`       //
		Percent      float64                  `json:"customer_percent"` //
		StartDate    time.Time                `json:"start_date"`       //
		EndDate      time.Time                `json:"end_date"`         //
		IsActive     bool                     `json:"is_active"`        // Идэвхитэй эсэх
		CreatedUser  *databases.MedSystemUser `json:"created_user"`     // Үүсгэсэн хэрэглэгч
		ModifiedUser *databases.MedSystemUser `json:"modified_user"`    // Өөрчилсөн хэрэглэгч
	}

	type Parent struct {
		ID           int                      `json:"id"`
		Item         *databases.MedItem       `json:"item"`             //
		Customer     *databases.MedCustomer   `json:"customer"`         // Үүсгэсэн хэрэглэгч
		PriceType    *databases.MedPriceType  `json:"price_type"`       // Зарах үнэ, солилцооны үнэ
		SalesPrice   float64                  `json:"sales_price"`      //
		IsPercent    bool                     `json:"is_percent"`       //
		Percent      float64                  `json:"customer_percent"` //
		StartDate    time.Time                `json:"start_date"`       //
		EndDate      time.Time                `json:"end_date"`         //
		IsActive     bool                     `json:"is_active"`        // Идэвхитэй эсэх
		CreatedUser  *databases.MedSystemUser `json:"created_user"`     // Үүсгэсэн хэрэглэгч
		ModifiedUser *databases.MedSystemUser `json:"modified_user"`    // Өөрчилсөн хэрэглэгч
		Childrens    []Child                  `json:"children"`
	}

	// ListCustomersPrice ...
	type ListCustomersPrice struct {
		Total int      `json:"total"`
		List  []Parent `json:"list"`
	}

	var listRepsonse ListCustomersPrice

	var customerPrices []databases.MedPriceCustomer
	db.Preload("CreatedUser.Person").
		Preload("ModifiedUser.Person").
		Preload("Item").
		Preload("Customer").
		Where("is_active = ?", true).
		Where("is_percent = ?", params.Filter.ExternalIsPercent).
		Find(&customerPrices)

	var parents []Parent
	temp := make(map[uint]int)
	indexCounter := 1

	for _, i := range customerPrices {
		if temp[i.CustomerID] == 0 {

			temp[i.CustomerID] = indexCounter
			indexCounter = indexCounter + 1

			var eachParent Parent
			eachParent.ID = int(i.Base.ID)
			eachParent.Customer = i.Customer

			var eachChild Child
			eachChild.ID = int(i.Item.Base.ID)
			eachChild.Item = i.Item
			eachChild.SalesPrice = i.SalesPrice
			eachChild.IsPercent = i.IsPercent
			eachChild.Percent = i.Percent
			eachChild.StartDate = i.StartDate
			eachChild.EndDate = i.EndDate
			eachChild.IsActive = i.IsActive
			eachChild.CreatedUser = i.CreatedUser
			eachChild.ModifiedUser = i.ModifiedUser

			eachParent.Childrens = append(eachParent.Childrens, eachChild)
			parents = append(parents, eachParent)
		} else {
			var eachChild Child
			eachChild.ID = int(i.Item.Base.ID)
			eachChild.Item = i.Item
			eachChild.SalesPrice = i.SalesPrice
			eachChild.IsPercent = i.IsPercent
			eachChild.Percent = i.Percent
			eachChild.StartDate = i.StartDate
			eachChild.EndDate = i.EndDate
			eachChild.IsActive = i.IsActive
			eachChild.CreatedUser = i.CreatedUser
			eachChild.ModifiedUser = i.ModifiedUser

			parents[temp[i.CustomerID]-1].Childrens = append(parents[temp[i.CustomerID]-1].Childrens, eachChild)
		}
	}

	listRepsonse.List = parents
	listRepsonse.Total = len(parents)

	co.SetBody(listRepsonse)
	return
}

// ListActive customer
// @Summary ListActive customer
// @Description List customer
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/list/active [get]
func (co CustomerController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customers []databases.MedCustomer
	co.DB.Where("company_registry_number = ?", "").Find(&customers)
	co.SetBody(customers)
	return
}

// ChangeActive customer
// @Summary ChangeActive customer
// @Description ChangeActive customer
// @Tags City
// @Accept json
// @Produce json
// @Param id path uint true "customer ID"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/list/active/{id} [get]
func (co CustomerController) ChangeActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customer databases.MedCustomer
	result := co.DB.First(&customer, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	user := co.GetAuth(c)

	customer.IsActive = !customer.IsActive
	customer.ModifiedDate = time.Now()
	customer.ModifiedUser = &user
	result = co.DB.Save(&customer)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(customer)
	return
}

// ListStatuses customer
// @Summary ListStatuses customer
// @Description ListStatuses customer
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/status/list [get]
func (co CustomerController) ListStatuses(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()
	var statuses []databases.MedStatus
	co.DB.Where("status_type_id = ?", 3).Find(&statuses)
	co.SetBody(statuses)
	return
}

// CustomerCountry customer
// @Summary CustomerCountry customer
// @Description CustomerCountry customer
// @Tags Customer
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/country/{id} [get]
func (co CustomerController) CustomerCountry(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()
	var customers []databases.MedCustomer
	co.DB.Where("is_active = ?", true).Where("company_registry_number = ?", "").Where("country_id = ?", c.Param("id")).Find(&customers)
	co.SetBody(customers)
	return
}

// List customer
// @Summary List customer
// @Description Get customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param filter body form.CustomerFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/list [post]
func (co CustomerController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CustomerFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListCustomers
	var customers []databases.MedCustomer

	db.Preload("CustomerAddress.Country").
		Preload("CustomerAddress.City").
		Preload("CustomerAddress.District").
		Preload("CustomerAddress.Street").
		Preload("CustomerAddress.AddressType").
		Preload("Country").
		Preload("City").
		Preload("Status").
		Preload("District").
		Preload("PaymentType").
		Preload("Classification").
		Preload("Parent").
		Preload("CreatedUser.Person").
		Preload("ModifiedUser.Person").
		Where("company_registry_number = ?", "")

	if params.Filter.ExternalRegistryNumber != "" {
		var customerIDs []uint
		co.DB.Table("med_customers").Where("company_registry_number like '%"+params.Filter.ExternalRegistryNumber+"%'").Pluck("id", &customerIDs)

		if len(customerIDs) != 0 {
			db.Where("parent_id IN (?)", customerIDs)
		} else {
			co.SetBody(ListCustomers{})
			return
		}
	}

	if params.Filter.ExternalCustomerTypeID != 0 {
		var typeCustomerIDs []uint

		co.DB.Table("med_customer_type_dtl").Where("med_customer_type_id = ?", params.Filter.ExternalCustomerTypeID).Order("med_customer_id asc").Pluck("med_customer_id", &typeCustomerIDs)
		if len(typeCustomerIDs) != 0 {
			db.Preload("Types", "id IN (?)", params.Filter.ExternalCustomerTypeID).Where("id IN (?)", typeCustomerIDs)
		} else {
			co.SetBody(ListCustomers{})
			return
		}
	} else {
		db.Preload("Types")
	}

	if params.Filter.ExternalContactPhone != "" {
		var phoneCustomerIDs []uint
		co.DB.Table("med_customer_contacts").Where("phone_number1 like '%"+params.Filter.ExternalContactPhone+"%'").Or("phone_number2 like '%"+params.Filter.ExternalContactPhone+"%'").Pluck("customer_id", &phoneCustomerIDs)

		if len(phoneCustomerIDs) != 0 {
			db.Where("id IN (?)", phoneCustomerIDs)
		} else {
			co.SetBody(ListCustomers{})
			return
		}
	}

	if params.Filter.ExternalContactPositionID != 0 {
		var customerIDs []uint
		co.DB.Table("med_customer_contacts").Where("position_id = ?", params.Filter.ExternalContactPositionID).Pluck("customer_id", &customerIDs)

		if len(customerIDs) != 0 {
			db.Preload("Contacts.Position", "id IN (?)", params.Filter.ExternalContactPositionID).Where("id IN (?)", customerIDs)
		} else {
			co.SetBody(ListCustomers{})
			return
		}
	} else {
		db.Preload("Contacts.Position")
	}

	db.Find(&customers)

	listRepsonse.List = customers
	listRepsonse.Total = int64(len(customers))

	co.SetBody(listRepsonse)
	return
}

// Get customer
// @Summary Get customer
// @Description Show customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path uint true "customer ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/{id} [get]
func (co CustomerController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customer databases.MedCustomer
	result := co.DB.
		Preload("CustomerAddress.Country").
		Preload("CustomerAddress.City").
		Preload("CustomerAddress.District").
		Preload("CustomerAddress.Street").
		Preload("CustomerAddress.AddressType").
		Preload("Types").
		Preload("Country").
		Preload("City").
		Preload("Contacts.Position").
		Preload("Status").
		Preload("District").
		Preload("PaymentType").
		Preload("Classification").
		Preload("Parent").
		Preload("Files.ContentType").
		Preload("CreatedUser.Person").
		Preload("ModifiedUser.Person").
		First(&customer, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(customer)
	return
}

// GetPriceDetail customer
// @Summary GetPriceDetail customer
// @Description GetPriceDetail customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param Customer body form.GetPriceDetailParam true "Customer"
// @Success 200 {object} structs.ResponseBody{body=databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/get/price [post]
func (co CustomerController) GetPriceDetail(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.GetPriceDetailParam
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var customerPrice databases.MedPriceCustomer
	result := co.DB.
		Where("customer_id = ?", params.CustomerID).
		Where("item_id = ?", params.ItemID).
		Where("is_active = ?", true).
		Preload("CreatedUser.Person").
		Preload("ModifiedUser.Person").
		Last(&customerPrice)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	type Reponse struct {
		WareHouseItem *databases.MedWarehouseItem `json:"warehouse_item"` // Агуулах ID
		Price         float64                     `json:"price"`
	}

	var responses []Reponse

	var warehousePrices []databases.MedSalesPriceDtl

	var ids []int
	result = co.DB.
		Where("warehouse_item_id IN (?)", co.DB.Table("med_warehouse_items").Where("item_id = ?", params.ItemID).Not("total_qty = ?", 0).Pluck("id", &ids)).
		Where("is_active = ?", true).
		Where("price_type_id = ?", 1).
		Preload("WarehouseItem.Item").
		Find(&warehousePrices)

	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	for _, i := range warehousePrices {
		var each Reponse
		each.WareHouseItem = i.WarehouseItem

		if customerPrice.IsPercent {
			each.Price = i.SalesPrice + i.SalesPrice*customerPrice.Percent/100
		} else {
			each.Price = i.SalesPrice
		}

		responses = append(responses, each)
	}

	co.SetBody(responses)
	return
}

// ChangeStatus Customer
// @Summary ChangeStatus Customer
// @Description ChangeStatus Customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param Customer body form.UpdateCustomerStatus true "Customer"
// @Success 200 {object} structs.ResponseBody{body=databases.MedCustomer}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/status/change [post]
func (co CustomerController) ChangeStatus(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.UpdateCustomerStatus
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	tx := co.DB.Begin()

	var customer databases.MedCustomer
	result := co.DB.First(&customer, params.CustomerID)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	customer.StatusID = uint(params.StatusID)
	resultOldSave := tx.Save(&customer)
	if resultOldSave.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, "Харилцагчийн төлөв өөрчлөхөд алдаа гарлаа "+resultOldSave.Error.Error())
		return
	}

	authUser := co.GetAuth(c)
	customerStatusLog := databases.MedStatusLog{
		CreatedUser:  &authUser,
		RecordID:     uint(params.CustomerID),
		StatusID:     uint(params.StatusID),
		Description:  params.Description,
		HdrTableName: "med_customers",
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result = tx.Create(&customerStatusLog)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	tx.Commit()
	co.SetBody(customer)
	return
}

// Create customer
// @Summary Create customer
// @Description Add customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body form.CustomerCreateParams true "customer"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer [post]
func (co CustomerController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	name, errName := c.GetPostForm("name")
	if !errName {
		co.SetError(http.StatusInternalServerError, "Харилцагч нэр оруулна уу!")
		return
	}

	countryID, errCountryID := c.GetPostForm("country_id")
	if !errCountryID {
		co.SetError(http.StatusInternalServerError, "Харилцагчийн улсыг оруулна уу")
		return
	}

	description, _ := c.GetPostForm("description")
	companyRd, errCompanyRd := c.GetPostForm("company_rd")
	if countryID == strconv.Itoa(constracts.MongoliaCountryCode) {
		if !errCompanyRd {
			co.SetError(http.StatusInternalServerError, "Байгууллагын РД оруулна уу!")
			return
		}
	}

	companyName, errCompanyName := c.GetPostForm("company_name")
	if countryID == strconv.Itoa(constracts.MongoliaCountryCode) {
		if !errCompanyName {
			co.SetError(http.StatusInternalServerError, "Компани нэр оруулна уу!")
			return
		}
	}

	cityID := ""
	districtID := ""
	errCityID := false
	addressDescription := ""
	intCountryID, _ := strconv.Atoi(countryID)

	if intCountryID == constracts.MongoliaCountryCode {
		cityID, errCityID = c.GetPostForm("city_id")
		if !errCityID {
			co.SetError(http.StatusInternalServerError, "Харилцагч аймгыг оруулна уу")
			return
		}
		districtID, errCityID = c.GetPostForm("district_id")
		if !errCityID {
			co.SetError(http.StatusInternalServerError, "Харилцагч сум/дүүрэг оруулна уу")
			return
		}
	} else {
		addressDescription, _ = c.GetPostForm("address_description")
	}

	paymentTypeID, errPaymentTypeID := c.GetPostForm("payment_type_id")
	if !errPaymentTypeID {
		co.SetError(http.StatusInternalServerError, "Төлбөрийн нөхцөл")
		return
	}

	maximumPurchase, _ := c.GetPostForm("maximum_purchase")
	maximumReceivables, _ := c.GetPostForm("maximum_receivables")
	oneTimePurchaseLimit, _ := c.GetPostForm("one_time_purchase_limit")
	classificationID, errClassificationID := c.GetPostForm("classification_id")
	if !errClassificationID {
		co.SetError(http.StatusInternalServerError, "Ангилал сонгоно уу")
		return
	}

	customerTypesString, errCustomerTypes := c.GetPostForm("customer_types")
	if !errCustomerTypes {
		co.SetError(http.StatusInternalServerError, "Төрөл сонгоно уу")
		return
	}
	addressesString, _ := c.GetPostForm("addresses")
	contactsString, _ := c.GetPostForm("contacts")

	var customerTypes []*databases.MedCustomerType
	json.Unmarshal([]byte(customerTypesString), &customerTypes)

	intDistrictID, _ := strconv.Atoi(districtID)
	intClassificationID, _ := strconv.Atoi(classificationID)
	intOneTimePurchaseLimit, _ := strconv.Atoi(oneTimePurchaseLimit)
	intMaximumReceivables, _ := strconv.Atoi(maximumReceivables)
	intMaximumPurchase, _ := strconv.Atoi(maximumPurchase)
	intCityID, _ := strconv.Atoi(cityID)

	intPaymentTypeID, _ := strconv.Atoi(paymentTypeID)

	// types, _ := utils.JSONUnmarshal([]byte(customerTypes))
	tx := co.DB.Begin()
	authUser := co.GetAuth(c)
	formdata, _ := c.MultipartForm()
	minioClient := co.MinioClinet()

	var preCustomer databases.MedCustomer
	todayTimeStump := time.Now().Format("200601")
	resultPre := co.DB.Where("code LIKE ?", "%"+todayTimeStump+"%").Last(&preCustomer)

	var parentCustomer databases.MedCustomer

	if intCountryID == constracts.MongoliaCountryCode {
		result := co.DB.Where("company_registry_number = ?", companyRd).First(&parentCustomer)
		// 202101020002
		if result.Error != nil {
			if result.Error.Error() == "record not found" {
				directorCards := formdata.File["director_cards"]
				certs := formdata.File["certifications"]
				licenses := formdata.File["licenses"]
				var allFiles []*databases.MedContent
				// Parent бүртгэх хэрэгтэй байна

				newCode := time.Now().Format("200601")

				if resultPre.Error != nil {
					newCode = newCode + "001"
				} else {
					// 202012290001
					tempstr := preCustomer.Code[6:9]
					intCode, intErr := strconv.Atoi(tempstr)
					if intErr != nil {
						co.SetError(http.StatusBadRequest, intErr.Error())
						return
					}
					s := fmt.Sprintf("%03d", intCode+1)
					newCode = newCode + s
				}

				if intCountryID == constracts.MongoliaCountryCode {
					companyName = name
					// if directorCards == nil {
					// 	tx.Rollback()
					// 	co.SetError(http.StatusInternalServerError, "Захиралын иргэнийн үнэмлэх оруулана уу ")
					// 	return
					// }

					// if certs == nil {
					// 	tx.Rollback()
					// 	co.SetError(http.StatusInternalServerError, "Улсын бүртгэлийн гэрчилгээ хуулбар оруулна уу ")
					// 	return
					// }

					// if licenses == nil {
					// 	tx.Rollback()
					// 	co.SetError(http.StatusInternalServerError, "Тусгай зөвшөөрөл оруулна уу ")
					// 	return
					// }
				}

				createdParentCustomer := databases.MedCustomer{
					Name:                  companyName,
					Code:                  newCode,
					Description:           description,
					IsActive:              true,
					AddressDescription:    addressDescription,
					CountryID:             uint(intCountryID),
					CityID:                uint(intCityID),
					DistrictID:            uint(intDistrictID),
					ClassificationID:      4,
					PaymentTypeID:         uint(intPaymentTypeID),
					MaximumReceivables:    float64(intMaximumReceivables),
					OneTimePurchaseLimit:  float64(intOneTimePurchaseLimit),
					MaximumPurchase:       float64(intMaximumPurchase),
					CompanyRegistryNumber: companyRd,
					StatusID:              uint(constracts.CustomerAccountConfirmed),
					CreatedUser:           &authUser,
					ModifiedUser:          &authUser,
					Base: databases.Base{
						CreatedDate: time.Now(),
					},
				}

				result = tx.Create(&createdParentCustomer)
				if result.Error != nil {
					tx.Rollback()
					co.SetError(http.StatusInternalServerError, "Харилцагч бүртгэж чадсангүй"+result.Error.Error())
					return
				}

				customerStatusLog := databases.MedStatusLog{
					CreatedUser:  &authUser,
					RecordID:     uint(createdParentCustomer.Base.ID),
					StatusID:     uint(constracts.CustomerAccountConfirmed),
					Description:  "",
					HdrTableName: "med_customers",
					Base: databases.Base{
						CreatedDate: time.Now(),
					},
				}

				result = tx.Create(&customerStatusLog)
				if result.Error != nil {
					tx.Rollback()
					co.SetError(http.StatusInternalServerError, result.Error.Error())
					return
				}

				bucketName := companyRd
				bucketError := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
				if bucketError != nil {
					fmt.Println("bucketError", bucketError.Error())
				}

				for _, cardFile := range directorCards {
					types := cardFile.Header["Content-Type"]
					openFile, errorfileOpen := cardFile.Open()
					if errorfileOpen != nil {
						co.SetError(http.StatusInternalServerError, "Файл нээх үед алдаа гарлаа "+errorfileOpen.Error())
					}

					fileName := time.Now().Format("20060102150405") + " - " + cardFile.Filename

					uploadInfo, errUploadInfo := minioClient.PutObject(context.Background(), bucketName, fileName, openFile, cardFile.Size, minio.PutObjectOptions{ContentType: types[0]})
					if errUploadInfo != nil {
						co.SetError(http.StatusInternalServerError, "Файлын санруу хуулах үед алдаа гарлаа "+errUploadInfo.Error())
						return
					}
					eachFile := databases.MedContent{
						FileName:      fileName,
						PhysicalPath:  uploadInfo.Bucket + "/" + uploadInfo.Key,
						ContentTypeID: uint(constracts.ContentDirectorCards),
						FileSize:      float64(cardFile.Size),
						Extention:     filepath.Ext(cardFile.Filename),
						CreatedUser:   &authUser,
						Base: databases.Base{
							CreatedDate: time.Now(),
						},
					}

					result := tx.Create(&eachFile)
					if result.Error != nil {
						tx.Rollback()
						co.SetError(http.StatusInternalServerError, "Файл датабайзруу хуулах үед алдаа гарлаа "+result.Error.Error())
						return
					}

					allFiles = append(allFiles, &eachFile)
				}

				for _, certFile := range certs {
					types := certFile.Header["Content-Type"]
					openFile, errorfileOpen := certFile.Open()
					if errorfileOpen != nil {
						co.SetError(http.StatusInternalServerError, "Файл нээх үед алдаа гарлаа "+errorfileOpen.Error())
					}
					fileName := time.Now().Format("20060102150405") + " - " + certFile.Filename

					uploadInfo, errUploadInfo := minioClient.PutObject(context.Background(), bucketName, fileName, openFile, certFile.Size, minio.PutObjectOptions{ContentType: types[0]})
					if errUploadInfo != nil {
						co.SetError(http.StatusInternalServerError, "Файлын санруу хуулах үед алдаа гарлаа "+errUploadInfo.Error())
						return
					}

					eachFile := databases.MedContent{
						FileName:      fileName,
						PhysicalPath:  uploadInfo.Bucket + "/" + uploadInfo.Key,
						ContentTypeID: uint(constracts.ContentCertification),
						FileSize:      float64(certFile.Size),
						Extention:     filepath.Ext(certFile.Filename),
						CreatedUser:   &authUser,
						Base: databases.Base{
							CreatedDate: time.Now(),
						},
					}

					result = tx.Create(&eachFile)
					if result.Error != nil {
						tx.Rollback()
						co.SetError(http.StatusInternalServerError, "Файл датабайзруу хуулах үед алдаа гарлаа "+result.Error.Error())
						return
					}

					allFiles = append(allFiles, &eachFile)

					createdParentCustomer.Files = allFiles
					createdParentCustomer.Types = customerTypes

					result = tx.Save(&createdParentCustomer)
					if result.Error != nil {
						tx.Rollback()
						co.SetError(http.StatusInternalServerError, result.Error.Error())
						return
					}
				}

				// 2021010222203
				tempstr := newCode
				intCode, intErr := strconv.Atoi(tempstr)
				if intErr != nil {
					co.SetError(http.StatusBadRequest, intErr.Error())
					return
				}
				s := fmt.Sprintf("%03d", intCode+1)

				var allLicenses []*databases.MedContent
				customer := databases.MedCustomer{
					Name:                 name,
					CompanyName:          createdParentCustomer.Name,
					Code:                 s,
					Description:          description,
					IsActive:             true,
					AddressDescription:   addressDescription,
					CountryID:            uint(intCountryID),
					CityID:               uint(intCityID),
					DistrictID:           uint(intDistrictID),
					ClassificationID:     uint(intClassificationID),
					PaymentTypeID:        uint(intPaymentTypeID),
					MaximumReceivables:   float64(intMaximumReceivables),
					OneTimePurchaseLimit: float64(intOneTimePurchaseLimit),
					MaximumPurchase:      float64(intMaximumPurchase),
					ParentID:             createdParentCustomer.Base.ID,
					StatusID:             uint(constracts.CustomerAccountConfirmed),
					CreatedUser:          &authUser,
					ModifiedUser:         &authUser,
					Base: databases.Base{
						CreatedDate: time.Now(),
					},
				}

				result = tx.Create(&customer)
				if result.Error != nil {
					tx.Rollback()
					co.SetError(http.StatusInternalServerError, "Харилцагч бүртгэж чадсангүй"+result.Error.Error())
					return
				}

				customersStatusLog := databases.MedStatusLog{
					CreatedUser:  &authUser,
					RecordID:     uint(customer.Base.ID),
					StatusID:     uint(constracts.CustomerAccountConfirmed),
					Description:  "",
					HdrTableName: "med_customers",
					Base: databases.Base{
						CreatedDate: time.Now(),
					},
				}

				result = tx.Create(&customersStatusLog)
				if result.Error != nil {
					tx.Rollback()
					co.SetError(http.StatusInternalServerError, result.Error.Error())
					return
				}

				for _, fileLicense := range licenses {
					types := fileLicense.Header["Content-Type"]
					openFile, errorfileOpen := fileLicense.Open()
					if errorfileOpen != nil {
						co.SetError(http.StatusInternalServerError, "Файл нээх үед алдаа гарлаа "+errorfileOpen.Error())
					}

					fileName := name + "/" + time.Now().Format("20060102150405") + " - " + fileLicense.Filename
					uploadInfo, errUploadInfo := minioClient.PutObject(context.Background(), bucketName, fileName, openFile, fileLicense.Size, minio.PutObjectOptions{ContentType: types[0]})
					if errUploadInfo != nil {
						co.SetError(http.StatusInternalServerError, "Файлын санруу хуулах үед алдаа гарлаа "+errUploadInfo.Error())
						return
					}
					eachFile := databases.MedContent{
						FileName:      fileName,
						PhysicalPath:  uploadInfo.Bucket + "/" + uploadInfo.Key,
						ContentTypeID: uint(constracts.ContentLicence),
						FileSize:      float64(fileLicense.Size),
						CreatedUser:   &authUser,
						Extention:     filepath.Ext(fileLicense.Filename),
						Base: databases.Base{
							CreatedDate: time.Now(),
						},
					}

					result = tx.Create(&eachFile)
					if result.Error != nil {
						tx.Rollback()
						co.SetError(http.StatusInternalServerError, "Файл датабайзруу хуулах үед алдаа гарлаа "+result.Error.Error())
						return
					}

					allLicenses = append(allLicenses, &eachFile)
				}

				if contactsString != "" {
					var customerContacts []*databases.MedCustomerContacts
					json.Unmarshal([]byte(contactsString), &customerContacts)

					for _, contanct := range customerContacts {
						eachContact := databases.MedCustomerContacts{
							CustomerID:     customer.Base.ID,
							LastName:       strings.Trim(contanct.LastName, " "),
							FirstName:      strings.Trim(contanct.FirstName, " "),
							RegisterNumber: strings.Trim(contanct.RegisterNumber, " "),
							PositionID:     contanct.PositionID,
							PhoneNumber1:   strings.Trim(contanct.PhoneNumber1, " "),
							PhoneNumber2:   strings.Trim(contanct.PhoneNumber2, " "),
							Email1:         strings.Trim(contanct.Email1, " "),
							Email2:         strings.Trim(contanct.Email2, " "),
							CreatedUser:    &authUser,
							ModifiedUser:   &authUser,
							Base: databases.Base{
								CreatedDate: time.Now(),
							},
						}

						result := tx.Create(&eachContact)
						if result.Error != nil {
							tx.Rollback()
							co.SetError(http.StatusInternalServerError, "Холбоо барих алдаа гарлаа "+result.Error.Error())
							return
						}
					}
				}

				if addressesString != "" {
					var customerAddresses []*databases.MedCustomerAddress
					json.Unmarshal([]byte(addressesString), &customerAddresses)

					for _, address := range customerAddresses {
						eachAddress := databases.MedCustomerAddress{
							CustomerID:    customer.Base.ID,
							CountryID:     uint(constracts.MongoliaCountryCode),
							CityID:        uint(intCityID),
							DistrictID:    uint(intDistrictID),
							StreetID:      address.StreetID,
							AddressTypeID: address.AddressTypeID,
							Description:   address.Description,
							CreatedUser:   &authUser,
							ModifiedUser:  &authUser,
							Base: databases.Base{
								CreatedDate: time.Now(),
							},
						}

						result := tx.Create(&eachAddress)
						if result.Error != nil {
							tx.Rollback()
							co.SetError(http.StatusInternalServerError, "Хаяг байршил нэмэх алдаа гарлаа "+result.Error.Error())
							return
						}
					}
				}

				if intCountryID != constracts.MongoliaCountryCode {
					eachAddress := databases.MedCustomerAddress{
						CustomerID:    customer.Base.ID,
						CountryID:     uint(intCountryID),
						AddressTypeID: 1,
						Description:   addressDescription,
						CreatedUser:   &authUser,
						ModifiedUser:  &authUser,
						Base: databases.Base{
							CreatedDate: time.Now(),
						},
					}

					result := tx.Create(&eachAddress)
					if result.Error != nil {
						tx.Rollback()
						co.SetError(http.StatusInternalServerError, "Хаяг байршил нэмэх алдаа гарлаа "+result.Error.Error())
						return
					}
				}

				customer.Files = allLicenses
				customer.Types = customerTypes

				result = tx.Save(&customer)
				if result.Error != nil {
					tx.Rollback()
					co.SetError(http.StatusInternalServerError, result.Error.Error())
					return
				}

				co.SetBody(structs.SuccessResponse{
					Success: true,
				})

				tx.Commit()
				return
			}
			tx.Rollback()
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}
	}

	newParentID := parentCustomer.Base.ID
	newParentName := parentCustomer.Name

	if intCountryID != constracts.MongoliaCountryCode {
		newParentID = 0
		newParentName = ""
	}

	newCode := time.Now().Format("200601")
	var allLicenses []*databases.MedContent

	if resultPre.Error != nil {
		newCode = newCode + "001"
	} else {
		// 202012290001
		fmt.Println(preCustomer.Code)
		tempstr := preCustomer.Code[6:9]
		fmt.Println(tempstr)
		intCode, intErr := strconv.Atoi(tempstr)
		if intErr != nil {
			co.SetError(http.StatusBadRequest, intErr.Error())
			return
		}
		s := fmt.Sprintf("%03d", intCode+1)
		newCode = newCode + s
	}

	fmt.Println(newCode)

	customer := databases.MedCustomer{
		Name:                 name,
		CompanyName:          newParentName,
		Code:                 newCode,
		Description:          description,
		IsActive:             true,
		AddressDescription:   addressDescription,
		CountryID:            uint(intCountryID),
		CityID:               uint(intCityID),
		DistrictID:           uint(intDistrictID),
		ClassificationID:     uint(intClassificationID),
		PaymentTypeID:        uint(intPaymentTypeID),
		MaximumReceivables:   float64(intMaximumReceivables),
		OneTimePurchaseLimit: float64(intOneTimePurchaseLimit),
		MaximumPurchase:      float64(intMaximumPurchase),
		ParentID:             newParentID,
		StatusID:             uint(constracts.CustomerAccountConfirmed),
		CreatedUser:          &authUser,
		ModifiedUser:         &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := tx.Create(&customer)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, "Харилцагч бүртгэж чадсангүй"+result.Error.Error())
		return
	}

	childCustomer := databases.MedStatusLog{
		CreatedUser:  &authUser,
		RecordID:     uint(customer.Base.ID),
		StatusID:     uint(constracts.CustomerAccountConfirmed),
		Description:  "",
		HdrTableName: "med_customers",
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result = tx.Create(&childCustomer)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	bucketName := companyRd
	licenses := formdata.File["licenses"]

	// if intCountryID == constracts.MongoliaCountryCode {
	// 	if len(licenses) == 0 {
	// 		tx.Rollback()
	// 		co.SetError(http.StatusInternalServerError, "Тусгай зөвшөөрөл оруулна уу")
	// 		return
	// 	}
	// }

	for _, fileLicense := range licenses {
		types := fileLicense.Header["Content-Type"]
		openFile, errorfileOpen := fileLicense.Open()
		if errorfileOpen != nil {
			co.SetError(http.StatusInternalServerError, "Файл нээх үед алдаа гарлаа "+errorfileOpen.Error())
		}

		fileName := name + "/" + time.Now().Format("20060102150405") + " - " + fileLicense.Filename
		uploadInfo, errUploadInfo := minioClient.PutObject(context.Background(), bucketName, fileName, openFile, fileLicense.Size, minio.PutObjectOptions{ContentType: types[0]})
		if errUploadInfo != nil {
			co.SetError(http.StatusInternalServerError, "Файлын санруу хуулах үед алдаа гарлаа "+errUploadInfo.Error())
			return
		}
		eachFile := databases.MedContent{
			FileName:      fileName,
			PhysicalPath:  uploadInfo.Bucket + "/" + uploadInfo.Key,
			ContentTypeID: uint(constracts.ContentLicence),
			FileSize:      float64(fileLicense.Size),
			CreatedUser:   &authUser,
			Extention:     filepath.Ext(fileLicense.Filename),
			Base: databases.Base{
				CreatedDate: time.Now(),
			},
		}

		result = tx.Create(&eachFile)
		if result.Error != nil {
			tx.Rollback()
			co.SetError(http.StatusInternalServerError, "Файл датабайзруу хуулах үед алдаа гарлаа "+result.Error.Error())
			return
		}

		allLicenses = append(allLicenses, &eachFile)
	}

	if contactsString != "" {
		var customerContacts []*databases.MedCustomerContacts
		json.Unmarshal([]byte(contactsString), &customerContacts)

		for _, contanct := range customerContacts {
			eachContact := databases.MedCustomerContacts{
				CustomerID:     customer.Base.ID,
				LastName:       strings.Trim(contanct.LastName, " "),
				FirstName:      strings.Trim(contanct.FirstName, " "),
				RegisterNumber: strings.Trim(contanct.RegisterNumber, " "),
				PositionID:     contanct.PositionID,
				PhoneNumber1:   strings.Trim(contanct.PhoneNumber1, " "),
				PhoneNumber2:   strings.Trim(contanct.PhoneNumber2, " "),
				Email1:         strings.Trim(contanct.Email1, " "),
				Email2:         strings.Trim(contanct.Email2, " "),
				CreatedUser:    &authUser,
				ModifiedUser:   &authUser,
				Base: databases.Base{
					CreatedDate: time.Now(),
				},
			}

			result := tx.Create(&eachContact)
			if result.Error != nil {
				tx.Rollback()
				co.SetError(http.StatusInternalServerError, "Холбоо барих алдаа гарлаа "+result.Error.Error())
				return
			}
		}
	}

	if addressesString != "" {
		var customerAddresses []*databases.MedCustomerAddress
		json.Unmarshal([]byte(addressesString), &customerAddresses)

		for _, address := range customerAddresses {
			eachAddress := databases.MedCustomerAddress{
				CustomerID:    customer.Base.ID,
				CountryID:     uint(constracts.MongoliaCountryCode),
				CityID:        uint(intCityID),
				DistrictID:    uint(intDistrictID),
				StreetID:      address.StreetID,
				AddressTypeID: address.AddressTypeID,
				Description:   address.Description,
				CreatedUser:   &authUser,
				ModifiedUser:  &authUser,
				Base: databases.Base{
					CreatedDate: time.Now(),
				},
			}

			result := tx.Create(&eachAddress)
			if result.Error != nil {
				tx.Rollback()
				co.SetError(http.StatusInternalServerError, "Хаяг байршил нэмэх алдаа гарлаа "+result.Error.Error())
				return
			}
		}
	}

	if intCountryID != constracts.MongoliaCountryCode {
		eachAddress := databases.MedCustomerAddress{
			CustomerID:    customer.Base.ID,
			CountryID:     uint(intCountryID),
			AddressTypeID: 1,
			Description:   addressDescription,
			CreatedUser:   &authUser,
			ModifiedUser:  &authUser,
			Base: databases.Base{
				CreatedDate: time.Now(),
			},
		}

		result := tx.Create(&eachAddress)
		if result.Error != nil {
			tx.Rollback()
			co.SetError(http.StatusInternalServerError, "Хаяг байршил нэмэх алдаа гарлаа "+result.Error.Error())
			return
		}
	}

	customer.Files = allLicenses
	customer.Types = customerTypes

	result = tx.Save(&customer)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	customerStatusLog := databases.MedStatusLog{
		CreatedUser:  &authUser,
		RecordID:     uint(customer.Base.ID),
		StatusID:     uint(constracts.CustomerAccountConfirmed),
		Description:  "",
		HdrTableName: "med_customers",
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result = tx.Create(&customerStatusLog)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})

	tx.Commit()
	return
}

// CreatePermission customer
// @Summary CreatePermission customer
// @Description CreatePermission customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body form.CustomerLoginPermissionParam true "customer"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/permission [post]
func (co CustomerController) CreatePermission(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CustomerLoginPermissionParam
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	var customer databases.MedCustomer
	result := co.DB.First(&customer, params.CustomerID)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	customer.StatusID = uint(constracts.CustomerPermissionCreated)
	customer.ModifiedDate = time.Now()
	customer.ModifiedUser = &authUser
	result = co.DB.Save(&customer)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	statusLog := databases.MedStatusLog{
		CreatedUser:  &authUser,
		Description:  "",
		RecordID:     customer.Base.ID,
		StatusID:     uint(constracts.CustomerPermissionCreated),
		HdrTableName: "med_customer",
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result = co.DB.Create(&statusLog)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	hashPwd, err := utils.GenerateHash(params.Password)
	if err != nil {
		co.SetError(http.StatusInternalServerError, err.Error())
		return
	}

	systemUser := databases.MedSystemUser{
		IsActive:     true,
		Username:     params.Username,
		StartDate:    time.Now(),
		PersonID:     uint(params.CustomerID),
		PersonType:   2,
		PasswordHash: hashPwd,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result = co.DB.Create(&systemUser)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, "Нэвтрэх эрх өгөхөд алдаа гарлаа "+result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})

	return
}

// Update customer
// @Summary Update customer
// @Description Edit customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path uint true "customer ID"
// @Param customer body form.CustomerCreateParams true "customer"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/{id} [put]
func (co CustomerController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	name, _ := c.GetPostForm("name")
	countryID, _ := c.GetPostForm("country_id")
	description, _ := c.GetPostForm("description")
	companyRd, _ := c.GetPostForm("company_rd")

	cityID, _ := c.GetPostForm("city_id")
	districtID, _ := c.GetPostForm("district_id")
	addressDescription, _ := c.GetPostForm("address_description")

	intCountryID, _ := strconv.Atoi(countryID)
	intCustomerID, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("intCustomerID", intCustomerID)

	paymentTypeID, _ := c.GetPostForm("payment_type_id")
	maximumPurchase, _ := c.GetPostForm("maximum_purchase")
	maximumReceivables, _ := c.GetPostForm("maximum_receivables")
	classificationID, _ := c.GetPostForm("classification_id")
	customerTypesString, _ := c.GetPostForm("customer_types")
	addressesString, _ := c.GetPostForm("addresses")
	contactsString, _ := c.GetPostForm("contacts")

	var customerTypes []*databases.MedCustomerType
	json.Unmarshal([]byte(customerTypesString), &customerTypes)

	intDistrictID, _ := strconv.Atoi(districtID)
	intClassificationID, _ := strconv.Atoi(classificationID)
	intMaximumReceivables, _ := strconv.Atoi(maximumReceivables)
	intMaximumPurchase, _ := strconv.Atoi(maximumPurchase)
	intCityID, _ := strconv.Atoi(cityID)
	intPaymentTypeID, _ := strconv.Atoi(paymentTypeID)

	// types, _ := utils.JSONUnmarshal([]byte(customerTypes))
	tx := co.DB.Begin()
	authUser := co.GetAuth(c)

	var customer databases.MedCustomer
	co.DB.First(&customer, intCustomerID)

	var parentCustomer databases.MedCustomer
	if companyRd != "" {
		co.DB.Where("company_registry_number = ?", companyRd).First(&parentCustomer)
	}

	customer.Name = name
	customer.Description = description
	customer.AddressDescription = addressDescription
	customer.CountryID = uint(intCountryID)
	customer.CityID = uint(intCityID)
	customer.DistrictID = uint(intDistrictID)
	customer.ClassificationID = uint(intClassificationID)
	customer.PaymentTypeID = uint(intPaymentTypeID)
	customer.MaximumReceivables = float64(intMaximumReceivables)
	customer.MaximumPurchase = float64(intMaximumPurchase)
	customer.ParentID = parentCustomer.Base.ID
	customer.ModifiedUser = &authUser
	customer.ModifiedDate = time.Now()

	result := tx.Save(&customer)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, "Харилцагч бүртгэж чадсангүй"+result.Error.Error())
		return
	}

	if contactsString != "" {
		co.DB.Where("customer_id = ?", intCustomerID).Delete(&databases.MedCustomerContacts{})

		var customerContacts []*databases.MedCustomerContacts
		json.Unmarshal([]byte(contactsString), &customerContacts)

		for _, contanct := range customerContacts {
			eachContact := databases.MedCustomerContacts{
				CustomerID:     customer.Base.ID,
				LastName:       strings.Trim(contanct.LastName, " "),
				FirstName:      strings.Trim(contanct.FirstName, " "),
				RegisterNumber: strings.Trim(contanct.RegisterNumber, " "),
				PositionID:     contanct.PositionID,
				PhoneNumber1:   strings.Trim(contanct.PhoneNumber1, " "),
				PhoneNumber2:   strings.Trim(contanct.PhoneNumber2, " "),
				Email1:         strings.Trim(contanct.Email1, " "),
				Email2:         strings.Trim(contanct.Email2, " "),
				IsActive:       true,
				CreatedUser:    &authUser,
				ModifiedUser:   &authUser,
				Base: databases.Base{
					CreatedDate: time.Now(),
				},
			}

			result := tx.Create(&eachContact)
			if result.Error != nil {
				tx.Rollback()
				co.SetError(http.StatusInternalServerError, "Холбоо барих алдаа гарлаа "+result.Error.Error())
				return
			}
		}
	}

	if addressesString != "" {
		co.DB.Where("customer_id = ?", intCustomerID).Delete(&databases.MedCustomerAddress{})

		var customerAddresses []*databases.MedCustomerAddress
		json.Unmarshal([]byte(addressesString), &customerAddresses)

		for _, address := range customerAddresses {
			eachAddress := databases.MedCustomerAddress{
				CustomerID:    customer.Base.ID,
				CountryID:     uint(constracts.MongoliaCountryCode),
				CityID:        address.CityID,
				DistrictID:    address.DistrictID,
				StreetID:      address.StreetID,
				IsActive:      true,
				AddressTypeID: address.AddressTypeID,
				Description:   address.Description,
				CreatedUser:   &authUser,
				ModifiedUser:  &authUser,
				Base: databases.Base{
					CreatedDate: time.Now(),
				},
			}

			result := tx.Create(&eachAddress)
			if result.Error != nil {
				tx.Rollback()
				co.SetError(http.StatusInternalServerError, "Хаяг байршил нэмэх алдаа гарлаа "+result.Error.Error())
				return
			}
		}
	}

	// formdata, _ := c.MultipartForm()
	// minioClient := co.MinioClinet()
	// var allFiles []*databases.MedContent
	// directorCards := formdata.File["director_cards"]
	// certs := formdata.File["certifications"]
	// bucketName := companyRd
	// // licenses := formdata.File["licenses"]

	// for _, cardFile := range directorCards {
	// 	types := cardFile.Header["Content-Type"]
	// 	openFile, errorfileOpen := cardFile.Open()
	// 	if errorfileOpen != nil {
	// 		co.SetError(http.StatusInternalServerError, "Файл нээх үед алдаа гарлаа "+errorfileOpen.Error())
	// 	}

	// 	fileName := time.Now().Format("20060102150405") + " - " + cardFile.Filename

	// 	uploadInfo, errUploadInfo := minioClient.PutObject(context.Background(), bucketName, fileName, openFile, cardFile.Size, minio.PutObjectOptions{ContentType: types[0]})
	// 	if errUploadInfo != nil {
	// 		co.SetError(http.StatusInternalServerError, "Файлын санруу хуулах үед алдаа гарлаа "+errUploadInfo.Error())
	// 		return
	// 	}
	// 	eachFile := databases.MedContent{
	// 		FileName:      fileName,
	// 		PhysicalPath:  uploadInfo.Bucket + "/" + uploadInfo.Key,
	// 		ContentTypeID: uint(constracts.ContentDirectorCards),
	// 		FileSize:      float64(cardFile.Size),
	// 		Extention:     filepath.Ext(cardFile.Filename),
	// 		CreatedUser:   &authUser,
	// 		Base: databases.Base{
	// 			CreatedDate: time.Now(),
	// 		},
	// 	}

	// 	result := tx.Create(&eachFile)
	// 	if result.Error != nil {
	// 		tx.Rollback()
	// 		co.SetError(http.StatusInternalServerError, "Файл датабайзруу хуулах үед алдаа гарлаа "+result.Error.Error())
	// 		return
	// 	}

	// 	allFiles = append(allFiles, &eachFile)
	// }

	// for _, certFile := range certs {
	// 	types := certFile.Header["Content-Type"]
	// 	openFile, errorfileOpen := certFile.Open()
	// 	if errorfileOpen != nil {
	// 		co.SetError(http.StatusInternalServerError, "Файл нээх үед алдаа гарлаа "+errorfileOpen.Error())
	// 	}
	// 	fileName := time.Now().Format("20060102150405") + " - " + certFile.Filename

	// 	uploadInfo, errUploadInfo := minioClient.PutObject(context.Background(), bucketName, fileName, openFile, certFile.Size, minio.PutObjectOptions{ContentType: types[0]})
	// 	if errUploadInfo != nil {
	// 		co.SetError(http.StatusInternalServerError, "Файлын санруу хуулах үед алдаа гарлаа "+errUploadInfo.Error())
	// 		return
	// 	}

	// 	eachFile := databases.MedContent{
	// 		FileName:      fileName,
	// 		PhysicalPath:  uploadInfo.Bucket + "/" + uploadInfo.Key,
	// 		ContentTypeID: uint(constracts.ContentCertification),
	// 		FileSize:      float64(certFile.Size),
	// 		Extention:     filepath.Ext(certFile.Filename),
	// 		CreatedUser:   &authUser,
	// 		Base: databases.Base{
	// 			CreatedDate: time.Now(),
	// 		},
	// 	}

	// 	result = tx.Create(&eachFile)
	// 	if result.Error != nil {
	// 		tx.Rollback()
	// 		co.SetError(http.StatusInternalServerError, "Файл датабайзруу хуулах үед алдаа гарлаа "+result.Error.Error())
	// 		return
	// 	}

	// 	allFiles = append(allFiles, &eachFile)
	// }

	co.DB.Model(&databases.MedCustomer{}).Association("Types").Clear()

	customer.Types = customerTypes
	result = tx.Save(&customer)
	if result.Error != nil {
		tx.Rollback()
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})

	tx.Commit()
	return
}

// UpdateItemPrice customer
// @Summary UpdateItemPrice customer
// @Description UpdateItemPrice customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body form.CustomerUpdatePriceParams true "customer"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/item/price [post]
func (co CustomerController) UpdateItemPrice(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CustomerUpdatePriceParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	tx := co.DB.Begin()

	for _, customerID := range params.CustomerIDs {
		var customer databases.MedCustomer
		result := co.DB.First(&customer, customerID)
		if result.Error != nil {
			tx.Rollback()
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}

		authUser := co.GetAuth(c)

		for _, item := range params.Items {
			discountPrice := item.DiscountPrice
			discountPercent := item.DiscountPercent

			var checkItem databases.MedPriceCustomer
			result = co.DB.Where("customer_id = ?", customerID).Where("item_id = ?", item.ItemID).Where("start_date <= ?", time.Now()).
				Where("end_date <= ?", time.Now()).Find(&checkItem)

			if result.Error == nil {
				tx.Rollback()
				co.SetError(http.StatusInternalServerError, "("+customer.Name+") - "+"Энэхүү бараан дээр харицлагчийн үнэ байна. ")
				return
			}

			if params.IsPercent {
				discountPrice = 0
			} else {
				discountPercent = 0
			}

			salesPriceDtl := databases.MedPriceCustomer{
				ItemID:      uint(item.ItemID),
				CustomerID:  uint(customerID),
				PriceTypeID: 3,
				SalesPrice:  discountPrice,
				IsPercent:   params.IsPercent,
				Percent:     discountPercent,
				StartDate:   params.StartDate,
				EndDate:     params.EndDate,
				IsActive:    true,
				CreatedUser: &authUser,
				Base: databases.Base{
					CreatedDate: time.Now(),
				},
			}

			result = tx.Create(&salesPriceDtl)
			if result.Error != nil {
				tx.Rollback()
				co.SetError(http.StatusInternalServerError, "Үнэ оруулах үед алдаа гарлаа "+result.Error.Error())
				return
			}
		}
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})

	tx.Commit()
	return
}

// Delete customer
// @Summary Delete customer
// @Description Remove customer
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body form.DeleteParams true "customer"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/{id} [delete]
func (co CustomerController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedCustomer{}, v)
		if result.Error != nil {
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Refs get customer reference
// @Summary Refs customer
// @Description Refs customer
// @Tags OrderBook
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=AllRefs}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/refs [get]
func (co CustomerController) Refs(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var response AllRefs

	var classification []databases.MedCustomerClassification
	co.DB.Find(&classification)

	var paymentTypes []databases.MedPaymentType
	co.DB.Where("is_active = ?", true).Find(&paymentTypes)

	var types []databases.MedCustomerType
	co.DB.Where("is_active = ?", true).Find(&types)

	var customers []databases.MedCustomer
	co.DB.Where("company_registry_number = ?", "").Where("is_active = ?", true).Find(&customers)

	var cities []databases.RefCity
	co.DB.Where("country_id = ?", 67).Find(&cities)

	response.Classification = classification
	response.Types = types
	response.Parents = customers
	response.Cities = cities
	response.PaymentTypes = paymentTypes

	co.SetBody(response)
	return
}

// WarehouseItemPriceList get customer reference
// @Summary WarehouseItemPriceList customer
// @Description WarehouseItemPriceList customer
// @Tags OrderBook
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=AllRefs}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customer/warehouse/item/list [post]
func (co CustomerController) WarehouseItemPriceList(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.SortItems
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	type ListResponse struct {
		List  []databases.MedSalesPriceDtl `json:"list"`
		Total int                          `json:"total"`
	}

	var listRepsonse ListResponse
	var itemIds []int
	var priceDtl []databases.MedSalesPriceDtl
	co.DB.
		Preload("WarehouseItem.Item.Measure").
		Preload("WarehouseItem.Item.Country").
		Preload("PriceType").
		Where("is_active = ?", true).
		Where("price_type_id = ?", 1).
		Not("sales_price = ?", 0).
		Where("item_id IN (?)", co.DB.Table("med_warehouse_items").Not("total_qty = ?", 0).Pluck("item_id", &itemIds)).
		Find(&priceDtl)

	listRepsonse.List = priceDtl
	listRepsonse.Total = len(priceDtl)

	co.SetBody(listRepsonse)
	return
}
