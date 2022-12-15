package handler

var MISSING_FIELD = newError("missing required field(s)")
var NOT_FOUND = newError("resource not found")
var UNAUTHORIZED = newError("you don't have the permission to perform such task")
var INTERNAL_ERROR = newError("something went wrong, please try again later")
var INVALID_TOKEN = newError("unable to cast token, maybe expired?")

type response struct {
	Err  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// swagger:model
type SuccessResponse struct {
	// the data
	// the data is in interface, it means it can be anything
	// example: {"id": 1, "name": "john doe", ...}
	Data interface{} `json:"data,omitempty"`
}

// swagger:model
type BadRequestResponse struct {
	// the error message
	// in: body
	// example: missing required field(s)
	Err string `json:"message"`
}

// swagger:model
type InternalErrorResponse struct {
	// the error message
	// in: body
	// example: something went wrong, please try again later
	Err string `json:"message"`
}

// swagger:model
type UnauthorizedResponse struct {
	// the error message
	// in: body
	// example: you don't have the permission to perform such task
	Err string `json:"message"`
}

// swagger:model
type NotFoundResponse struct {
	// the error message
	// in: body
	// example: resource not found
	Err string `json:"message"`
}

func newError(err string) *response {
	return &response{err, nil}
}

func newResponse(data interface{}) *response {
	return &response{"", data}
}
