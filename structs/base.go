package structs

type (
	// Response struct
	Response struct {
		StatusCode int
		Body       ResponseBody
	}

	// ResponseBody struct
	ResponseBody struct {
		StatusCode int         `json:"status_code"`
		ErrorMsg   string      `json:"error_msg"`
		Body       interface{} `json:"body"`
	}

	// SuccessResponse struct
	SuccessResponse struct {
		Success bool `json:"success"`
	}

	// ErrorResponse struct
	ErrorResponse struct {
		StatusCode int    `json:"status_code"`
		ErrorMsg   string `json:"error_msg"`
		Body       string `json:"body"`
	}
)
