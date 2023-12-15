package code

// admin-srv: sysUser errors.
const (
	// ErrUserNotFound : User not found.
	ErrUserNotFound int = iota + 110101

	// ErrUserAlreadyExist : User already exist.
	ErrUserAlreadyExist

	// ErrFailedAuthentication : username or password error
	ErrFailedAuthentication
)

// admin-srv: baseMenu errors
const (

	// ErrCreateBaseMenu : 创建菜单失败
	ErrCreateBaseMenu int = iota + 110201
)
