package consts

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrServiceUnavailable      = errors.New("service is unavailable")
	ErrNilRequest              = errors.New("nil request")
	ErrMissingBasicAuthPrefix  = errors.New(`missing "Basic " prefix in "authorization" header`)
	ErrMissingAuthTokenPrefix  = errors.New(`missing "Auth Token " prefix in "authorization" header`)
	ErrMissingEmailTokenPrefix = errors.New(`missing "Email Token " prefix in "authorization" header`)
	ErrInvalidBase64Header     = errors.New("invalid base64 in header")
	ErrInvalidBasicAuthFormat  = errors.New("invalid basic auth format")
	ErrMissingHeadersFromCtx   = errors.New("no headers in request")
	ErrMissingAuthHeader       = errors.New(`no "authorization" header in request`)
	ErrMultipleAuthHeaders     = errors.New("more than 1 header in request")
	ErrMissingHeader           = errors.New("missing header")
	ErrMissingTag              = errors.New("missing tag")
	//ErrNilRequest                = errors.New("nil request")
	//ErrNilUserRequest            = errors.New("nil UserRequest/User")
	ErrUnableToUpdateAuthSecret = errors.New("unable to update auth secret")
	StatusUnauthenticated       = status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
)
