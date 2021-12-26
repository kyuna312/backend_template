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

// StatusTypeController struct
type StatusTypeController struct {
	shared.BaseController
}

// ListStatusType ...
type ListStatusType struct {
	Total int64                     `json:"total"`
	List  []databases.MedStatusType `json:"list"`
}

// Init Controller
func (co StatusTypeController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)             // List
	router.GET("/list/active", co.ListActive) // ListActive
	router.GET("get/:id", co.Get)             // Show
	router.POST("", co.Create)                // Create
	router.PUT("/:id", co.Update)             // Update
	router.DELETE("", co.Delete)              // Delete
}

// List statusType
// @Summary List statusType
// @Description Get statusType
// @Tags StatusType
// @Accept json
// @Produce json
// @Param filter body form.MeasureFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedStatusType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /statusType/list [post]
func (co StatusTypeController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.MeasureFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListStatusType

	var statusType []databases.MedStatusType
	db.Find(&statusType)

	db.Table("med_status_types").Count(&count)

	listRepsonse.List = statusType
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// ListActive statusType
// @Summary ListActive statusType
// @Description Get statusType
// @Tags StatusType
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedStatusType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /statusType/list/active [get]
func (co StatusTypeController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var statusTypes []databases.MedStatusType
	co.DB.Where("is_active = ?", true).Preload("Status").Find(&statusTypes)

	co.SetBody(statusTypes)
	return
}

// Get statusType
// @Summary Get statusType
// @Description Show statusType
// @Tags StatusType
// @Accept json
// @Produce json
// @Param id path uint true "statusType ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedStatusType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /statusType/{id} [get]
func (co StatusTypeController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var statusType databases.MedStatusType
	result := co.DB.First(&statusType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(statusType)
	return
}

// Create statusType
// @Summary Create statusType
// @Description Add statusType
// @Tags StatusType
// @Accept json
// @Produce json
// @Param statusType body form.MeasureCreateParams true "statusType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /statusType [post]
func (co StatusTypeController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	statusType := databases.MedStatusType{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&statusType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update statusType
// @Summary Update statusType
// @Description Edit statusType
// @Tags StatusType
// @Accept json
// @Produce json
// @Param id path uint true "statusType ID"
// @Param statusType body form.MeasureUpdateParams true "statusType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /statusType/{id} [put]
func (co StatusTypeController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var statusType databases.MedStatusType
	result := co.DB.First(&statusType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	statusType.Name = params.Name
	statusType.Description = params.Description
	statusType.IsActive = params.IsActive
	statusType.Base.ModifiedDate = time.Now()
	statusType.ModifiedUser = &authUser

	result = co.DB.Save(&statusType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete statusType
// @Summary Delete statusType
// @Description Remove statusType
// @Tags StatusType
// @Accept json
// @Produce json
// @Param statusType body form.DeleteParams true "statusType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /statusType/{id} [delete]
func (co StatusTypeController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedStatusType{}, v)
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
