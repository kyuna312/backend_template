package reference

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	gin "github.com/gin-gonic/gin"
	"gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// CityController struct
type CityController struct {
	shared.BaseController
}

// RegionStruct ...
type RegionStruct struct {
	ID           int         `json:"id"`
	RegionName   string      `json:"region_name"`
	RegionNameEn interface{} `json:"region_name_en"`
	RegionTypeID int         `json:"region_type_id"`
	ParentID     int         `json:"parent_id"`
	Order        int         `json:"order"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	IsDeleted    bool        `json:"is_deleted"`
	CreateID     int         `json:"create_id"`
	UpdateID     int         `json:"update_id"`
}

// ListCities ...
type ListCities struct {
	Total int64               `json:"total"`
	List  []databases.RefCity `json:"list"`
}

// Init Controller
func (co CityController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)                        // List
	router.GET("get/:id", co.Get)                        // Show
	router.GET("/list/active/:county_id", co.ListActive) // ListActive
	router.POST("", co.Create)                           // Create
	router.PUT("/:id", co.Update)                        // Update
	router.DELETE("", co.Delete)                         // Delete
	router.GET("/initdb", co.InitDB)                     // InitDB
}

// List city
// @Summary List city
// @Description Get city
// @Tags City
// @Accept json
// @Produce json
// @Param filter body form.CityFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefCity}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/list [post]
func (co CityController) List(c *gin.Context) {
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

	var listRepsonse ListCities

	var cities []databases.RefCity
	db.Preload("Country").Find(&cities)

	db.Table("ref_cities").Count(&count)

	listRepsonse.List = cities
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
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefCity}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/list/active/{county_id} [get]
func (co CityController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var cities []databases.RefCity

	if c.Param("county_id") == "" {
		co.DB.Find(&cities)
	} else {
		co.DB.Where("country_id = ?", c.Param("county_id")).Find(&cities)
	}

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
// @Success 200 {object} structs.ResponseBody{body=databases.RefCity}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /city/{id} [get]
func (co CityController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var city databases.RefCity
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
func (co CityController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	authUser := co.GetAuth(c)
	city := databases.RefCity{
		Name:         params.Name,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		CountryID:    params.CountryID,
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
func (co CityController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.CityParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var city databases.RefCity
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
func (co CityController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.RefCity{}, v)
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
func (co CityController) InitDB(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	file, _ := ioutil.ReadFile("files/region.json")
	authUser := co.GetAuth(c)

	var data []RegionStruct

	_ = json.Unmarshal([]byte(file), &data)
	for _, region := range data {
		if region.RegionTypeID == 39 {
			// aiamg = hot
			city := databases.RefCity{
				Name:         region.RegionName,
				CreatedUser:  &authUser,
				ModifiedUser: &authUser,
				CountryID:    67,
				Base: databases.Base{
					ID:          uint(region.ID),
					CreatedDate: time.Now(),
				},
			}

			result := co.DB.Create(&city)
			if result.Error != nil {
				co.SetError(http.StatusInternalServerError, result.Error.Error())
				return
			}
		}

		if region.RegionTypeID == 58 {
			// aiamg = hot
			district := databases.RefDistrict{
				Name:         region.RegionName,
				CreatedUser:  &authUser,
				ModifiedUser: &authUser,
				CityID:       uint(region.ParentID),
				Base: databases.Base{
					ID:          uint(region.ID),
					CreatedDate: time.Now(),
				},
			}

			result := co.DB.Create(&district)
			if result.Error != nil {
				co.SetError(http.StatusInternalServerError, result.Error.Error())
				return
			}
		}

		if region.RegionTypeID == 102 {
			// aiamg = hot
			street := databases.RefStreet{
				Name:         region.RegionName,
				CreatedUser:  &authUser,
				ModifiedUser: &authUser,
				DistrictID:   uint(region.ParentID),
				Base: databases.Base{
					ID:          uint(region.ID),
					CreatedDate: time.Now(),
				},
			}

			result := co.DB.Create(&street)
			if result.Error != nil {
				co.SetError(http.StatusInternalServerError, result.Error.Error())
				return
			}
		}
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}
