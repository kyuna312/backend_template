package controllers

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	reference "gitlab.com/fibocloud/medtech/gin/controllers/reference"
	shared "gitlab.com/fibocloud/medtech/gin/controllers/shared"
	User "gitlab.com/fibocloud/medtech/gin/controllers/user"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	middlewares "gitlab.com/fibocloud/medtech/gin/middlewares"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
)

// Init Controller
func Init(router *gin.RouterGroup) {
	db := databases.InitDB()
	bc := shared.BaseController{
		Response: &structs.Response{
			StatusCode: http.StatusOK,
			Body: structs.ResponseBody{
				StatusCode: 0,
				ErrorMsg:   "",
				Body:       nil,
			},
		},
		DB: db,
	}
	AuthController{bc}.Init(router.Group("/auth"))
	authRouter := router.Group("")
	authRouter.Use(middlewares.Authenticate(db))
	// authRouter.Use(middlewares.Checkrole())

	{

		CsvController{bc}.Init(authRouter.Group("/csv"))

		// region [ User ]
		User.CustomerController{bc}.Init(authRouter.Group("/customer"))
		User.CustomerTypeController{bc}.Init(authRouter.Group("/customerType"))
		User.CustomerClassificationController{bc}.Init(authRouter.Group("/customerClassification"))
		User.DepartmentController{bc}.Init(authRouter.Group("/department"))
		User.PersonController{bc}.Init(authRouter.Group("/person"))
		User.PermissionController{bc}.Init(authRouter.Group("/permission"))
		User.RoleController{bc}.Init(authRouter.Group("/role"))
		User.PositionController{bc}.Init(authRouter.Group("/position"))
		User.PositionTypeController{bc}.Init(authRouter.Group("/positionType"))
		User.PositionKeyController{bc}.Init(authRouter.Group("/positionKey"))
		// endregion

		// region [ Reference ]
		reference.PaymentTypeController{bc}.Init(authRouter.Group("/paymentType"))
		reference.PaymentMethodController{bc}.Init(authRouter.Group("/paymentMethod"))
		reference.CountryController{bc}.Init(authRouter.Group("/country"))
		reference.CityController{bc}.Init(authRouter.Group("/city"))
		reference.DistrictController{bc}.Init(authRouter.Group("/district"))
		reference.StreetController{bc}.Init(authRouter.Group("/street"))
		reference.StatusController{bc}.Init(authRouter.Group("/status"))
		reference.StatusTypeController{bc}.Init(authRouter.Group("/statusType"))
		reference.ValuteController{bc}.Init(authRouter.Group("/valute"))
		reference.AddressTypeController{bc}.Init(authRouter.Group("/addressType"))
		// endregion
	}
}
