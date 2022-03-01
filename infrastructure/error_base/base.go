package error_base

type Error struct {
	Error error
	Code  uint32
	Line  string
}

func New(code uint32, err error) *Error {
	return &Error{
		Error: err,
		Code:  code,
		Line:  MapError(code),
	}
}

func MapError(errorCode uint32) string {
	switch errorCode {
	case ErrorBindDataCode:
		return ErrorBindDataCodeMsg
	case ErrorValidDataCode:
		return ErrorValidDataMsg
	case ErrorSendDataCode:
		return ErrorSendDataMsg
	}

	return "Unknown error"
}
