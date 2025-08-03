package types

// validator package used to validate the student request
type Student struct {
	Id    int64
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int    `validate:"required"`
}
