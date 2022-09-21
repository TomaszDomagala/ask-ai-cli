package errs

type CmdError struct {
	// Err is the error that occurred.
	Err error
	// Msg is a friendly message to the user.
	Msg string
}

func (e *CmdError) Error() string {
	return e.Err.Error()
}

func (e *CmdError) Unwrap() error {
	return e.Err
}

func New(err error, msg string) *CmdError {
	return &CmdError{
		Err: err,
		Msg: msg,
	}
}
