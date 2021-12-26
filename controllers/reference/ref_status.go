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

// StatusController struct
type StatusController struct {
	shared.BaseController
}

// ListStatus ...
type ListStatus struct {
	Total int64                 `json:"total"`
	List  []databases.MedStatus `json:"list"`
}

// Init Controller
func (co StatusController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List status
// @Summary List status
// @Description Get status
// @Tags Status
// @Accept json
// @Produce json
// @Param filter body form.StatusFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedStatus}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /status/list [post]
func (co StatusController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.StatusFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListStatus

	var status []databases.MedStatus
	db.Preload("StatusType").Find(&status)
	db.Table("med_statuses").Count(&count)

	listRepsonse.List = status
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get status
// @Summary Get status
// @Description Show status
// @Tags Status
// @Accept json
// @Produce json
// @Param id path uint true "status ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedStatus}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /status/{id} [get]
func (co StatusController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var status databases.MedStatus
	result := co.DB.First(&status, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(status)
	return
}

// Create status
// @Summary Create status
// @Description Add status
// @Tags Status
// @Accept json
// @Produce json
// @Param status body form.StatusParams true "status"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /status [post]
func (co StatusController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.StatusParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	colorCode := "#ebebeb"
	if params.ColorCode != "" {
		colorCode = params.ColorCode
	}

	authUser := co.GetAuth(c)
	status := databases.MedStatus{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		ColorCode:    colorCode,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		StatusTypeID: uint(params.TypeID),
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&status)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update status
// @Summary Update status
// @Description Edit status
// @Tags Status
// @Accept json
// @Produce json
// @Param id path uint true "status ID"
// @Param status body form.StatusParams true "status"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /status/{id} [put]
func (co StatusController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.StatusParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var status databases.MedStatus
	result := co.DB.First(&status, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	status.Name = params.Name
	status.Description = params.Description
	status.IsActive = params.IsActive
	status.ColorCode = params.ColorCode
	status.Base.ModifiedDate = time.Now()
	status.ModifiedUser = &authUser
	status.StatusTypeID = uint(params.TypeID)

	result = co.DB.Save(&status)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete status
// @Summary Delete status
// @Description Remove status
// @Tags Status
// @Accept json
// @Produce json
// @Param status body form.DeleteParams true "status"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /status/{id} [delete]
func (co StatusController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedStatus{}, v)
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
