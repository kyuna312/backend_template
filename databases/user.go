package databases

import "time"

type (
	// MedBasePerson [ Ажилтан ]
	MedBasePerson struct {
		Base
		IsActive       bool             `gorm:"column:is_active;default:false" json:"is_active"`          // Идэвхтэй эсэх
		LastName       string           `gorm:"column:last_name;not null" json:"last_name"`               //
		FirstName      string           `gorm:"column:first_name;not null" json:"first_name"`             //
		MobileNumber   string           `gorm:"column:mobile_number;" json:"mobile_number"`               // Утас
		StateRegNumber string           `gorm:"column:state_reg_number;not null" json:"state_reg_number"` //
		SystemUsers    []*MedSystemUser `gorm:"foreignKey:PersonID" json:"system_users"`                  //
	}

	// MedSystemUser [ Хэрэглэгч ]
	MedSystemUser struct {
		Base
		IsActive     bool             `gorm:"column:is_active;default:false" json:"is_active"`    // Идэвхтэй эсэх
		Username     string           `gorm:"column:username;unique;not null" json:"username"`    // Нэр
		StartDate    time.Time        `gorm:"column:start_date;not null" json:"start_date"`       // Эхлэх огноо
		EndDate      time.Time        `gorm:"column:end_date;" json:"end_date"`                   // Дуусах огноо
		PasswordSalt string           `gorm:"column:password_salt;" json:"password_salt"`         //
		PasswordHash string           `gorm:"column:password_hash;not null" json:"password_hash"` //
		Person       *MedBasePerson   `gorm:"foreignKey:PersonID" json:"person"`
		PersonID     uint             `gorm:"column:person_id;" json:"person_id"`
		PersonType   uint             `gorm:"column:person_type;" json:"person_type"`
		Roles        []*MedSystemRole `gorm:"many2many:map_users_roles;"`
	}

	// MedHrmDepartment [ Ажилтны салбар, нэгж ]
	MedHrmDepartment struct {
		Base
		Code           string         `gorm:"column:code;not null" json:"code"`                // Барааны код
		Name           string         `gorm:"column:name;not null" json:"name"`                // Код
		Description    string         `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"` // Идэвхтэй эсэх
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` // Өөрчилсөн хэрэглэгч
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}

	// MedHrmPosition [ Албан тушаал ]
	MedHrmPosition struct {
		Base
		Code           string              `gorm:"column:code;not null" json:"code"`                // Барааны код
		Name           string              `gorm:"column:name;not null" json:"name"`                // Код
		PositionTypeID uint                `gorm:"column:position_type_id" json:"position_type_id"` // Өөрчилсөн хэрэглэгч
		PositionType   *MedHrmPositionType `gorm:"foreignKey:PositionTypeID" json:"position_type"`  // Үүсгэсэн хэрэглэгч
		Description    string              `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool                `gorm:"column:is_active;default:false" json:"is_active"` // Идэвхтэй эсэх
		CreatedUserID  uint                `gorm:"column:created_user_id" json:"created_user_id"`   // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint                `gorm:"column:modified_user_id" json:"modified_user_id"` // Өөрчилсөн хэрэглэгч
		CreatedUser    *MedSystemUser      `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser      `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}

	// MedHrmPositionType [ Албан тушаал төрөл ]
	MedHrmPositionType struct {
		Base
		Code           string         `gorm:"column:code;not null" json:"code"`                // Барааны код
		Name           string         `gorm:"column:name;not null" json:"name"`                // Код
		Description    string         `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"` // Идэвхтэй эсэх
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` // Өөрчилсөн хэрэглэгч
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}

	// MedHrmPositionKey [ Ажлын байр ]
	MedHrmPositionKey struct {
		Base
		Code           string            `gorm:"column:code;not null" json:"code"`                // Барааны код
		Name           string            `gorm:"column:name;not null" json:"name"`                // Код
		Description    string            `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool              `gorm:"column:is_active;default:false" json:"is_active"` // Идэвхтэй эсэх
		CreatedUserID  uint              `gorm:"column:created_user_id" json:"created_user_id"`   // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint              `gorm:"column:modified_user_id" json:"modified_user_id"` // Өөрчилсөн хэрэглэгч
		CreatedUser    *MedSystemUser    `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser    `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		Position       *MedHrmPosition   `gorm:"foreignKey:PositionID" json:"position"`           // positionid
		PositionID     uint              `gorm:"column:position_id" json:"position_id"`           // position_id
		Department     *MedHrmDepartment `gorm:"foreignKey:DepartmentID" json:"department"`       // positionid
		DepartmentID   uint              `gorm:"column:department_id" json:"department_id"`       // department_id
	}

	// MedSystemPermission [ Эрхийн тохиргоо ]
	MedSystemPermission struct {
		Base
		Code           string         `gorm:"column:code;not null" json:"code"`                // Барааны код
		Path           string         `gorm:"column:path;not null" json:"path"`                // API  зам
		Description    string         `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"` // Идэвхтэй эсэх
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` // Өөрчилсөн хэрэглэгч
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		// Roles          []*MedSystemRole `gorm:"many2many:roles;"`                                //
	}

	// MedSystemRole [ Эрхийн групп ]
	MedSystemRole struct {
		Base
		Code           string                 `gorm:"column:code;not null" json:"code"`                // Барааны код
		Name           string                 `gorm:"column:name;not null" json:"name"`                // Код
		Description    string                 `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool                   `gorm:"column:is_active;default:false" json:"is_active"` // Идэвхтэй эсэх
		CreatedUserID  uint                   `gorm:"column:created_user_id" json:"created_user_id"`   // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint                   `gorm:"column:modified_user_id" json:"modified_user_id"` // Өөрчилсөн хэрэглэгч
		CreatedUser    *MedSystemUser         `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser         `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		Permissions    []*MedSystemPermission `gorm:"many2many:map_permissions_roles;"`                //
		Users          []*MedSystemUser       `gorm:"many2many:map_users_roles;"`                      //
	}

	// MedCustomer ...
	MedCustomer struct {
		Base
		Code                  string                     `gorm:"column:code" json:"code"`                                        //
		Name                  string                     `gorm:"column:name;not null" json:"name"`                               //
		CompanyName           string                     `gorm:"column:company_name;not null" json:"company_name"`               //
		IsActive              bool                       `gorm:"column:is_active;default:false" json:"is_active"`                //
		Description           string                     `gorm:"column:description;" json:"description"`                         //
		ClassificationID      uint                       `gorm:"column:classification_id;" json:"classification_id"`             //
		Classification        *MedCustomerClassification `gorm:"foreignKey:ClassificationID" json:"classification"`              //
		CompanyRegistryNumber string                     `gorm:"column:company_registry_number;" json:"company_registry_number"` //
		CountryID             uint                       `gorm:"column:country_id;" json:"country_id"`                           //
		Country               *RefCountry                `gorm:"foreignKey:CountryID" json:"country"`                            //
		CityID                uint                       `gorm:"column:city_id;" json:"city_id"`                                 //
		City                  *RefCity                   `gorm:"foreignKey:CityID" json:"city"`                                  //
		DistrictID            uint                       `gorm:"column:district_id;" json:"district_id"`                         //
		District              *RefDistrict               `gorm:"foreignKey:DistrictID" json:"district"`                          //
		AddressDescription    string                     `gorn:"address_description" json:"address_description"`                 //
		Contacts              []*MedCustomerContacts     `gorm:"foreignKey:CustomerID" json:"contacts"`                          //
		CustomerAddress       []*MedCustomerAddress      `gorm:"foreignKey:CustomerID" json:"addresses"`                         //
		Types                 []*MedCustomerType         `gorm:"many2many:med_customer_type_dtl;" json:"types"`                  //
		Files                 []*MedContent              `gorm:"many2many:med_content_map" json:"files"`                         //
		MaximumPurchase       float64                    `gorm:"column:maximum_purchase" json:"maximum_purchase"`                //
		MaximumReceivables    float64                    `gorm:"column:maximum_receivables" json:"maximum_receivables"`          //
		OneTimePurchaseLimit  float64                    `gorm:"column:one_time_purchase_limit" json:"one_time_purchase_limit"`  //
		CreatedUserID         uint                       `gorm:"column:created_user_id" json:"created_user_id"`                  //
		ModifiedUserID        uint                       `gorm:"column:modified_user_id" json:"modified_user_id"`                //
		CreatedUser           *MedSystemUser             `gorm:"foreignKey:CreatedUserID" json:"created_user"`                   // Үүсгэсэн хэрэглэгч
		ModifiedUser          *MedSystemUser             `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`                 // Өөрчилсөн хэрэглэгч
		ParentID              uint                       `gorm:"column:parent_id" json:"parent_id"`                              //
		Parent                *MedCustomer               `gorm:"foreignKey:ParentID" json:"parent"`                              // Харьяалагдах
	}

	// MedCustomerAddress [  ]
	MedCustomerAddress struct {
		Base
		CustomerID     uint            `gorm:"column:customer_id" json:"customer_id"`           //
		Customer       *MedCustomer    `gorm:"foreignKey:CustomerID" json:"customer"`           //
		CountryID      uint            `gorm:"column:country_id;" json:"country_id"`            //
		Country        *RefCountry     `gorm:"foreignKey:CountryID" json:"country"`             //
		CityID         uint            `gorm:"column:city_id;" json:"city_id"`                  //
		City           *RefCity        `gorm:"foreignKey:CityID" json:"city"`                   //
		DistrictID     uint            `gorm:"column:district_id;" json:"district_id"`          //
		District       *RefDistrict    `gorm:"foreignKey:DistrictID" json:"district"`           //
		StreetID       uint            `gorm:"column:street_id;" json:"street_id"`              //
		Street         *RefStreet      `gorm:"foreignKey:StreetID" json:"street"`               //
		AddressTypeID  uint            `gorm:"column:address_type_id;" json:"address_type_id"`  //
		AddressType    *MedAddressType `gorm:"foreignKey:AddressTypeID" json:"address_type"`    //
		Description    string          `gorm:"column:description;" json:"description"`          // Тайлбар
		IsActive       bool            `gorm:"column:is_active;default:false" json:"is_active"` //
		CreatedUserID  uint            `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint            `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser  `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser  `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}

	// MedAddressType [ ]
	MedAddressType struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"` //
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}

	// MedCustomerType [ Нийлүүлэгч, Худалдан авагч ]
	MedCustomerType struct {
		Base
		ColorCode      string         `gorm:"column:color_code" json:"color_code"`             //
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"` //
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		Customers      []*MedCustomer `gorm:"many2many:med_customer_type_dtl" json:"customers"`
		// MedCustomers   []*MedCustomer `gorm:"foreignKey:type_id" json:"med_customers"`         //
	}

	// MedCustomerClassification [ ]
	MedCustomerClassification struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                 //
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"`  //
		Description    string         `gorm:"column:description;" json:"description"`           // Тайлбар
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`    //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"`  //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`     // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`   // Өөрчилсөн хэрэглэгч
		MedCustomers   []*MedCustomer `gorm:"foreignKey:ClassificationID" json:"med_customers"` //
	}

	// MedCustomerContacts [ ]
	MedCustomerContacts struct {
		Base
		CustomerID     uint            `gorm:"column:customer_id" json:"customer_id"`           //
		Customer       *MedCustomer    `gorm:"foreignKey:CustomerID" json:"customer"`           //
		LastName       string          `gorm:"column:last_name" json:"last_name"`               //
		FirstName      string          `gorm:"column:first_name" json:"first_name"`             //
		RegisterNumber string          `gorm:"column:register_number" json:"register_number"`   //
		PositionID     int             `gorm:"column:position_id" json:"position_id"`           //
		Position       *MedHrmPosition `gorm:"foreignKey:PositionID" json:"position"`           // Үүсгэсэн хэрэглэгч
		PhoneNumber1   string          `gorm:"column:phone_number1" json:"phone_number1"`       //
		PhoneNumber2   string          `gorm:"column:phone_number2" json:"phone_number2"`       //
		Email1         string          `gorm:"column:email1" json:"email1"`                     //
		Email2         string          `gorm:"column:email2" json:"email2"`                     //
		IsActive       bool            `gorm:"column:is_active;default:false" json:"is_active"` //
		CreatedUserID  uint            `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint            `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser  `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser  `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}

	// MedContent [ ]
	MedContent struct {
		Base
		FileName       string          `gorm:"column:file_name" json:"file_name"`               //
		Extention      string          `gorm:"column:extention" json:"extention"`               //
		PhysicalPath   string          `gorm:"column:physical_path" json:"physical_path"`       //
		FileSize       float64         `gorm:"column:file_size" json:"file_size"`               //
		ContentTypeID  uint            `gorm:"column:content_type_id" json:"content_type_id"`   //
		ContentType    *MedContentType `gorm:"foreignKey:ContentTypeID" json:"content_type"`    // Үүсгэсэн хэрэглэгч
		ModifiedUserID uint            `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUserID  uint            `gorm:"column:created_user_id" json:"created_user_id"`   //
		CreatedUser    *MedSystemUser  `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser  `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
		Customer       []*MedCustomer  `gorm:"many2many:med_content_map;"`
		// ContentType    *MedContentType `gorm:"foreignKey:ContentTypeID" json:"content_type"`          //
	}

	// MedContentType [ ]
	MedContentType struct {
		Base
		Name           string         `gorm:"column:name;not null" json:"name"`                //
		IsActive       bool           `gorm:"column:is_active;default:false" json:"is_active"` //
		ParentID       uint           `gorm:"column:parent_id" json:"parent_id"`
		CreatedUserID  uint           `gorm:"column:created_user_id" json:"created_user_id"`   //
		ModifiedUserID uint           `gorm:"column:modified_user_id" json:"modified_user_id"` //
		CreatedUser    *MedSystemUser `gorm:"foreignKey:CreatedUserID" json:"created_user"`    // Үүсгэсэн хэрэглэгч
		ModifiedUser   *MedSystemUser `gorm:"foreignKey:ModifiedUserID" json:"modified_user"`  // Өөрчилсөн хэрэглэгч
	}
)
