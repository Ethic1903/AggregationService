package errors

import "errors"

//var (
//	ErrSubscriptionNotFound = NewAppError("subscription not found", http.StatusNotFound)
//	ErrInvalidInput         = NewAppError("invalid input", http.StatusBadRequest)
//	ErrInternalServer       = NewAppError("internal server error", http.StatusInternalServerError)
//	ErrDuplicateFound       = NewAppError("this subscription is already active", http.StatusConflict)
//)
//
//type AppError struct {
//	Message string `json:"message"`
//	Code    int    `json:"code"`
//}
//
//func NewAppError(message string, code int) *AppError {
//	return &AppError{
//		Message: message,
//		Code:    code,
//	}
//}
//
//func (e *AppError) Error() string {
//	return e.Message
//}
//
//func IsAppError(err error) (*AppError, bool) {
//	var appErr *AppError
//	if errors.As(err, &appErr) {
//		return appErr, true
//	}
//	return nil, false
//}

var (
	ErrSubscriptionNotFound     = errors.New("subscription not found")
	ErrInvalidRequest           = errors.New("invalid request")
	ErrSubscriptionAlreadyFound = errors.New("this subscription already active")
	ErrInvalidDateFormat        = errors.New("invalid date format")
	ErrInternalServer           = errors.New("internal server error occurred")
	ErrInvalidUUID              = errors.New("invalid UUID")
	ErrNoSubscriptionsFound     = errors.New("0 subscriptions were found")
	ErrInvalidPagination        = errors.New("invalid pagination parameters")
	ErrInvalidServiceName       = errors.New("invalid service name")
)
