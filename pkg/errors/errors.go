package errors

import (
	"fmt"
	"strings"

	"github.com/esfands/retpaladinbot/pkg/utils"
	"github.com/valyala/fasthttp"
)

type APIError interface {
	Error() string
	Message() string
	Code() int
	SetDetail(str string, a ...any) APIError
	SetFields(d Fields) APIError
	GetFields() Fields
	ExpectedHTTPStatus() int
	WithHTTPStatus(s int) APIError
}

type apiErrorFunc func() APIError

var (
	// Generic client errors
	ErrUnauthorized            apiErrorFunc = DefineError(10401, "Authorization Required", fasthttp.StatusUnauthorized)
	ErrInsufficientPermissions apiErrorFunc = DefineError(10403, "Insufficient Permissions", fasthttp.StatusForbidden)
	ErrBadRequest              apiErrorFunc = DefineError(10404, "Bad Request", fasthttp.StatusBadRequest)

	// Client not found
	ErrUnknownAccount apiErrorFunc = DefineError(10440, "Unknown Account", fasthttp.StatusNotFound)

	// Client type errors
	ErrValidationRejected apiErrorFunc = DefineError(10410, "Validation Rejected", fasthttp.StatusBadRequest)

	// Other client errors

	// Server errors
	ErrInternalServerError apiErrorFunc = DefineError(10500, "Internal Server Error", fasthttp.StatusInternalServerError)
	ErrFeatureDisabled     apiErrorFunc = DefineError(10555, "Feature Disabled", fasthttp.StatusForbidden)

	// Feature flag specific errors
	ErrFeatureFlagNotFound     apiErrorFunc = DefineError(10460, "Feature Flag Not Found", fasthttp.StatusNotFound)
	ErrFeatureFlagConflict     apiErrorFunc = DefineError(10461, "Feature Flag Name Conflict", fasthttp.StatusConflict)
	ErrFeatureFlagCreateFailed apiErrorFunc = DefineError(10510, "Failed to Create Feature Flag", fasthttp.StatusInternalServerError)
	ErrFeatureFlagUpdateFailed apiErrorFunc = DefineError(10511, "Failed to Update Feature Flag", fasthttp.StatusInternalServerError)
	ErrFeatureFlagDeleteFailed apiErrorFunc = DefineError(10512, "Failed to Delete Feature Flag", fasthttp.StatusInternalServerError)
	ErrFeatureFlagFetchFailed  apiErrorFunc = DefineError(10513, "Failed to Fetch Feature Flags", fasthttp.StatusInternalServerError)

	// Tag specific errors
	ErrTagNotFound     apiErrorFunc = DefineError(10450, "Tag Not Found", fasthttp.StatusNotFound)
	ErrTagNameConflict apiErrorFunc = DefineError(10451, "Tag Name Conflict", fasthttp.StatusConflict)
	ErrTagColorInvalid apiErrorFunc = DefineError(10452, "Tag Color Invalid", fasthttp.StatusBadRequest)
	ErrTagCreateFailed apiErrorFunc = DefineError(10501, "Failed to Create Tag", fasthttp.StatusInternalServerError)
	ErrTagUpdateFailed apiErrorFunc = DefineError(10502, "Failed to Update Tag", fasthttp.StatusInternalServerError)
	ErrTagDeleteFailed apiErrorFunc = DefineError(10503, "Failed to Delete Tag", fasthttp.StatusInternalServerError)
	ErrTagFetchFailed  apiErrorFunc = DefineError(10504, "Failed to Fetch Tags", fasthttp.StatusInternalServerError)

	// Event tag specific errors
	ErrEventTagNotFound     apiErrorFunc = DefineError(10470, "Event Tag Not Found", fasthttp.StatusNotFound)
	ErrEventTagAttachFailed apiErrorFunc = DefineError(10520, "Failed to Attach Event Tag", fasthttp.StatusInternalServerError)
	ErrEventTagDetachFailed apiErrorFunc = DefineError(10521, "Failed to Detach Event Tag", fasthttp.StatusInternalServerError)
	ErrEventTagFetchFailed  apiErrorFunc = DefineError(10522, "Failed to Fetch Event Tags", fasthttp.StatusInternalServerError)
	ErrEventTagCreateFailed apiErrorFunc = DefineError(10523, "Failed to Create Event Tag", fasthttp.StatusInternalServerError)
	ErrEventTagDeleteFailed apiErrorFunc = DefineError(10524, "Failed to Delete Event Tag", fasthttp.StatusInternalServerError)
)

type apiError struct {
	message            string
	code               int
	fields             Fields
	expectedHTTPStatus int
}

type Fields map[string]interface{}

func (e *apiError) Error() string {
	return fmt.Sprintf("[%d] %s", e.code, strings.ToLower(e.message))
}

func (e *apiError) Message() string {
	return e.message
}

func (e *apiError) Code() int {
	return e.code
}

func (e *apiError) SetDetail(str string, a ...any) APIError {
	e.message = e.message + ": " + utils.Ternary(len(a) > 0, fmt.Sprintf(str, a...), str)
	return e
}

func (e *apiError) SetFields(d Fields) APIError {
	e.fields = d
	return e
}

func (e *apiError) GetFields() Fields {
	return e.fields
}

func (e *apiError) ExpectedHTTPStatus() int {
	return e.expectedHTTPStatus
}

func (e *apiError) WithHTTPStatus(s int) APIError {
	e.expectedHTTPStatus = s
	return e
}

func DefineError(code int, message string, httpStatus int) func() APIError {
	return func() APIError {
		return &apiError{
			message:            message,
			code:               code,
			fields:             Fields{},
			expectedHTTPStatus: httpStatus,
		}
	}
}
