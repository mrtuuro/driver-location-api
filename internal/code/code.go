package code

import "net/http"

const (
	SuccessHealthCheck = "HEALTHCHECK_SUCCESS"

	SuccessOperationCompleted = "OPERATION_COMPLETED"
	SuccessDriversCreated     = "DRIVER_LOCATIONS_CREATED"
	SuccessNearbyDriverFound  = "NEARBY_DRIVER_FOUND"
)

var SuccessMessages = map[string]string{

	SuccessHealthCheck:        "Healtcheck success.",
	SuccessOperationCompleted: "Operation completed successfully.",
	SuccessDriversCreated:     "Driver locations successfully created.",
	SuccessNearbyDriverFound:  "Driver nearby found.",
}

var StatusCodes = map[string]int{

	SuccessHealthCheck:        http.StatusOK,
	SuccessOperationCompleted: http.StatusOK,
	SuccessDriversCreated:     http.StatusCreated,
	SuccessNearbyDriverFound:  http.StatusOK,

	// ----- DRIVER ---------
	ErrNearbyDriverNotFound: http.StatusNotFound,

	// ───── AUTHENTICATION ─────
	ErrAuthInvalidCredentials: http.StatusUnauthorized, // 401
	ErrAuthInvalidToken:       http.StatusUnauthorized, // 401
	ErrAuthInvalidProtocol:    http.StatusUnauthorized,
	ErrAuthMissingToken:       http.StatusUnauthorized,    // 401
	ErrAuthForbidden:          http.StatusForbidden,       // 403
	ErrAuthTooManyAttempts:    http.StatusTooManyRequests, // 429

	// ───── VALIDATION ─────
	ErrValidationFailed: http.StatusBadRequest, // 400

	// ───── SYSTEM / INTERNAL ─────
	ErrSystemInternal:     http.StatusInternalServerError, // 500
	ErrSystemDBFailure:    http.StatusInternalServerError, // 500
	ErrSystemCacheFailure: http.StatusInternalServerError, // 500
	ErrSystemTimeout:      http.StatusGatewayTimeout,      // 504
}

const (
	ErrNearbyDriverNotFound = "DRIVER_NEARBY_NOT_FOUND"
)

// Authentication
const (
	ErrAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	ErrAuthInvalidToken       = "AUTH_INVALID_TOKEN"
	ErrAuthTokenRevoked       = "AUTH_TOKEN_REVOKED"
	ErrAuthMissingToken       = "AUTH_MISSING_TOKEN"
	ErrAuthForbidden          = "AUTH_FORBIDDEN"
	ErrAuthTooManyAttempts    = "AUTH_TOO_MANY_ATTEMPTS"
	ErrAuthInvalidProtocol    = "AUTH_INVALID_PROTOCOL"
)

// Validation
const (
	ErrValidationFailed         = "VALIDATION_FAILED"
	ErrInvalidJSON              = "INVALID_JSON"
	ErrInvalidLatitudeLongitude = "INVALID_LAT_LONG"
)

// System / Internal
const (
	ErrSystemInternal     = "SYSTEM_INTERNAL"
	ErrSystemDBFailure    = "SYSTEM_DB_FAILURE"
	ErrSystemCacheFailure = "SYSTEM_CACHE_FAILURE"
	ErrSystemTimeout      = "SYSTEM_TIMEOUT"
)

var ErrorMessages = map[string]string{
	// Authentication
	ErrAuthInvalidCredentials: "Invalid email or password.",
	ErrAuthInvalidToken:       "Access token is invalid.",
	ErrAuthTokenRevoked:       "Access token has been revoked.",
	ErrAuthMissingToken:       "Authentication token is missing.",
	ErrAuthForbidden:          "You do not have permission to access this resource.",
	ErrAuthTooManyAttempts:    "Too many failed attempts. Please try again later.",

	// Driver
	ErrNearbyDriverNotFound: "No driver found nearby.",

	// Validation
	ErrValidationFailed:         "Validation failed for input.",
	ErrInvalidJSON:              "Bad request for json body.",
	ErrInvalidLatitudeLongitude: "Driver has invalid coordinates",

	// System
	ErrSystemInternal:     "An unexpected error occurred. Please try again later.",
	ErrSystemDBFailure:    "Database error occurred.",
	ErrSystemCacheFailure: "Cache service is unavailable.",
	ErrSystemTimeout:      "The request timed out. Please try again.",
}

func GetErrorMessage(code string) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "An unknown error occurred."
}

func GetSuccessMessage(code string) string {
	if msg, ok := SuccessMessages[code]; ok {
		return msg
	}
	return "Operation completed successfully."
}

func GetStatusCode(code string) int {
	if status, ok := StatusCodes[code]; ok {
		return status
	}
	return http.StatusOK // Fallback for success, or 500 for error if needed
}
