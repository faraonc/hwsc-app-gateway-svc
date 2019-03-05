package consts

import "errors"

var (
	ErrMissingBasicAuthPrefix = errors.New(`missing "Basic " prefix in "Authorization" header`)
	ErrInvalidBase64Header    = errors.New("invalid base64 in header")
	ErrInvalidBasicAuthFormat = errors.New("invalid basic auth format")
	ErrMissingAuthHeaders     = errors.New("no headers in request")
	ErrMissingAuthHeader      = errors.New("no header in request")
	ErrMultipleAuthHeaders    = errors.New("more than 1 header in request")
)
