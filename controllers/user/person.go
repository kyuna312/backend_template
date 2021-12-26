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

// PersonController struct
type PersonController struct {
	shared.BaseController
}

// ListPersons ...
type ListPersons struct {
	Total int64                     `json:"total"`
	List  []databases.MedBasePerson `json:"list"`
}

// Init Controller
func (co PersonController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)    // List
	router.GET("get/:id", co.Get)    // Show
	router.POST("", co.Create)       // Create
	router.PUT("/:id", co.Update)    // Update
	router.DELETE("/:id", co.Delete) // Delete
	router.GET("/me", co.Me)         // Me
}

// List person
// @Summary List person
// @Description Get person
// @Tags Person
// @Accept json
// @Produce json
// @Param filter body form.PersonFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedBasePerson}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /person/list [post]
func (co PersonController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.PersonFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListPersons

	var persons []databases.MedBasePerson
	db.Find(&persons)

	db.Table("med_base_people").Count(&count)

	listRepsonse.List = persons
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return
}

// Get person
// @Summary Get person
// @Description Show person
// @Tags Person
// @Accept json
// @Produce json
// @Param id path uint true "person ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedBasePerson}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /person/{id} [get]
func (co PersonController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var person databases.MedBasePerson
	result := co.DB.First(&person, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(person)
	return
}

// Create person
// @Summary Create person
// @Description Add person
// @Tags Person
// @Accept json
// @Produce json
// @Param person body form.PersonCreateParams true "person"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /person [post]
func (co PersonController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PersonCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	person := databases.MedBasePerson{
		IsActive:       params.IsActive,
		LastName:       params.LastName,
		FirstName:      params.FirstName,
		StateRegNumber: params.StateRegNumber,
		MobileNumber:   params.MobileNumber,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&person)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update person
// @Summary Update person
// @Description Edit person
// @Tags Person
// @Accept json
// @Produce json
// @Param id path uint true "person ID"
// @Param person body form.PersonUpdateParams true "person"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /person/{id} [put]
func (co PersonController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PersonUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var person databases.MedBasePerson
	result := co.DB.First(&person, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	person.IsActive = params.IsActive
	person.LastName = params.LastName
	person.FirstName = params.FirstName
	person.StateRegNumber = params.StateRegNumber
	person.MobileNumber = params.MobileNumber
	// person.Username = params.Username
	// person.StartDate = params.StartDate
	// person.EndDate = params.EndDate

	person.Base.ModifiedDate = time.Now()

	result = co.DB.Save(&person)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete person
// @Summary Delete person
// @Description Remove person
// @Tags Person
// @Accept json
// @Produce json
// @Param person body form.DeleteParams true "person"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /person/{id} [delete]
func (co PersonController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedBasePerson{}, v)
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

// Me get auth person
// @Summary Get auth
// @Description Show auth
// @Tags Person
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=databases.MedBasePerson}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /person/me [get]
func (co PersonController) Me(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	co.SetBody(co.GetAuth(c))

	return
}
