package reference

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	gin "github.com/gin-gonic/gin"
	"gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// AddressTypeController struct
type AddressTypeController struct {
	shared.BaseController
}

// ListAddressTypes ...
type ListAddressTypes struct {
	Total int64                      `json:"total"`
	List  []databases.MedAddressType `json:"list"`
}

// Init Controller
func (co AddressTypeController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)             // List
	router.GET("get/:id", co.Get)             // Show
	router.GET("/list/active", co.ListActive) // ListActive
	router.POST("", co.Create)                // Create
	router.PUT("/:id", co.Update)             // Update
	router.DELETE("", co.Delete)              // Delete
	router.GET("/initdb", co.InitDB)          // InitDB
}

// List city
// @Summary List city
// @Description Get city
// @Tags City
// @Accept json
// @Produce json
// @Param filter body form.CityFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedAddressType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/list [post]
func (co AddressTypeController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.CityFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg

	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListAddressTypes

	var types []databases.MedAddressType
	db.Find(&types)

	db.Table("ref_cities").Count(&count)

	listRepsonse.List = types
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// ListActive city
// @Summary ListActive city
// @Description ListActive city
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedAddressType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/list/active [get]
func (co AddressTypeController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var cities []databases.MedAddressType
	co.DB.Find(&cities)
	co.SetBody(cities)
	return
}

// Get city
// @Summary Get city
// @Description Show city
// @Tags City
// @Accept json
// @Produce json
// @Param id path uint true "city ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedAddressType}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/{id} [get]
func (co AddressTypeController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var city databases.MedAddressType
	result := co.DB.First(&city, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(city)
	return
}

// Create city
// @Summary Create city
// @Description Add city
// @Tags City
// @Accept json
// @Produce json
// @Param city body form.CityParams true "city"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city [post]
func (co AddressTypeController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	authUser := co.GetAuth(c)
	city := databases.MedAddressType{
		Name:         params.Name,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		IsActive:     true,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&city)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update city
// @Summary Update city
// @Description Edit city
// @Tags City
// @Accept json
// @Produce json
// @Param id path uint true "city ID"
// @Param city body form.CityParams true "city"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/{id} [put]
func (co AddressTypeController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var city databases.MedAddressType
	result := co.DB.First(&city, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	city.Name = params.Name
	city.Base.ModifiedDate = time.Now()
	city.ModifiedUser = &authUser

	result = co.DB.Save(&city)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete city
// @Summary Delete city
// @Description Remove city
// @Tags City
// @Accept json
// @Produce json
// @Param county body form.DeleteParams true "county"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/{id} [delete]
func (co AddressTypeController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedAddressType{}, v)
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

// InitDB classification
// @Summary Delete classification
// @Description Remove classification
// @Tags Classification
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /initdb [get]
func (co AddressTypeController) InitDB(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	authUser := co.GetAuth(c)
	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.
	file, err := os.Open("files/countries.txt")

	if err != nil {
		log.Fatalf("failed to open")

	}

	// The bufio.NewScanner() function is called in which the
	// object os.File passed as its parameter and this returns a
	// object bufio.Scanner which is further used on the
	// bufio.Scanner.Split() method.
	scanner := bufio.NewScanner(file)

	// The bufio.ScanLines is used as an
	// input to the method bufio.Scanner.Split()
	// and then the scanning forwards to each
	// new line using the bufio.Scanner.Scan()
	// method.
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	// The method os.File.Close() is called
	// on the os.File object to close the file
	file.Close()

	// and then a loop iterates through
	// and prints each of the slice values.
	for _, eachLn := range text {
		city := databases.MedAddressType{
			Name:         eachLn,
			CreatedUser:  &authUser,
			ModifiedUser: &authUser,
			Base: databases.Base{
				CreatedDate: time.Now(),
			},
		}

		result := co.DB.Create(&city)
		if result.Error != nil {
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}
	}
	return
}
