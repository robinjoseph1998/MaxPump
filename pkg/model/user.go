package model

// type Users struct {
// 	ID       uint   `json:"-"`
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

type Signup struct {
	FirstName string `json:"firstname" validate:"required,min=2,max=100"`
	LastName  string `json:"lastname" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required,len=10"`
	Password  string `json:"password" validate:"required,min=6"`
}
