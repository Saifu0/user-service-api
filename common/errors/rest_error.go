package errors

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewBadRequest(message string) RestErr {
	return RestErr{
		Message: message,
		Status:  400,
		Error:   "bad request",
	}
}
