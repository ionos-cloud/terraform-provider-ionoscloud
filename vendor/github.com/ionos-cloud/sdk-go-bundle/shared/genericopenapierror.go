package shared

// GenericOpenAPIError Provides access to the body, error and model on returned errors.
type GenericOpenAPIError struct {
	statusCode int
	body       []byte
	error      string
	model      interface{}
}

// NewGenericOpenAPIError - constructor for GenericOpenAPIError
func NewGenericOpenAPIError(message string, body []byte, model interface{}, statusCode int) *GenericOpenAPIError {
	return &GenericOpenAPIError{
		statusCode: statusCode,
		body:       body,
		error:      message,
		model:      model,
	}
}

// Error returns non-empty string if there was an error.
func (e GenericOpenAPIError) Error() string {
	return e.error
}

// SetError sets the error string
func (e *GenericOpenAPIError) SetError(error string) {
	e.error = error
}

// Body returns the raw bytes of the response
func (e GenericOpenAPIError) Body() []byte {
	return e.body
}

// SetBody sets the raw body of the error
func (e *GenericOpenAPIError) SetBody(body []byte) {
	e.body = body
}

// Model returns the unpacked model of the error
func (e GenericOpenAPIError) Model() interface{} {
	return e.model
}

// SetModel sets the model of the error
func (e *GenericOpenAPIError) SetModel(model interface{}) {
	e.model = model
}

// StatusCode returns the status code of the error
func (e GenericOpenAPIError) StatusCode() int {
	return e.statusCode
}

// SetStatusCode sets the status code of the error
func (e *GenericOpenAPIError) SetStatusCode(statusCode int) {
	e.statusCode = statusCode
}
