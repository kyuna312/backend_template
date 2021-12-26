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

// PositionKeyController struct
type PositionKeyController struct {
	shared.BaseController
}

// ListPositionKey ...
type ListPositionKey struct {
	Total int64                         `json:"total"`
	List  []databases.MedHrmPositionKey `json:"list"`
}

// Init Controller
func (co PositionKeyController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List positionKeys
// @Summary List positionKeys
// @Description Get positionKeys
// @Tags PositionKey
// @Accept json
// @Produce json
// @Param filter body form.MeasureFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedHrmPositionKey}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionKey/list [post]
func (co PositionKeyController) List(c *gin.Context) {
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

	var listRepsonse ListPositionKey

	var positionKeys []databases.MedHrmPositionKey
	db.Find(&positionKeys)

	db.Table("med_positionKey").Count(&count)

	listRepsonse.List = positionKeys
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get positionKey
// @Summary Get positionKey
// @Description Show positionKey
// @Tags PositionKey
// @Accept json
// @Produce json
// @Param id path uint true "positionKey ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedHrmPositionKey}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionKey/{id} [get]
func (co PositionKeyController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var positionKeys databases.MedHrmPositionKey
	result := co.DB.First(&positionKeys, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(positionKeys)
	return
}

// Create positionKey
// @Summary Create positionKey
// @Description Add positionKey
// @Tags PositionKey
// @Accept json
// @Produce json
// @Param positionKey body form.MeasureCreateParams true "PositionKey"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionKey [post]
func (co PositionKeyController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	positionKeys := databases.MedHrmPositionKey{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&positionKeys)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update positionKey
// @Summary Update positionKey
// @Description Edit positionKey
// @Tags PositionKey
// @Accept json
// @Produce json
// @Param id path uint true "PositionKey ID"
// @Param positionKey body form.MeasureCreateParams true "PositionKey"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionKey/{id} [put]
func (co PositionKeyController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var positionKey databases.MedHrmPositionKey
	result := co.DB.First(&positionKey, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	positionKey.Name = params.Name
	positionKey.Description = params.Description
	positionKey.IsActive = params.IsActive
	positionKey.Base.ModifiedDate = time.Now()
	positionKey.ModifiedUser = &authUser

	result = co.DB.Save(&positionKey)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete positionKey
// @Summary Delete positionKey
// @Description Remove positionKey
// @Tags PositionKey
// @Accept json
// @Produce json
// @Param positionKey body form.DeleteParams true "positionKey"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /positionKey/{id} [delete]
func (co PositionKeyController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedHrmPositionKey{}, v)
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
