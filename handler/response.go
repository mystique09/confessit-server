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

func newError(err string) *response {
	return &response{err, nil}
}

func newResponse(data interface{}) *response {
	return &response{"", data}
}
