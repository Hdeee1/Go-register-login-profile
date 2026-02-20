package response

type ErrorDetail struct {
	Code	string	`json:"status_code"`
	Message	string	`json:"message"`
}

type APIResponse struct {
	Success	bool			`json:"success"`
	Message	string			`json:"message,omitempty"`
	Data	any				`json:"data,omitempty"`
	Error	*ErrorDetail	`json:"error,omitempty"`
}

func BuildSuccessResponse(message string, data any) APIResponse {
	response := APIResponse{
		Success: true,
		Message: message,
		Data: data,
		Error: nil,
	}
	
	return response
}

func BuildErrorResponse(code, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &ErrorDetail{
			Code: code,
			Message: message,
		},
	}
}