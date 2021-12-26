package shared

import (
	"log"
	"net/http"
	"reflect"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/fibocloud/medtech/gin/databases"
	"gitlab.com/fibocloud/medtech/gin/form"
	structs "gitlab.com/fibocloud/medtech/gin/structs"
	gorm "gorm.io/gorm"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// BaseController struct
type BaseController struct {
	Response    *structs.Response
	DB          *gorm.DB
	minioClient *minio.Client
}

// MinioClinet ..
func (co BaseController) MinioClinet() *minio.Client {
	minioClient, err := minio.New(viper.GetString("minio.endpoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(viper.GetString("minio.accessKeyID"), viper.GetString("minio.secretAccessKey"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln("err minio", err)
	}

	return minioClient
}

// ListResponse ...
type ListResponse struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

// SetBody successfully response
func (co BaseController) SetBody(body interface{}) {
	co.Response.StatusCode = http.StatusOK
	co.Response.Body.StatusCode = 0
	co.Response.Body.ErrorMsg = ""
	co.Response.Body.Body = body
}

// SetError rrror response
func (co BaseController) SetError(code int, message string) {
	co.Response.StatusCode = 200
	co.Response.Body.StatusCode = code
	co.Response.Body.ErrorMsg = message
	co.Response.Body.Body = nil
}

// GetBody in response
func (co BaseController) GetBody() (int, interface{}) {
	return co.Response.StatusCode, co.Response.Body
}

// GetAuth get auth user
func (co BaseController) GetAuth(c *gin.Context) databases.MedSystemUser {
	if iauth, exists := c.Get("auth"); exists {
		return iauth.(databases.MedSystemUser)
	}
	return databases.MedSystemUser{}
}

// Paginate table
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// TableSearch undsen table search hiih
func TableSearch(v reflect.Value, sort form.SortColumn) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		typeOfS := v.Type()
		for i := 0; i < v.NumField(); i++ {

			filterName := typeOfS.Field(i).Name
			filterJSON := typeOfS.Field(i).Tag.Get("json")

			if !strings.Contains(filterName, "External") {
				switch v.Field(i).Interface().(type) {
				case int:
					fitlerValue := v.Field(i).Interface().(int)
					if fitlerValue > 0 {
						db = db.Where(filterJSON+" = ?", fitlerValue)
					}

				case string:
					fitlerValue := v.Field(i).Interface().(string)

					isDate := strings.Contains(strings.ToLower(filterName), "date")

					if fitlerValue != "" && fitlerValue != "0" {
						if filterName == "IsImport" || filterName == "IsActive" {
							boolValue := fitlerValue == "true"
							db = db.Where(filterJSON+" = ?", boolValue)
						}
						if isDate {
							db = db.Where(filterJSON+" >= ?", fitlerValue)
						}

						if !isDate && filterName != "IsImport" && filterName != "IsActive" {
							db = db.Where("LOWER(CAST("+filterJSON+" as TEXT))"+" LIKE ?", "%"+strings.ToLower(fitlerValue)+"%")
						}
					}
				}
			}
		}

		if sort.Field != "" {
			db = db.Order(sort.Field + " " + sort.Order)
		} else {
			db = db.Order("created_date desc")
		}
		return db
	}
}
