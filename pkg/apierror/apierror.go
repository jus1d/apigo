package apierror

type Type string

const (
	TypeInvalidRequest Type = "invalid_request_error"
	TypeAuthentication Type = "authentication_error"
	TypePermission     Type = "permission_error"
	TypeNotFound       Type = "not_found_error"
	TypeConflict       Type = "conflict_error"
	TypeRateLimit      Type = "rate_limit_error"
	TypeInternal       Type = "internal_error"
)

type Error struct {
	Status  int
	Type    Type
	Message string
	Hint    string
}

func New(status int, errType Type, message string, hint string) error {
	return &Error{
		Status:  status,
		Type:    errType,
		Message: message,
		Hint:    hint,
	}
}

func (e *Error) Error() string {
	return e.Message
}
