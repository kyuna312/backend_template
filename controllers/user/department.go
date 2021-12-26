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

// DepartmentController struct
type DepartmentController struct {
	shared.BaseController
}

// ListDepartment ...
type ListDepartment struct {
	Total int64                        `json:"total"`
	List  []databases.MedHrmDepartment `json:"list"`
}

// Init Controller
func (co DepartmentController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List) // List
	router.GET("get/:id", co.Get) // Show
	router.POST("", co.Create)    // Create
	router.PUT("/:id", co.Update) // Update
	router.DELETE("", co.Delete)  // Delete
}

// List departments
// @Summary List departments
// @Description Get departments
// @Tags Department
// @Accept json
// @Produce json
// @Param filter body form.MeasureFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedHrmDepartment}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /department/list [post]
func (co DepartmentController) List(c *gin.Context) {
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

	var listRepsonse ListDepartment

	var departments []databases.MedHrmDepartment
	db.Find(&departments)

	db.Table("med_department").Count(&count)

	listRepsonse.List = departments
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get department
// @Summary Get department
// @Description Show department
// @Tags Department
// @Accept json
// @Produce json
// @Param id path uint true "department ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedHrmDepartment}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /department/{id} [get]
func (co DepartmentController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var departments databases.MedHrmDepartment
	result := co.DB.First(&departments, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(departments)
	return
}

// Create department
// @Summary Create department
// @Description Add department
// @Tags Department
// @Accept json
// @Produce json
// @Param department body form.MeasureCreateParams true "Department"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /department [post]
func (co DepartmentController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	departments := databases.MedHrmDepartment{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&departments)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update department
// @Summary Update department
// @Description Edit department
// @Tags Department
// @Accept json
// @Produce json
// @Param id path uint true "Department ID"
// @Param department body form.MeasureCreateParams true "Department"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /department/{id} [put]
func (co DepartmentController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.MeasureCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var department databases.MedHrmDepartment
	result := co.DB.First(&department, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	department.Name = params.Name
	department.Description = params.Description
	department.IsActive = params.IsActive
	department.Base.ModifiedDate = time.Now()
	department.ModifiedUser = &authUser

	result = co.DB.Save(&department)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete department
// @Summary Delete department
// @Description Remove department
// @Tags Department
// @Accept json
// @Produce json
// @Param department body form.DeleteParams true "department"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /department/{id} [delete]
func (co DepartmentController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedHrmDepartment{}, v)
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
