package config

type HttpError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func UnauthorizedRoleError() HttpError {
	return HttpError{
		401,
		"Unauthorized",
		"Your user role is not authorized to access this resource",
	}
}
func InvalidTokenError() HttpError {
	return HttpError{
		401,
		"Unauthorized",
		"This token is not valid",
	}
}
func UnauthorizedError() HttpError {
	return HttpError{
		401,
		"Unauthorized",
		"You are not authorized to access this resource",
	}
}
func NotFoundError() *HttpError {
	return &HttpError{
		404,
		"Not found",
		"The requested resource was not found",
	}
}
func DataAccessLayerError(message string) *HttpError {
	return &HttpError{
		500,
		"Data access error",
		message,
	}
}
func BadRequestError(message string) *HttpError {
	return &HttpError{
		400,
		"Bad Request",
		message,
	}
}
func BadDepositError(message string) *HttpError {
	return &HttpError{
		400,
		"Only 5, 10, 20, 50 and 100 cent coins are accepted",
		message,
	}
}
func BadProductCostError(message string) *HttpError {
	return &HttpError{
		400,
		"Check product cost (must be an 5 multiple)",
		message,
	}
}
func BusinessLayerError(message string) *HttpError {
	return &HttpError{
		500,
		"Business Rule Constraint error",
		message,
	}
}
