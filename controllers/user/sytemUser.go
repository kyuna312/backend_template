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
	utils "gitlab.com/fibocloud/medtech/gin/utils"
)

// SystemUserController struct
type SystemUserController struct {
	shared.BaseController
}

// ListSystemUsers ...
type ListSystemUsers struct {
	Total int64                     `json:"total"`
	List  []databases.MedSystemUser `json:"list"`
}

// Init Controller
func (co SystemUserController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)    // List
	router.GET("get/:id", co.Get)    // Show
	router.POST("", co.Create)       // Create
	router.PUT("/:id", co.Update)    // Update
	router.DELETE("/:id", co.Delete) // Delete
	router.GET("/me", co.Me)         // Me
}

// List systemUser
// @Summary List systemUser
// @Description Get systemUser
// @Tags SystemUser
// @Accept json
// @Produce json
// @Param filter body form.SystemUserFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedSystemUser}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /systemUser/list [post]
func (co SystemUserController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.SystemUserFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListSystemUsers

	var systemUsers []databases.MedSystemUser
	db.Find(&systemUsers)

	db.Table("med_base_systemUser").Count(&count)

	listRepsonse.List = systemUsers
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return
}

// Get systemUser
// @Summary Get systemUser
// @Description Show systemUser
// @Tags SystemUser
// @Accept json
// @Produce json
// @Param id path uint true "systemUser ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedSystemUser}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /systemUser/{id} [get]
func (co SystemUserController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var systemUser databases.MedSystemUser
	result := co.DB.First(&systemUser, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(systemUser)
	return
}

// Create systemUser
// @Summary Create systemUser
// @Description Add systemUser
// @Tags SystemUser
// @Accept json
// @Produce json
// @Param systemUser body form.SystemUserParams true "systemUser"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /systemUser [post]
func (co SystemUserController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.SystemUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	hashPwd, err := utils.GenerateHash(params.Password)
	if err != nil {
		co.SetError(http.StatusInternalServerError, err.Error())
		return
	}

	systemUser := databases.MedSystemUser{
		IsActive:     params.IsActive,
		Username:     params.Username,
		StartDate:    params.StartDate,
		EndDate:      params.EndDate,
		PersonID:     params.PersondID,
		PersonType:   1,
		PasswordHash: hashPwd,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&systemUser)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update systemUser
// @Summary Update systemUser
// @Description Edit systemUser
// @Tags SystemUser
// @Accept json
// @Produce json
// @Param id path uint true "systemUser ID"
// @Param systemUser body form.SystemUserParams true "systemUser"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /systemUser/{id} [put]
func (co SystemUserController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.SystemUserParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var systemUser databases.MedSystemUser
	result := co.DB.First(&systemUser, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	systemUser.IsActive = params.IsActive
	systemUser.Username = params.Username
	systemUser.StartDate = params.StartDate
	systemUser.EndDate = params.EndDate

	systemUser.Base.ModifiedDate = time.Now()

	result = co.DB.Save(&systemUser)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete systemUser
// @Summary Delete systemUser
// @Description Remove systemUser
// @Tags SystemUser
// @Accept json
// @Produce json
// @Param systemUser body form.DeleteParams true "systemUser"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /systemUser/{id} [delete]
func (co SystemUserController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedSystemUser{}, v)
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

// Me get auth systemUser
// @Summary Get auth
// @Description Show auth
// @Tags SystemUser
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=databases.MedSystemUser}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /systemUser/me [get]
func (co SystemUserController) Me(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	co.SetBody(co.GetAuth(c))

	return
}
