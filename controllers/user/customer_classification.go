package user

import (
	"net/http"
	"reflect"
	"time"

	gin "github.com/gin-gonic/gin"
	shared "gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// CustomerClassificationController struct
type CustomerClassificationController struct {
	shared.BaseController
}

// ListCustomerClassification ...
type ListCustomerClassification struct {
	Total int64                                 `json:"total"`
	List  []databases.MedCustomerClassification `json:"list"`
}

// Init Controller
func (co CustomerClassificationController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List customerClassification
// @Summary List customerClassification
// @Description Get customerClassification
// @Tags CustomerClassification
// @Accept json
// @Produce json
// @Param filter body form.ClassFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedCustomerClassification}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerClassification/list [POST]
func (co CustomerClassificationController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.ClassFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListCustomerClassification

	var customerClassifications []databases.MedCustomerClassification
	db.Find(&customerClassifications)

	co.DB.Table("med_customer_classifications").Count(&count)

	listRepsonse.List = customerClassifications
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return
}

// Get customerClassification
// @Summary Get customerClassification
// @Description Show customerClassification
// @Tags CustomerClassification
// @Accept json
// @Produce json
// @Param id path uint true "customerClassification ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedCustomerClassification}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerClassification/{id} [get]
func (co CustomerClassificationController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customerClassification databases.MedCustomerClassification
	result := co.DB.First(&customerClassification, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(customerClassification)
	return
}

// GetCode customerClassification
// @Summary GetCode customerClassification
// @Description Show customerClassification
// @Tags CustomerClassification
// @Accept json
// @Produce json
// @Param code path uint true "customerClassification code"
// @Success 200 {object} structs.ResponseBody{body=databases.MedCustomerClassification}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerClassification/code/{code} [get]
func (co CustomerClassificationController) GetCode(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var customerClassification databases.MedCustomerClassification
	result := co.DB.First(&customerClassification, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(customerClassification)
	return
}

// Create customerClassification
// @Summary Create customerClassification
// @Description Add customerClassification
// @Tags CustomerClassification
// @Accept json
// @Produce json
// @Param customerClassification body form.SpecificationCreateParams true "customerClassification"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerClassification [post]
func (co CustomerClassificationController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.SpecificationCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	authUser := co.GetAuth(c)

	customerClassification := databases.MedCustomerClassification{
		Name:         params.Name,
		IsActive:     params.IsActive,
		Description:  params.Description,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&customerClassification)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update customerClassification
// @Summary Update customerClassification
// @Description Edit customerClassification
// @Tags CustomerClassification
// @Accept json
// @Produce json
// @Param id path uint true "customerClassification ID"
// @Param customerClassification body form.SpecificationUpdateParams true "customerClassification"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerClassification/{id} [put]
func (co CustomerClassificationController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.SpecificationUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var customerClassification databases.MedCustomerClassification
	result := co.DB.First(&customerClassification, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	customerClassification.Name = params.Name
	customerClassification.IsActive = params.IsActive
	customerClassification.Description = params.Description
	customerClassification.Base.ModifiedDate = time.Now()
	customerClassification.ModifiedUser = &authUser

	result = co.DB.Save(&customerClassification)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete customerClassification
// @Summary Delete customerClassification
// @Description Remove customerClassification
// @Tags CustomerClassification
// @Accept json
// @Produce json
// @Param customerClassification body form.DeleteParams true "customerClassification"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /customerClassification [delete]
func (co CustomerClassificationController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedCustomerClassification{}, v)
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
