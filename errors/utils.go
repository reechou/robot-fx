package errors

func GetErrorMessage(err error) string {
	switch err.(type) {
	case Error:
		e, _ := err.(Error)
		return e.Message

	case ErrorCode:
		ec, _ := err.(ErrorCode)
		return ec.Message()

	default:
		return err.Error()
	}
}
