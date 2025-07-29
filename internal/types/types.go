package types

// validator package used to validate the student request
type Student struct {
	Id    string
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   string `validate:"required"`
}
