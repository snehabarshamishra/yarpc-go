// Package types are for objects in the yarpctest API that implement
// multiple interfaces.  So if we want to reuse a function like "Name" across
// the Service and Procedure option patterns, we need to have specific structs
// that implement both of those types.  These structs also need to be public.
// By putting these structs into a single package, we can remove unusable
// functions from the exposed yarpctest package and hide them in a sub package.
package types
