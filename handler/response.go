package handler

type response struct {
	Err  string      `json:"error,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func newError(err string) *response {
	return &response{err, nil}
}

func newResponse(data interface{}) *response {
	return &response{"", data}
}
