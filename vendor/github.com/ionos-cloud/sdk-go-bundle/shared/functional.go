package shared

/*
 * This file contains a set of functional patterns that can be used to work with slices.
 * They are designed with code readability in mind. Use them to simplify common slice operations and provide error handling where appropriate.
 */

import (
	"errors"
)

// ApplyAndAggregateErrors applies the provided function for each element of the slice
// If the function returns an error, it accumulates the error and continues execution
// After all elements are processed, it returns the aggregated errors if any
func ApplyAndAggregateErrors[T any](xs []T, f func(T) error) error {
	return Fold(
		xs,
		func(errs error, x T) error {
			err := f(x)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			return errs
		},
		nil,
	)
}

// ApplyOrFail tries applying the provided function for each element of the slice
// If the function returns an error, we break execution and return the error
func ApplyOrFail[T any](xs []T, f func(T) error) error {
	return Fold(
		xs,
		func(err error, x T) error {
			// accumulate the error. If it's not nil break out of the fold
			if err != nil {
				return err
			}
			return f(x)
		},
		nil,
	)
}

// Fold (aka Reduce) accumulates the result of f into acc and returns acc by applying f over each element in the slice
func Fold[T any, Acc any](xs []T, f func(Acc, T) Acc, acc Acc) Acc {
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

// Map applies a function to each element of a slice and returns the modified slice without considering the index of each element.
func Map[T comparable, K any](s []T, f func(T) K) []K {
	return MapIdx(s, func(_ int, t T) K {
		return f(t)
	})
}

// MapIdx applies a function to each element and index of a slice, returning the modified slice with consideration of the index.
func MapIdx[V comparable, R any](s []V, f func(int, V) R) []R {
	sm := make([]R, len(s))
	for i, v := range s {
		sm[i] = f(i, v)
	}
	return sm
}

// Filter applies a function to each element of a slice, returning a new slice with only the elements for which the function returns true.
func Filter[T any](xs []T, f func(T) bool) []T {
	result := make([]T, 0, len(xs))
	for _, x := range xs {
		if f(x) {
			result = append(result, x)
		}
	}
	return result
}

// Any returns true if at least one element of a slice satisfies a given predicate function, and false otherwise.
func Any[T any](xs []T, f func(T) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}

// All returns true if all elements of a slice satisfy a given predicate function, and false otherwise.
func All[T any](xs []T, f func(T) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}
	return true
}
