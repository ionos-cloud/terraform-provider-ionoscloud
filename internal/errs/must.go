package errs

// Must panics if err is not nil, otherwise it returns x.
func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}
