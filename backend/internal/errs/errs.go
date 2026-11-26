package errs

import "errors"

type Err string

func (e Err) Error() string {
	return string(e)
}

// common errors
const (
	ServiceNA         = Err("service_not_available")
	NotImplemented    = Err("not_implemented")
	InvalidConfig     = Err("invalid_config")
	NoPermission      = Err("no_permission")
	ObjectNotFound    = Err("object_not_found")
	NoRows            = Err("err_no_rows")
	NotAuthorized     = Err("not_authorized")
	InvalidRequest    = Err("invalid_request")
	IncorrectPageSize = Err("incorrect_page_size")
	IdRequired        = Err("id_required")
)

// service errors
const (
	DomainNotAllowed   = Err("domain_not_allowed")
	EmailNotVerified   = Err("email_not_verified")
	AccountBlocked     = Err("account_blocked")
	FileTooLarge       = Err("file_too_large")
	FileTypeNotAllowed = Err("file_type_not_allowed")
)

type ErrFull struct {
	Err    error
	Desc   string
	Fields map[string]string
}

func (e ErrFull) Error() string {
	return e.Err.Error() + ", desc: " + e.Desc
}

// Field is a validation error for a single field (invalid_request + fields)
func Field(name, desc string) ErrFull {
	return ErrFull{Err: InvalidRequest, Desc: desc, Fields: map[string]string{name: desc}}
}

// AsErr extracts an Err from the error chain regardless of whether it was
// returned as a value (Err) or a pointer (*Err).
func AsErr(err error) (Err, bool) {
	if e, ok := errors.AsType[Err](err); ok {
		return e, true
	}
	if e, ok := errors.AsType[*Err](err); ok && e != nil {
		return *e, true
	}
	return "", false
}

// AsErrFull extracts an ErrFull from the error chain regardless of whether it
// was returned as a value (ErrFull) or a pointer (*ErrFull).
func AsErrFull(err error) (ErrFull, bool) {
	if e, ok := errors.AsType[ErrFull](err); ok {
		return e, true
	}
	if e, ok := errors.AsType[*ErrFull](err); ok && e != nil {
		return *e, true
	}
	return ErrFull{}, false
}

// Code is the response error code for logs: the service's errs.Err (including inside ErrFull),
// service_not_available for anything else.
func Code(err error) string {
	if fullErr, ok := AsErrFull(err); ok {
		err = fullErr.Err
	}
	if baseErr, ok := AsErr(err); ok {
		return baseErr.Error()
	}
	return ServiceNA.Error()
}

// IsFailure reports whether the code is a service failure (dependency down, unexpected error, not configured)
// rather than a business refusal (not found, bad request, no permission).
func IsFailure(code string) bool {
	switch Err(code) {
	case ServiceNA, NotImplemented, InvalidConfig:
		return true
	}
	return false
}

// HttpStatus maps an error code to an HTTP status: not signed in -> 401, no permission -> 403, else 400
func HttpStatus(code string) int {
	switch Err(code) {
	case NotAuthorized:
		return 401
	case NoPermission:
		return 403
	}
	return 400
}
