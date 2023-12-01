package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ../../../docs/guide/zh-CN/api/error_code_generated.md

// Common: basic errors
// Code must start with 1xxxxx
const (
	// Success - 200: OK.
	Success int = iota + 100001

	// ErrUnknown : 内部服务错误
	ErrUnknown

	// ErrBind : 请求参数绑定错误
	ErrBind

	// ErrValidation : 参数验证错误
	ErrValidation

	// ErrTokenInvalid : 无效的token
	ErrTokenInvalid
)

// common: database errors.
const (
	// ErrDatabase : Database error.
	ErrDatabase int = iota + 100101
)

// common: authorization and authentication errors.
const (
	// ErrEncrypt : Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201

	// ErrSignatureInvalid : Signature is invalid.
	ErrSignatureInvalid

	// ErrExpired : Token expired.
	ErrExpired

	// ErrInvalidAuthHeader : Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader : The `Authorization` header was empty.
	ErrMissingHeader

	// ErrPasswordIncorrect : Password was incorrect.
	ErrPasswordIncorrect

	// ErrPermissionDenied : Permission denied.
	ErrPermissionDenied
)

// common: encode/decode errors.
const (
	// ErrEncodingFailed : Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// ErrDecodingFailed : Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON : Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON : JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON : JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYaml : Data is not valid Yaml.
	ErrInvalidYaml

	// ErrEncodingYaml : Yaml data could not be encoded.
	ErrEncodingYaml

	// ErrDecodingYaml : Yaml data could not be decoded.
	ErrDecodingYaml
)
