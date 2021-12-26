package user

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	shared "gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	form "gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// PositionController struct
type PositionController struct {
	shared.BaseController
}

// ListPosition ...
type ListPosition struct {
	Total int64                      `json:"total"`
	List  []databases.MedHrmPosition `json:"list"`
}

// Init Controller
func (co PositionController) Init(router *gin.RouterGroup) {
	router.POST("/list", co.List)                 // List
	router.GET("get/:id", co.Get)                 // Show
	router.POST("", co.Create)                    // Create
	router.PUT("/:id", co.Update)                 // Update
	router.DELETE("", co.Delete)                  // Delete
	router.GET("/list/active/:id", co.ListActive) // ListActive
}

// List positions
// @Summary List positions
// @Description Get positions
// @Tags Position
// @Accept json
// @Produce json
// @Param filter body form.PositionFilter true "filter"
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedHrmPosition}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /position/list [post]
func (co PositionController) List(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var count int64
	var params form.PositionFilter
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	db := co.DB

	// filter hiij bgaa heseg
	v := reflect.ValueOf(params.Filter)

	db = db.Scopes(shared.TableSearch(v, params.Sort))
	db = db.Scopes(shared.Paginate(params.Page, params.Size))

	var listRepsonse ListPosition

	var positions []databases.MedHrmPosition
	db.Preload("PositionType").Find(&positions)

	db.Table("med_hrm_positions").Count(&count)

	listRepsonse.List = positions
	listRepsonse.Total = count

	co.SetBody(listRepsonse)
	return
}

// ListActive position
// @Summary ListActive position
// @Description ListActive position
// @Tags Position
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=[]databases.MedHrmPosition}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /position/list/active/{id} [get]
func (co PositionController) ListActive(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var positions []databases.MedHrmPosition
	co.DB.Where("is_active = ?", true).Where("position_type_id = ?", c.Param("id")).Find(&positions)
	co.SetBody(positions)
	return
}

// Get position
// @Summary Get position
// @Description Show position
// @Tags Position
// @Accept json
// @Produce json
// @Param id path uint true "position ID"
// @Success 200 {object} structs.ResponseBody{body=databases.MedHrmPosition}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /position/{id} [get]
func (co PositionController) Get(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var positions databases.MedHrmPosition
	result := co.DB.Preload("PositionType").First(&positions, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(positions)
	return
}

// Create position
// @Summary Create position
// @Description Add position
// @Tags Position
// @Accept json
// @Produce json
// @Param position body form.PositionCreateParams true "Position"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /position [post]
func (co PositionController) Create(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var prePosition databases.MedHrmPosition
	resultPre := co.DB.Last(&prePosition)
	newCode := ""

	if resultPre.Error != nil {
		newCode = "001"
	} else {
		intCode, intErr := strconv.Atoi(prePosition.Code)
		if intErr != nil {
			co.SetError(http.StatusBadRequest, intErr.Error())
			return
		}
		s := fmt.Sprintf("%03d", intCode+1)
		newCode = s
	}

	var params form.PositionCreateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}
	authUser := co.GetAuth(c)
	positions := databases.MedHrmPosition{
		Name:           params.Name,
		Code:           newCode,
		Description:    params.Description,
		PositionTypeID: uint(params.TypeID),
		IsActive:       params.IsActive,
		CreatedUser:    &authUser,
		ModifiedUser:   &authUser,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&positions)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Update position
// @Summary Update position
// @Description Edit position
// @Tags Position
// @Accept json
// @Produce json
// @Param id path uint true "Position ID"
// @Param position body form.PositionUpdateParams true "Position"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /position/{id} [put]
func (co PositionController) Update(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.PositionUpdateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var position databases.MedHrmPosition
	result := co.DB.First(&position, c.Param("id"))
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	authUser := co.GetAuth(c)

	position.Name = params.Name
	position.Description = params.Description
	position.IsActive = params.IsActive
	position.PositionTypeID = uint(params.TypeID)
	position.Base.ModifiedDate = time.Now()
	position.ModifiedUser = &authUser

	result = co.DB.Save(&position)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	co.SetBody(structs.SuccessResponse{
		Success: true,
	})
	return
}

// Delete position
// @Summary Delete position
// @Description Remove position
// @Tags Position
// @Accept json
// @Produce json
// @Param position body form.DeleteParams true "position"
// @Success 200 {object} structs.ResponseBody{body=structs.SuccessResponse}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /position/{id} [delete]
func (co PositionController) Delete(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params form.DeleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	for _, v := range params.IDs {
		result := co.DB.Delete(&databases.MedHrmPosition{}, v)
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
