package apierror

type Error struct {
	Status  int
	Message string
}

func New(code int, message string) error {
	return &Error{
		Status:  code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
