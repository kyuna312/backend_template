package reference

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

// PaymentTypeController struct
type PaymentTypeController struct {
	shared.BaseController
}

// ListPaymentType ...
type ListPaymentType struct {
	Total int64                      `json:"total"`
	List  []databases.MedPaymentType `json:"list"`
}

// Init Controller
func (co PaymentTypeController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List paymentType
// @Summary List paymentType
// @Description Get paymentType
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param filter body form.PaymentTypeFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedPaymentType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentType/list [post]
func (co PaymentTypeController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.PaymentTypeFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListPaymentType

	var paymentType []databases.MedPaymentType
	db.Find(&paymentType)

	db.Table("med_payment_types").Count(&count)

	listRepsonse.List = paymentType
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get paymentType
// @Summary Get paymentType
// @Description Show paymentType
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param id path uint true "paymentType ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedPaymentType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentType/{id} [get]
func (co PaymentTypeController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var paymentType databases.MedPaymentType
	result := co.DB.First(&paymentType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(paymentType)
	return
}

// Create paymentType
// @Summary Create paymentType
// @Description Add paymentType
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param paymentType body form.PaymentTypeParams true "paymentType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentType [post]
func (co PaymentTypeController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PaymentTypeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	paymentType := databases.MedPaymentType{
		Name:             params.Name,
		Description:      params.Description,
		IsActive:         params.IsActive,
		CreatedUser:      &authUser,
		ModifiedUser:     &authUser,
		PaymentDay:       uint(params.PaymentDay),
		PrepaidPercent:   uint(params.PrepaidPercent),
		PaymentCondition: params.PaymentCondition,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&paymentType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update paymentType
// @Summary Update paymentType
// @Description Edit paymentType
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param id path uint true "paymentType ID"
// @Param paymentType body form.PaymentTypeParams true "paymentType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentType/{id} [put]
func (co PaymentTypeController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PaymentTypeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var paymentType databases.MedPaymentType
	result := co.DB.First(&paymentType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	paymentType.Name = params.Name
	paymentType.Description = params.Description
	paymentType.IsActive = params.IsActive
	paymentType.Base.ModifiedDate = time.Now()
	paymentType.ModifiedUser = &authUser
	paymentType.PaymentDay = uint(params.PaymentDay)
	paymentType.PrepaidPercent = uint(params.PrepaidPercent)
	paymentType.PaymentCondition = params.PaymentCondition

	result = co.DB.Save(&paymentType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete paymentType
// @Summary Delete paymentType
// @Description Remove paymentType
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param paymentType body form.DeleteParams true "paymentType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentType/{id} [delete]
func (co PaymentTypeController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedPaymentType{}, v)
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
