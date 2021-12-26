package databases

import (
	"fmt"
	"time"

	viper "github.com/spf13/viper"
	postgres "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

type (
	// Base struct
	Base struct {
		ID           uint      `gorm:"primary_key" json:"id"`                              //
		ModifiedDate time.Time `gorm:"column:modified_date;not null" json:"modified_date"` // Өөрчилсөн огноо
		CreatedDate  time.Time `gorm:"column:created_date;not null" json:"created_date"`   // Үүсгэсэн огноо
	}

	// RefCountry struct
	RefCountry struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		Cities         []*RefCity     `gorm:"foreignKey:CountryID" json:"cities"`              //
	}

	// RefCity struct
	RefCity struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		Country        *RefCountry    `gorm:"foreignKey:CountryID" json:"country"`             // Улс
		CountryID      uint           `gorm:"column:country_id" json:"country_id"`             //
	}

	// RefDistrict struct
	RefDistrict struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		City           *RefCity       `gorm:"foreignKey:CityID" json:"city"`                   // Хот
		CityID         uint           `gorm:"column:city_id" json:"city_id"`                   //
	}

	// RefStreet struct
	RefStreet struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		District       *RefDistrict   `gorm:"foreignKey:DistrictID" json:"district"`           // Дүүрэг
		DistrictID     uint           `gorm:"column:district_id" json:"district_id"`           //
	}
)

// InitDB initialize databases and tables
func InitDB() *gorm.DB {
	DBUser := viper.GetString("database.user")
	DBPassword := viper.GetString("database.password")
	DBDatabase := viper.GetString("database.name")
	DBHost := viper.GetString("database.host")
	DBPort := viper.GetString("database.port")
	DBTimezone := viper.GetString("database.timezone")
	ConnectionString := fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable TimeZone=%v",
		DBHost,
		DBPort,
		DBUser,
		DBDatabase,
		DBPassword,
		DBTimezone,
	)

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold: time.Second,   // Slow SQL threshold
	// 		LogLevel:      logger.Silent, // Log level
	// 		Colorful:      false,         // Disable color
	// 	},
	// )

	db, err := gorm.Open(postgres.Open(ConnectionString), &gorm.Config{
		// Logger:                                   newLogger,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(
		&MedBasePerson{},
		&MedSystemUser{},
		&MedHrmDepartment{},
		&MedHrmPosition{},
		&MedHrmPositionType{},
		&MedHrmPositionKey{},
		&MedSystemPermission{},
		&MedSystemRole{},
		&MedCustomer{},
		&MedCustomerType{},
		&MedCustomerAddress{},
		&MedAddressType{},
		&MedCustomerClassification{},
		&MedCustomerContacts{},
		&MedContent{},
		&MedContentType{},
		&RefCountry{},
		&RefCity{},
	)
	return db
}

// func (u *MedSystemUser) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.UUID = uuid.New()

// 	if !u.IsValid() {
// 		err = errors.New("can't save invalid data")
// 	}
// 	return
// }
