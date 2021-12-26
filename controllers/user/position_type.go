package user

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	shared "gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// PositionTypeController struct
type PositionTypeController struct {
	shared.BaseController
}

// ListPositionType ...
type ListPositionType struct {
	Total int64                          `json:"total"`
	List  []databases.MedHrmPositionType `json:"list"`
}

// Init Controller
func (co PositionTypeController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)             // List
	router.GET("get/:id", co.Get)             // Show
	router.POST("", co.Create)                // Create
	router.PUT("/:id", co.Update)             // Update
	router.DELETE("", co.Delete)              // Delete
	router.GET("/list/active", co.ListActive) // ListActive
}

// List positionTypes
// @Summary List positionTypes
// @Description Get positionTypes
// @Tags PositionType
// @Accept json
// @Produce json
// @Param filter body form.MeasureFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedHrmPositionType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionType/list [post]
func (co PositionTypeController) List(c *gin.Context) {
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

	var listRepsonse ListPositionType

	var positionTypes []databases.MedHrmPositionType
	db.Find(&positionTypes)

	db.Table("med_hrm_position_types").Count(&count)

	listRepsonse.List = positionTypes
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return
}

// ListActive positionType
// @Summary ListActive positionType
// @Description ListActive positionType
// @Tags PositionType
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedHrmPositionType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionType/list/active [get]
func (co PositionTypeController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var positionTypes []databases.MedHrmPositionType
	co.DB.Where("is_active = ?", true).Find(&positionTypes)
	co.SetBody(positionTypes)
	return
}

// Get positionType
// @Summary Get positionType
// @Description Show positionType
// @Tags PositionType
// @Accept json
// @Produce json
// @Param id path uint true "positionType ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedHrmPositionType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionType/{id} [get]
func (co PositionTypeController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var positionTypes databases.MedHrmPositionType
	result := co.DB.First(&positionTypes, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(positionTypes)
	return
}

// Create positionType
// @Summary Create positionType
// @Description Add positionType
// @Tags PositionType
// @Accept json
// @Produce json
// @Param positionType body form.MeasureCreateParams true "PositionType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionType [post]
func (co PositionTypeController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var prePositionType databases.MedHrmPositionType
	resultPre := co.DB.Last(&prePositionType)
	newCode := ""

	if resultPre.Error != nil {
		newCode = "001"
	} else {
		intCode, intErr := strconv.Atoi(prePositionType.Code)
		if intErr != nil {
			co.SetError(http.StatusBadRequest, intErr.Error())
			return
		}
		s := fmt.Sprintf("%03d", intCode+1)
		newCode = s
	}

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	positionTypes := databases.MedHrmPositionType{
		Name:         params.Name,
		Code:         newCode,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&positionTypes)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update positionType
// @Summary Update positionType
// @Description Edit positionType
// @Tags PositionType
// @Accept json
// @Produce json
// @Param id path uint true "PositionType ID"
// @Param positionType body form.MeasureCreateParams true "PositionType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionType/{id} [put]
func (co PositionTypeController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var positionType databases.MedHrmPositionType
	result := co.DB.First(&positionType, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	positionType.Name = params.Name
	positionType.Description = params.Description
	positionType.IsActive = params.IsActive
	positionType.Base.ModifiedDate = time.Now()
	positionType.ModifiedUser = &authUser

	result = co.DB.Save(&positionType)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete positionType
// @Summary Delete positionType
// @Description Remove positionType
// @Tags PositionType
// @Accept json
// @Produce json
// @Param positionType body form.DeleteParams true "positionType"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionType/{id} [delete]
func (co PositionTypeController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedHrmPositionType{}, v)
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
