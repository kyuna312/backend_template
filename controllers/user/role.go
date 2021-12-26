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

// RoleController struct
type RoleController struct {
	shared.BaseController
}

// ListRole ...
type ListRole struct {
	Total int64                     `json:"total"`
	List  []databases.MedSystemRole `json:"list"`
}

// Init Controller
func (co RoleController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List roles
// @Summary List roles
// @Description Get roles
// @Tags Role
// @Accept json
// @Produce json
// @Param filter body form.MeasureFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedSystemRole}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /role/list [post]
func (co RoleController) List(c *gin.Context) {
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

	var listRepsonse ListRole

	var roles []databases.MedSystemRole
	db.Find(&roles)

	db.Table("med_role").Count(&count)

	listRepsonse.List = roles
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get role
// @Summary Get role
// @Description Show role
// @Tags Role
// @Accept json
// @Produce json
// @Param id path uint true "role ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedSystemRole}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /role/{id} [get]
func (co RoleController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var roles databases.MedSystemRole
	result := co.DB.First(&roles, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(roles)
	return
}

// Create role
// @Summary Create role
// @Description Add role
// @Tags Role
// @Accept json
// @Produce json
// @Param role body form.RoleCreateParams true "Role"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /role [post]
func (co RoleController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.RoleCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	roles := databases.MedSystemRole{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&roles)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update role
// @Summary Update role
// @Description Edit role
// @Tags Role
// @Accept json
// @Produce json
// @Param id path uint true "Role ID"
// @Param role body form.RoleUpdateParams true "Role"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /role/{id} [put]
func (co RoleController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.RoleUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var role databases.MedSystemRole
	result := co.DB.First(&role, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	role.Name = params.Name
	role.Description = params.Description
	role.IsActive = params.IsActive
	role.Base.ModifiedDate = time.Now()
	role.ModifiedUser = &authUser
	role.Permissions = params.Permissions

	result = co.DB.Save(&role)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete role
// @Summary Delete role
// @Description Remove role
// @Tags Role
// @Accept json
// @Produce json
// @Param role body form.DeleteParams true "role"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /role/{id} [delete]
func (co RoleController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedSystemRole{}, v)
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
