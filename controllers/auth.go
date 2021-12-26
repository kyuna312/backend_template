package controllers

import (
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
	"gitlab.com/fibocloud/medtech/gin/controllers/shared"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	utils "gitlab.com/fibocloud/medtech/gin/utils"
)

// AuthController struct
type AuthController struct {
	shared.BaseController
}

// Init Controller
func (co AuthController) Init(router *gin.RouterGroup) {
	router.GET("/admin", co.Admin)  // Admin create
	router.POST("/login", co.Login) // Login
}

// LoginParams create body params
type LoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResult create body params
type LoginResult struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

// Login user
// @Summary Sign in user
// @Description Sign in user
// @Tags Auth
// @Accept json
// @Produce json
// @Param auth body LoginParams true "Auth"
// @Success 200 {object} structs.ResponseBody{body=LoginResult}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /auth/login [post]
func (co AuthController) Login(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	var params LoginParams
	if err := c.ShouldBindJSON(&params); err != nil {
		co.SetError(http.StatusBadRequest, err.Error())
		return
	}

	var user databases.MedSystemUser
	result := co.DB.Where("person_type = ?", 1).Where("username = ? AND is_active = true AND end_date >= ?", params.Username, time.Now()).Preload("Person").First(&user)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			co.SetError(http.StatusNotFound, "Хэрэглэгч олдсонгүй")
		}
		return
	}

	if valid, err := utils.ComparePassword(user.PasswordHash, params.Password); !valid {
		if err.Error() == "record not found" {
			co.SetError(http.StatusNotFound, "Нэвтрэх нэр эсвэл нууц үг буруу байна")
		}
		// co.SetError(http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, refreshToken := utils.GenerateToken(user)
	co.SetBody(LoginResult{Token: accessToken, Refresh: refreshToken})
	return
}

// Admin create
// @Summary Init admin account
// @Description Init admin account
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} structs.ResponseBody{body=databases.MedSystemUser}
// @Failure 400 {object} structs.ErrorResponse
// @Failure 500 {object} structs.ErrorResponse
// @Router /auth/admin [get]
func (co AuthController) Admin(c *gin.Context) {
	defer func() {
		c.JSON(co.GetBody())
	}()

	hashPwd, err := utils.GenerateHash("Mongol123@")
	if err != nil {
		co.SetError(http.StatusInternalServerError, err.Error())
		return
	}

	baseUser := databases.MedBasePerson{
		LastName:       "Erdenebat",
		FirstName:      "Darkhanbayar",
		StateRegNumber: "TA97072630",
		IsActive:       true,
		MobileNumber:   "88878150",
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	result := co.DB.Create(&baseUser)
	if result.Error != nil {
		co.SetError(http.StatusInternalServerError, result.Error.Error())
		return
	}

	var baseAdminUser databases.MedBasePerson
	resultBaseuser := co.DB.Where("state_reg_number = ?", "TA97072630").First(&baseAdminUser)
	if resultBaseuser.Error != nil {
		co.SetError(http.StatusInternalServerError, resultBaseuser.Error.Error())
		return
	}

	user := databases.MedSystemUser{
		PersonID:     baseAdminUser.Base.ID,
		IsActive:     true,
		Username:     "admin",
		StartDate:    time.Now(),
		EndDate:      time.Now().Add(8640 * time.Hour),
		PasswordSalt: "",
		PasswordHash: hashPwd,
		Base: databases.Base{
			CreatedDate: time.Now(),
		},
	}

	resultSystemUser := co.DB.Create(&user)
	if resultSystemUser.Error != nil {
		co.SetError(http.StatusInternalServerError, resultSystemUser.Error.Error())
		return
	}

	co.SetBody(user)

	return
}
