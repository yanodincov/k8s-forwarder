package helper

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func If[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

func IfFn[T any](cond bool, t, f func() T) T {
	if cond {
		return t()
	}
	return f()
}

func StaticFn[T any](val T) func() T {
	return func() T {
		return val
	}
}

func IfFnOrDef[T any](cond bool, t func() T) T {
	if !cond {
		return Empty[T]()
	}

	return t()
}
