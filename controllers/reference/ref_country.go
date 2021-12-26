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

// CountryController struct
type CountryController struct {
	shared.BaseController
}

// ListCountries ...
type ListCountries struct {
	Total int64                  `json:"total"`
	List  []databases.RefCountry `json:"list"`
}

// Init Controller
func (co CountryController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)             // List
	router.GET("/list/active", co.ListActive) // ListActive
	router.GET("get/:id", co.Get)             // Show
	router.POST("", co.Create)                // Create
	router.PUT("/:id", co.Update)             // Update
	router.DELETE("", co.Delete)              // Delete
	router.GET("/initdb", co.InitDB)          // InitDB
}

// List country
// @Summary List country
// @Description Get country
// @Tags Country
// @Accept json
// @Produce json
// @Param filter body form.CountryFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefCountry}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /country/list [post]
func (co CountryController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.CountryFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg

	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListCountries

	var countries []databases.RefCountry
	db.Find(&countries)

	db.Table("ref_countries").Count(&count)

	listRepsonse.List = countries
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// ListActive country
// @Summary ListActive country
// @Description List country
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefCountry}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /country/list/active [get]
func (co CountryController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var counties []databases.RefCountry
	co.DB.Find(&counties)
	co.SetBody(counties)
	return
}

// Get country
// @Summary Get country
// @Description Show country
// @Tags Country
// @Accept json
// @Produce json
// @Param id path uint true "country ID"
// @Success 200 {object} structs.ResponseBody{body=databases.RefCountry}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /country/{id} [get]
func (co CountryController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var country databases.RefCountry
	result := co.DB.First(&country, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(country)
	return
}

// Create country
// @Summary Create country
// @Description Add country
// @Tags Country
// @Accept json
// @Produce json
// @Param country body form.CountryParams true "country"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /country [post]
func (co CountryController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CountryParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	authUser := co.GetAuth(c)
	country := databases.RefCountry{
		Name:         params.Name,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&country)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update country
// @Summary Update country
// @Description Edit country
// @Tags Country
// @Accept json
// @Produce json
// @Param id path uint true "country ID"
// @Param country body form.CountryParams true "country"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /country/{id} [put]
func (co CountryController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CountryParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var country databases.RefCountry
	result := co.DB.First(&country, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	country.Name = params.Name
	country.Base.ModifiedDate = time.Now()
	country.ModifiedUser = &authUser

	result = co.DB.Save(&country)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete country
// @Summary Delete country
// @Description Remove country
// @Tags Country
// @Accept json
// @Produce json
// @Param county body form.DeleteParams true "county"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /country/{id} [delete]
func (co CountryController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.RefCountry{}, v)
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
func (co CountryController) InitDB(c *gin.Context) {
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
		country := databases.RefCountry{
			Name:         eachLn,
			CreatedUser:  &authUser,
			ModifiedUser: &authUser,
			Base: databases.Base{
				CreatedDate: time.Now(),
			},
		}

		result := co.DB.Create(&country)
		if result.Error != nil {
			co.SetError(http.StatusInternalServerError, result.Error.Error())
			return
		}
	}
	return
}
