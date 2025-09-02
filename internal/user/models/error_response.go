package models

// ---------------- Error Response ----------------

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Implement The Error Interface
func (e ErrorResponse) Error() string {
	return e.Message
}
