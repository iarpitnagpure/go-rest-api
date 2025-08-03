package storage

// interfaces are a powerful way to specify the behavior of an object.
// An interface defines a set of method signatures (i.e., method names and their parameters) but does not implement them.
// Any type that implements those methods satisfies the interface—implicitly.
type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
}
