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

// ValuteController struct
type ValuteController struct {
	shared.BaseController
}

// ListValute ...
type ListValute struct {
	Total int64                 `json:"total"`
	List  []databases.RefValute `json:"list"`
}

// Init Controller
func (co ValuteController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)                  // List
	router.GET("get/:id", co.Get)                  // Show
	router.POST("", co.Create)                     // Create
	router.PUT("/:id", co.Update)                  // Update
	router.DELETE("", co.Delete)                   // Delete
	router.GET("/mongolBank/:name", co.MongolBank) //MongolBank
}

// List valute
// @Summary List valute
// @Description Get valute
// @Tags Valute
// @Accept json
// @Produce json
// @Param filter body form.ValuteFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.RefValute}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /valute/list [post]
func (co ValuteController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.ValuteFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListValute

	var valute []databases.RefValute
	db.Find(&valute)

	db.Table("ref_valutes").Count(&count)

	listRepsonse.List = valute
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return

}

// Get valute
// @Summary Get valute
// @Description Show valute
// @Tags Valute
// @Accept json
// @Produce json
// @Param id path uint true "valute ID"
// @Success 200 {object} structs.ResponseBody{body=databases.RefValute}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /valute/{id} [get]
func (co ValuteController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var valute databases.RefValute
	result := co.DB.First(&valute, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(valute)
	return
}

// Create valute
// @Summary Create valute
// @Description Add valute
// @Tags Valute
// @Accept json
// @Produce json
// @Param valute body form.ValuteCreateParams true "valute"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /valute [post]
func (co ValuteController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.ValuteCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	valute := databases.RefValute{
		Name:         params.Name,
		Description:  params.Description,
		IsActive:     params.IsActive,
		Symbol:       params.Symbol,
		CreatedUser:  &authUser,
		ModifiedUser: &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&valute)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update valute
// @Summary Update valute
// @Description Edit valute
// @Tags Valute
// @Accept json
// @Produce json
// @Param id path uint true "valute ID"
// @Param valute body form.ValuteUpdateParams true "valute"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /valute/{id} [put]
func (co ValuteController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.ValuteUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var valute databases.RefValute
	result := co.DB.First(&valute, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	valute.Name = params.Name
	valute.Description = params.Description
	valute.IsActive = params.IsActive
	valute.Symbol = params.Symbol
	valute.Base.ModifiedDate = time.Now()
	valute.ModifiedUser = &authUser

	result = co.DB.Save(&valute)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete valute
// @Summary Delete valute
// @Description Remove valute
// @Tags Valute
// @Accept json
// @Produce json
// @Param valute body form.DeleteParams true "valute"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /valute/{id} [delete]
func (co ValuteController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.RefValute{}, v)
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

// MongolBank valute
// @Summary MongolBank valute
// @Description MongolBank valute
// @Tags Valute
// @Accept json
// @Produce json
// @Param name path uint true "valute name"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /valute/mongolBank/{name} [get]
func (co ValuteController) MongolBank(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	type MongolBankResponse struct {
		LastDate  string  `json:"last_date"`
		Code      string  `json:"code"`
		Rate      string  `json:"rate"`
		RateFloat float64 `json:"rate_float"`
		Name      string  `json:"name"`
	}

	var repsonse []MongolBankResponse

	url := "http://monxansh.appspot.com/xansh.json?currency=" + c.Param("name")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if errUnmarshal := json.Unmarshal(body, &repsonse); err != nil {
		co.SetError(http.StatusBadRequest, errUnmarshal.Error())
		return
	}

	co.SetBody(repsonse)

	return
}
