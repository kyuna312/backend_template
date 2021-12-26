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

// PermissionController struct
type PermissionController struct {
	shared.BaseController
}

// ListPermission ...
type ListPermission struct {
	Total int64                           `json:"total"`
	List  []databases.MedSystemPermission `json:"list"`
}

// Init Controller
func (co PermissionController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List permissions
// @Summary List permissions
// @Description Get permissions
// @Tags Permission
// @Accept json
// @Produce json
// @Param filter body form.MeasureFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedSystemPermission}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /permission/list [post]
func (co PermissionController) List(c *gin.Context) {
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

	var listRepsonse ListPermission

	var permissions []databases.MedSystemPermission
	db.Find(&permissions)

	db.Table("med_system_permissions").Count(&count)

	listRepsonse.List = permissions
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get permission
// @Summary Get permission
// @Description Show permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path uint true "permission ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedSystemPermission}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /permission/{id} [get]
func (co PermissionController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var permissions databases.MedSystemPermission
	result := co.DB.First(&permissions, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(permissions)
	return
}

// Create permission
// @Summary Create permission
// @Description Add permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param permission body form.PermissionParams true "Permission"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /permission [post]
func (co PermissionController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PermissionParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	permissions := databases.MedSystemPermission{
		Path:         params.Path,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&permissions)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update permission
// @Summary Update permission
// @Description Edit permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path uint true "Permission ID"
// @Param permission body form.PermissionParams true "Permission"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /permission/{id} [put]
func (co PermissionController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PermissionParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var permission databases.MedSystemPermission
	result := co.DB.First(&permission, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	permission.Path = params.Path
	permission.Description = params.Description
	permission.IsActive = params.IsActive
	permission.Base.ModifiedDate = time.Now()
	permission.ModifiedUser = &authUser

	result = co.DB.Save(&permission)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete permission
// @Summary Delete permission
// @Description Remove permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param permission body form.DeleteParams true "permission"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /permission/{id} [delete]
func (co PermissionController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedSystemPermission{}, v)
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
