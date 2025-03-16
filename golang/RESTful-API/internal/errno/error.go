package errno

import (
	"errors"
)

var (
	NonPointerOrEmptyError     = errors.New("value can not be empty or non-pointer")
	InvalidParamError          = errors.New("invalid param")
	ErrInternal                = errors.New("internal error")
	ErrRequestParamOriginal    = errors.New("参数错误")
	ErrCacheVal                = errors.New("缓存数据不能为空")
	ErrCacheAcceptor           = errors.New("缓存数据接收对象不能为空")
	ErrMergeFieldLessThanValue = errors.New("fields个数大于values")
	ErrRedisNoValue            = errors.New("redis no value")
	ErrEmptyResponse           = errors.New("empty response")
	ErrNotInArray              = errors.New("not in array")

	ErrIntegerVal = errors.New("invalid integer")
	ErrFloatVal   = errors.New("invalid float val")

	ErrOrmUpdateDeny = errors.New("update deny")
)
