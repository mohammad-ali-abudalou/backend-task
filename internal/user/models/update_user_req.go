package models

type UpdateUserReq struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	// To Enforce Read-Only, The Group Is Purposefully Left Out.
}
