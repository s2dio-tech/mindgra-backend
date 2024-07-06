package common

import "errors"

var (
	ErrInternalServerError = errors.New("internalServerError")
	ErrNotFound            = errors.New("notFound")
	ErrConflict            = errors.New("alreadyExist")
	ErrBadParamInput       = errors.New("badParamInput")
	ErrEmailDuplicate      = errors.New("emailDuplicated")
	ErrTokenInvalid        = errors.New("tokenInvalid")
	ErrUnauthentication    = errors.New("unAuthentication")
	ErrUnauthorization     = errors.New("unAuthorization")
	ErrInvalidCredential   = errors.New("invalidCredential")
)
