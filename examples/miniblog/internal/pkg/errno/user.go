package errno

var (
	// ErrUserAlreadyExist 代表用户已经存在.
	ErrUserAlreadyExist  = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: "User already exist."}
	ErrUserNotFound      = &Errno{HTTP: 400, Code: "FailedOperation.UserNotFound", Message: "User not found."}
	ErrPasswordIncorrect = &Errno{HTTP: 400, Code: "FailedOperation.PasswordIncorrect", Message: "Password incorrect,"}
)
