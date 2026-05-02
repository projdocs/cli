package response

// Response is a generic envelope returned by every handler.
// Exactly one of Data and Err will be non-zero; the other will be its zero value.
type Response[T any] struct {
	Data *T     `json:"data"`
	Err  *Error `json:"error"`
}

// Error carries a human-readable message and an optional machine-readable code.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// Success constructs a Response carrying data and a nil error.
func Success[T any](data T) Response[T] {
	return Response[T]{Data: &data}
}

// Failure constructs a Response carrying an error and a nil data pointer.
func Failure(code, message string) Response[any] {
	return Response[any]{Err: &Error{Code: code, Message: message}}
}
