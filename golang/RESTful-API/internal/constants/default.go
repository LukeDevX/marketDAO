package constants

import "time"

const (
	DefaultZero        = 0
	DefaultEmptyString = ""

	DefaultLayout   = "2006-01-02 15:04:05"
	Default1970Time = "1970-01-01 23:59:59"

	DefaultErrTemplate = `%v: %v error, detail: %v`

	DefaultCacheExpired = time.Duration(0)

	DefaultRedisIndex = 0

	DefaultSuccessCode = 0
	DefaultSuccessMsg  = "success"

	DefaultMysqlUniqueErr        = "Error 1062"
	DefaultMysqlCardNumUniqueErr = "unique_idx_card_num"
	DefaultDBPort                = "3306"

	DefaultRequestTimeout = 10000000 * time.Second
	DefaultNonceRegExp    = "^[a-z0-9A-Z]{16}$"

	DefaultTagName = "json"

	DefaultErrRequestParamCode = 10000
	DefaultMinModelId          = 1

	DefaultOne = 1

	CodeListInsertLimit = 50000

	UploadGoodsThumbImageSuffix = "_thumb" //缩略图文件名后置

	MAX_UINT = ^uint(0)
	MIN_UINT = 0
	MAX_INT  = int(MAX_UINT >> 1)
	MIN_INT  = -MAX_INT - 1

	EntNoQuotes = 1

	DefaultCacheExpireTime = 86400 * time.Second

	DefaultPage     = 1
	DefaultPageSize = 10
)
