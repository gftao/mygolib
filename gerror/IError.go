package gerror

type IError interface {
	Error() string
	GetErrorCode() string
	GetErrorString() string
}
