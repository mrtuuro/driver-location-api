package response

// NOTE: This structs are only for swagger documentation
// SwaggerSuccess is a concrete, non-generic shell so swag can render it.
// Replace Data with the real payload shape per endpoint (or omit when nil).
type SwaggerSuccess struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// SwaggerError mirrors APIError so swag has a concrete type.
type SwaggerError struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
}
