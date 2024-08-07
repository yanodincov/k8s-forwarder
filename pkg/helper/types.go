package helper

func Empty[T any]() T {
	var empty T

	return empty
}

func Ptr[T any](v T) *T {
	return &v
}
