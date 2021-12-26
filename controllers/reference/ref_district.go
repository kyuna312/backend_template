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

// DistrictController struct
type DistrictController struct {
	shared.BaseController
}

// ListDistricts ...
type ListDistricts struct {
	Total int64                   `json:"total"`
	List  []databases.RefDistrict `json:"list"`
}

// Init Controller
func (co DistrictController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)                      // List
	router.GET("/list/active/:city_id", co.ListActive) // ListActive
	router.GET("get/:id", co.Get)                      // Show
	router.POST("", co.Create)                         // Create
	router.PUT("/:id", co.Update)                      // Update
	router.DELETE("", co.Delete)                       // Delete
	router.GET("/initdb", co.InitDB)                   // InitDB
}

// List districts
// @Summary List districts
// @Description Get districts
// @Tags District
// @Accept json
// @Produce json
// @Param filter body form.DistrictFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefDistrict}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /districts/list [post]
func (co DistrictController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.DistrictFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg

	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListDistricts

	var countries []databases.RefDistrict
	db.Preload("City").Preload("City.Country").Find(&countries)

	co.DB.Table("ref_districts").Count(&count)

	listRepsonse.List = countries
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// ListActive districts
// @Summary ListActive districts
// @Description List districts
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefDistrict}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /districts/list/active/{city_id} [get]
func (co DistrictController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var counties []databases.RefDistrict
	if c.Param("city_id") == "" {
		co.DB.Find(&counties)
	} else {
		co.DB.Where("city_id = ?", c.Param("city_id")).Find(&counties)
	}

	co.SetBody(counties)
	return
}

// Get districts
// @Summary Get districts
// @Description Show districts
// @Tags District
// @Accept json
// @Produce json
// @Param id path uint true "districts ID"
// @Success 200 {object} structs.ResponseBody{body=databases.RefDistrict}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /districts/{id} [get]
func (co DistrictController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var districts databases.RefDistrict
	result := co.DB.First(&districts, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(districts)
	return
}

// Create districts
// @Summary Create districts
// @Description Add districts
// @Tags District
// @Accept json
// @Produce json
// @Param districts body form.DistrictParams true "districts"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /districts [post]
func (co DistrictController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DistrictParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	authUser := co.GetAuth(c)
	districts := databases.RefDistrict{
		Name:         params.Name,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		CityID:       uint(params.CityID),
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&districts)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update districts
// @Summary Update districts
// @Description Edit districts
// @Tags District
// @Accept json
// @Produce json
// @Param id path uint true "districts ID"
// @Param districts body form.DistrictParams true "districts"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /districts/{id} [put]
func (co DistrictController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DistrictParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var district databases.RefDistrict
	result := co.DB.First(&district, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	district.Name = params.Name
	district.Base.ModifiedDate = time.Now()
	district.ModifiedUser = &authUser
	district.CityID = uint(params.CityID)

	result = co.DB.Save(&district)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete districts
// @Summary Delete districts
// @Description Remove districts
// @Tags District
// @Accept json
// @Produce json
// @Param county body form.DeleteParams true "county"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /districts/{id} [delete]
func (co DistrictController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.RefDistrict{}, v)
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
func (co DistrictController) InitDB(c *gin.Context) {
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
		districts := databases.RefDistrict{
			Name:         eachLn,
			CreatedUser:  &authUser,
			ModifiedUser: &authUser,
			Base: databases.Base{
				CreatedDate: time.Now(),
			},
		}

		result := co.DB.Create(&districts)
		if result.Error != nil {
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}
	}
	return
}
