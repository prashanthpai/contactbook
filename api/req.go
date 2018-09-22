// Package api provides HTTP request and response structs.
package api

// CreateReq is request sent by client to create a new contact.
type CreateReq struct {
	Name  string `json:"name" valid:"required,runelength(1|255)"`
	Email string `json:"email" valid:"email,required,runelength(3|255)"`
	Phone string `json:"phone" valid:"runelength(0|20)"`
}

// UpdateReq is request sent by client to update an existing contact.
type UpdateReq struct {
	Name  string `json:"name" valid:"runelength(0|255)"`
	Email string `json:"email" valid:"email,runelength(0|255)"`
	Phone string `json:"phone" valid:"runelength(0|20)"`
}
