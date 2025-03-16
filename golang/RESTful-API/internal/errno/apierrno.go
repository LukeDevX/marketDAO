package errno

import "fmt"

// ServerError server error interface
type ServerError interface {
	String() string
}

// ServerErrno server error struct
type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

const (
	TgBotAPICode = 80000
)

var (
	ErrServerBusy           = &ErrMsg{Code: 503, Msg: "Server is busy now!!!"}
	ErrNeedLogin            = &ErrMsg{Code: 400, Msg: "请登录后再操作"}
	ErrNotFound             = &ErrMsg{Code: 404, Msg: "页面不存在"}
	ErrServer               = &ErrMsg{Code: 500, Msg: "服务器错误"}
	ErrServerNotImplemented = &ErrMsg{Code: 501, Msg: "服务器错误"}
	ErrServerLoginAuthError = &ErrMsg{Code: 502, Msg: "登录校验失败，请重新登录"}
	ErrRecordNotFoundError  = &ErrMsg{Code: 505, Msg: "记录不存在"}

	ErrRequestTimeout = &ErrMsg{Code: 408, Msg: "请求超时"}

	ErrQuery                    = &ErrMsg{Code: 409, Msg: "服务器异常"}
	ErrUpdate                   = &ErrMsg{Code: 410, Msg: "服务器异常"}
	ErrVerify                   = &ErrMsg{Code: 411, Msg: "服务器异常"}
	ErrParam                    = &ErrMsg{Code: 412, Msg: "params error"}
	ErrParamLost                = &ErrMsg{Code: 413, Msg: "参数丢失"}
	ErrHandleInvalid            = &ErrMsg{Code: 414, Msg: "操作非法"}
	ErrUserFromNotExists        = &ErrMsg{Code: 418, Msg: "请求来源不存在"}
	ErrSourceTagNotExists       = &ErrMsg{Code: 419, Msg: "请求来源不存在"}
	ErrUserHandleTooOftenExists = &ErrMsg{Code: 420, Msg: "Operation is too frequent, please try later"}

	ErrGetUserCredEmpty = &ErrMsg{Code: 10001, Msg: "获取登录凭证为空"}

	ErrAddressSubmitRepeatError         = &ErrMsg{Code: 20001, Msg: "This address has already been submitted. Do not submit it again"}
	ErrTwitterUsernameSubmitRepeatError = &ErrMsg{Code: 20002, Msg: "Twitter Username has already been submitted. Do not submit it again"}
	ErrMustAddWorkError                 = &ErrMsg{Code: 20003, Msg: "You must add a work"}
	ErrAddExceedMaxNumberOnceError      = &ErrMsg{Code: 20004, Msg: "Add up to 10 works at a time"}
	ErrWorkTypeError                    = &ErrMsg{Code: 20005, Msg: "Work type error"}
	ErrWorkReleaseTimeError             = &ErrMsg{Code: 20006, Msg: "Work release time error"}
	ErrCreatorNameMaxError              = &ErrMsg{Code: 20007, Msg: "The name of the creator cannot exceed 3"}
	ErrDuplicateNameOfWorkError         = &ErrMsg{Code: 20008, Msg: "Duplicate name of work"}
	ErrInvalidEmailAddressError         = &ErrMsg{Code: 20010, Msg: "Invalid email address"}
	ErrOperationTypeError               = &ErrMsg{Code: 20011, Msg: "Operation type error"}
	ErrYouMustCancelTheLikeBeforeError  = &ErrMsg{Code: 20012, Msg: "You must cancel the like before you can do anything else"}
	ErrYouMustCancelTheTapBeforeError   = &ErrMsg{Code: 20013, Msg: "You must cancel the tap before you can do anything else"}
	ErrWorkNotExistsError               = &ErrMsg{Code: 20014, Msg: "Work does not exist"}
	ErrAirdropSubmitTypeError           = &ErrMsg{Code: 20015, Msg: "You've already submitted another type of airdrop"}
	ErrUploadMustMintError              = &ErrMsg{Code: 20016, Msg: "You must mint before you can upload your work"}
	ErrChainIDInvalidError              = &ErrMsg{Code: 20017, Msg: "The chain ID is invalid"}
	ErrHashNotExists                    = &ErrMsg{Code: 20018, Msg: "tx hash not exists"}

	ErrUploadErrorError = &ErrMsg{Code: 80000, Msg: "Upload Error"}
	ErrFileTypeError    = &ErrMsg{Code: 80001, Msg: "File type error"}

	ErrTokenVerificationError = &ErrMsg{Code: 30001, Msg: "token verification fails or has expired. Please log in again"}
	ErrNotLoginError          = &ErrMsg{Code: 30002, Msg: "You are not logged in yet, please log in first"}
	ErrTonLoginAuthError      = &ErrMsg{Code: 30003, Msg: "Login authentication failure"}
	ErrSignVerifyError        = &ErrMsg{Code: 30004, Msg: "Signature verification failed"}

	ErrRepeatSubmitChannelInfoError = &ErrMsg{Code: 40001, Msg: "You have already submitted the information of the channel provider, please do not submit it again"}
)

// String encode struct to string
func (e *ErrMsg) String() string {
	return fmt.Sprintf("%s(%d)", e.Msg, e.Code)
}
