package dto

// Credential have the information for an users log in
type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
