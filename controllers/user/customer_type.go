package user

import (
	"net/http"
	"reflect"
	"time"

	gin "github.com/gin-gonic/gin"
	"gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// CustomerTypeController struct
type CustomerTypeController struct {
	shared.BaseController
}

// ListCustomerType ...
type ListCustomerType struct {
	Total int64                       `json:"total"`
	List  []databases.MedCustomerType `json:"list"`
}

// Init Controller
func (co CustomerTypeController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)             // List
	router.GET("/list/active", co.ListActive) // ListActive
	router.GET("get/:id", co.Get)             // Show
	router.POST("", co.Create)                // Create
	router.PUT("/:id", co.Update)             // Update
	router.DELETE("/:id", co.Delete)          // Delete
}

// List customerType
// @Summary List customerType
// @Description Get customerType
// @Tags CustomerType
// @Accept json
// @Produce json
// @Param filter body form.CustomerTypeFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomerType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerType/list [post]
func (co CustomerTypeController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.CustomerTypeFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg

	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListCustomerType

	var customerTypes []databases.MedCustomerType
	db.Find(&customerTypes)
	db.Table("med_customer_types").Count(&count)

	listRepsonse.List = customerTypes
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return
}

// ListActive customerType
// @Summary ListActive customerType
// @Description Get customerType
// @Tags CustomerType
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomerType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerType/list/active [get]
func (co CustomerTypeController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customerTypes []databases.MedCustomerType
	co.DB.Where("is_active = ?", true).Find(&customerTypes)

	co.SetBody(customerTypes)
	return
}

// Get customerType
// @Summary Get customerType
// @Description Show customerType
// @Tags CustomerType
// @Accept json
// @Produce json
// @Param id path uint true "customerType ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedCustomerType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerType/{id} [get]
func (co CustomerTypeController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customerType databases.MedCustomerType
	result := co.DB.First(&customerType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(customerType)
	return
}

// Create customerType
// @Summary Create customerType
// @Description Add customerType
// @Tags CustomerType
// @Accept json
// @Produce json
// @Param customerType body form.CustomerTypeCreateParams true "customerType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerType [post]
func (co CustomerTypeController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CustomerTypeCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	colorCode := "#ebebeb"
	if params.ColorCode != "" {
		colorCode = params.ColorCode
	}

	authUser := co.GetAuth(c)
	customerType := databases.MedCustomerType{
		Name:         params.Name,
		IsActive:     params.IsActive,
		ColorCode:    colorCode,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&customerType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update customerType
// @Summary Update customerType
// @Description Edit customerType
// @Tags CustomerType
// @Accept json
// @Produce json
// @Param id path uint true "customerType ID"
// @Param customerType body form.CustomerTypeUpdateParams true "customerType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerType/{id} [put]
func (co CustomerTypeController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CustomerTypeUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var customerType databases.MedCustomerType
	result := co.DB.First(&customerType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	customerType.Name = params.Name
	customerType.ColorCode = params.ColorCode
	customerType.IsActive = params.IsActive
	customerType.Base.ModifiedDate = time.Now()
	customerType.ModifiedUser = &authUser

	result = co.DB.Save(&customerType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete customerType
// @Summary Delete customerType
// @Description Remove customerType
// @Tags CustomerType
// @Accept json
// @Produce json
// @Param customerType body form.DeleteParams true "customerTypecustomer"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerType/{id} [delete]
func (co CustomerTypeController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedCustomerType{}, v)
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
