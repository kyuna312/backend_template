package middlewares

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
	databases "gitlab.com/fibocloud/medtech/gin/databases"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
	"gitlab.com/fibocloud/medtech/gin/utils"
	gorm "gorm.io/gorm"
)

// Response return auth middleware error
func Response(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, structs.ResponseBody{
		StatusCode: code,
		ErrorMsg:   message,
		Body:       nil,
	})
}

// Authenticate fetches user details from token
func Authenticate(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requiredToken := c.Request.Header["Authorization"]
		if len(requiredToken) == 0 {
			Response(c, http.StatusUnauthorized, "Please login to your account")
			return
		}

		claims, err := utils.ExtractJWTString(requiredToken[0][7:])
		if err != nil {
			Response(c, http.StatusUnauthorized, err.Error())
			return
		}

		var user databases.MedSystemUser
		result := db.Where(databases.MedSystemUser{Username: claims.Username}).Preload("Person").First(&user)
		if result.Error != nil {
			Response(c, http.StatusNotFound, result.Error.Error())
			return
		}

		c.Set("auth", user)
		c.Next()
	}
}

// Checkrole [ Орж ирсэн хэрэглэгчийн эрхийг шалгах ]
func Checkrole() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if iauth, exists := c.Get("auth"); exists {
		// 	// user := iauth.(databases.MedSystemUser)

		// }

		Response(c, http.StatusNotFound, "aldaa")
		return
	}
}
