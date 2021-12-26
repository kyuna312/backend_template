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

// StreetController struct
type StreetController struct {
	shared.BaseController
}

// ListStreets ...
type ListStreets struct {
	Total int64                 `json:"total"`
	List  []databases.RefStreet `json:"list"`
}

// Init Controller
func (co StreetController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)                          // List
	router.GET("/list/active/:district_id", co.ListActive) // ListActive
	router.GET("get/:id", co.Get)                          // Show
	router.POST("", co.Create)                             // Create
	router.PUT("/:id", co.Update)                          // Update
	router.DELETE("", co.Delete)                           // Delete
	router.GET("/initdb", co.InitDB)                       // InitDB
}

// List streets
// @Summary List streets
// @Description Get streets
// @Tags Street
// @Accept json
// @Produce json
// @Param filter body form.StreetFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefStreet}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /streets/list [post]
func (co StreetController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.StreetFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg

	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListStreets

	var countries []databases.RefStreet
	db.Preload("District").Preload("District.City").Preload("District.City.Country").Find(&countries)

	db.Table("ref_streets").Count(&count)

	listRepsonse.List = countries
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// ListActive streets
// @Summary ListActive streets
// @Description List streets
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefStreet}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /streets/list/active/{district_id} [get]
func (co StreetController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var counties []databases.RefStreet
	if c.Param("district_id") == "" {
		co.DB.Find(&counties)
	} else {
		co.DB.Where("district_id = ?", c.Param("district_id")).Find(&counties)
	}

	co.SetBody(counties)
	return
}

// Get streets
// @Summary Get streets
// @Description Show streets
// @Tags Street
// @Accept json
// @Produce json
// @Param id path uint true "streets ID"
// @Success 200 {object} structs.ResponseBody{body=databases.RefStreet}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /streets/{id} [get]
func (co StreetController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var streets databases.RefStreet
	result := co.DB.First(&streets, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(streets)
	return
}

// Create streets
// @Summary Create streets
// @Description Add streets
// @Tags Street
// @Accept json
// @Produce json
// @Param streets body form.StreetParams true "streets"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /streets [post]
func (co StreetController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.StreetParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	authUser := co.GetAuth(c)
	streets := databases.RefStreet{
		Name:         params.Name,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		DistrictID:   uint(params.DistrictID),
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&streets)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update streets
// @Summary Update streets
// @Description Edit streets
// @Tags Street
// @Accept json
// @Produce json
// @Param id path uint true "streets ID"
// @Param streets body form.StreetParams true "streets"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /streets/{id} [put]
func (co StreetController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.StreetParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var street databases.RefStreet
	result := co.DB.First(&street, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	street.Name = params.Name
	street.Base.ModifiedDate = time.Now()
	street.ModifiedUser = &authUser
	street.DistrictID = uint(params.DistrictID)

	result = co.DB.Save(&street)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete streets
// @Summary Delete streets
// @Description Remove streets
// @Tags Street
// @Accept json
// @Produce json
// @Param county body form.DeleteParams true "county"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /streets/{id} [delete]
func (co StreetController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.RefStreet{}, v)
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
func (co StreetController) InitDB(c *gin.Context) {
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
		streets := databases.RefStreet{
			Name:         eachLn,
			CreatedUser:  &authUser,
			ModifiedUser: &authUser,
			Base: databases.Base{
				CreatedDate: time.Now(),
			},
		}

		result := co.DB.Create(&streets)
		if result.Error != nil {
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}
	}
	return
}
