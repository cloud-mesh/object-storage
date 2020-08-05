package model

const (
	ErrCodeOK      = 0
	ErrCodeUnknown = 99999

	errorCodeBase = 15000
)

var (
	ErrInvalidParam = newError(errorCodeBase+0, "invalid parameters")

	ErrVendorGetBucket                 = newError(errorCodeBase+500, "vendor get bucket")
	ErrVendorListBucket                = newError(errorCodeBase+501, "vendor get bucket")
	ErrVendorMakeBucket                = newError(errorCodeBase+502, "vendor make bucket")
	ErrVendorRemoveBucket              = newError(errorCodeBase+503, "vendor remove bucket")
	ErrVendorHeadObject                = newError(errorCodeBase+504, "vendor head object")
	ErrVendorPutObject                 = newError(errorCodeBase+505, "vendor put object")
	ErrVendorRemoveObject              = newError(errorCodeBase+506, "vendor remove object")
	ErrVendorMultipartUploadInit       = newError(errorCodeBase+507, "vendor init multipart upload")
	ErrVendorMultipartUploadComplete   = newError(errorCodeBase+508, "vendor complete multipart upload")
	ErrVendorMultipartUploadAbort      = newError(errorCodeBase+509, "vendor abort multipart upload")
	ErrVendorMultipartUploadUploadPart = newError(errorCodeBase+510, "vendor multipart upload upload part")
	ErrVendorMultipartUploadListPart   = newError(errorCodeBase+511, "vendor multipart upload list part")
)

type Error struct {
	errCode int
	text    string
}

func newError(errCode int, text string) Error {
	return Error{
		errCode: errCode,
		text:    text,
	}
}

func (e Error) Code() int {
	return e.errCode
}

func (e Error) Error() string {
	return e.text
}

func GetCode(err error) int {
	if err == nil {
		return ErrCodeOK
	}

	errCode := ErrCodeUnknown // 未知错误
	if err, ok := err.(interface {
		Code() int
	}); ok {
		errCode = err.Code()
	}

	return errCode
}
