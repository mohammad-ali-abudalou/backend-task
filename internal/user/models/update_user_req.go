package models

// Update User Req Represents The Payload To Update A User's Name Or Email.
// Group Field Is Intentionally Omitted ( Read-Only ).
type UpdateUserReq struct {
	Name  *string `json:"name,omitempty" example:"Jane Doe"`
	Email *string `json:"email,omitempty" example:"jane@example.com"`
}
