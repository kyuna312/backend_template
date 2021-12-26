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

// PaymentMethodController struct
type PaymentMethodController struct {
	shared.BaseController
}

// ListPaymentMethod ...
type ListPaymentMethod struct {
	Total int64                        `json:"total"`
	List  []databases.MedPaymentMethod `json:"list"`
}

// Init Controller
func (co PaymentMethodController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List paymentMethod
// @Summary List paymentMethod
// @Description Get paymentMethod
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param filter body form.PaymentTypeFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedPaymentMethod}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentMethod/list [post]
func (co PaymentMethodController) List(c *gin.Context) {
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

	var listRepsonse ListPaymentMethod

	var paymentMethod []databases.MedPaymentMethod
	db.Find(&paymentMethod)

	db.Table("med_payment_methods").Count(&count)

	listRepsonse.List = paymentMethod
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get paymentMethod
// @Summary Get paymentMethod
// @Description Show paymentMethod
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param id path uint true "paymentMethod ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedPaymentMethod}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentMethod/{id} [get]
func (co PaymentMethodController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var paymentMethod databases.MedPaymentMethod
	result := co.DB.First(&paymentMethod, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(paymentMethod)
	return
}

// Create paymentMethod
// @Summary Create paymentMethod
// @Description Add paymentMethod
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param paymentMethod body form.PaymentTypeParams true "paymentMethod"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentMethod [post]
func (co PaymentMethodController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PaymentTypeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	paymentMethod := databases.MedPaymentMethod{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&paymentMethod)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update paymentMethod
// @Summary Update paymentMethod
// @Description Edit paymentMethod
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param id path uint true "paymentMethod ID"
// @Param paymentMethod body form.PaymentTypeParams true "paymentMethod"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentMethod/{id} [put]
func (co PaymentMethodController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PaymentTypeParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var paymentMethod databases.MedPaymentMethod
	result := co.DB.First(&paymentMethod, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	paymentMethod.Name = params.Name
	paymentMethod.Description = params.Description
	paymentMethod.IsActive = params.IsActive
	paymentMethod.Base.ModifiedDate = time.Now()
	paymentMethod.ModifiedUser = &authUser

	result = co.DB.Save(&paymentMethod)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete paymentMethod
// @Summary Delete paymentMethod
// @Description Remove paymentMethod
// @Tags PaymentType
// @Accept json
// @Produce json
// @Param paymentMethod body form.DeleteParams true "paymentMethod"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /paymentMethod/{id} [delete]
func (co PaymentMethodController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedPaymentMethod{}, v)
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
