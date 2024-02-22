package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400
	// ErrorExistUser user模块错误 3xxxx
	ErrorExistUser             = 30001
	ErrorFailEncryption        = 30002
	ErrorExistUserNotFound     = 30003
	ErrorNotCompare            = 30004
	ErrorAuthToken             = 30005
	ErrorAuthCheckTokenTimeout = 30006
	ErrorUploadFile            = 30007
	ErrorSendEmail             = 30008
	ErrorSetKey                = 30009
	//product 模块错误
	ErrorProductImgUpload = 40001
	ErrorDatabase         = 40002

	//成员错误
	ErrorExistFavorite    = 10001
	ErrorProductExistCart = 10002
	ErrorProductMoreCart  = 10003
)
