package echo_swagger

import (
	"fmt"
)

type AttributeError struct {
	AttributeName string
}

// An error that returned whenever there is a duplicate attribute in a handler documentation
type DuplicateAttributeError struct {
	AttributeError
}

func (e DuplicateAttributeError) Error() string {
	return fmt.Sprintf("duplicate attribute `%s`", e.AttributeName)
}

// An error that returned whenever there is an attribute with invalid value in a handler documentation
type InvalidAttributeValueError struct {
	AttributeError
	Value string
}

func (e InvalidAttributeValueError) Error() string {
	return fmt.Sprintf("invalid attribute value for `%s`: `%s`", e.AttributeName, e.Value)
}

// An error that returned whenever there is a required attribute that is missing in a handler documentation
type MissingAttributeError struct {
	AttributeError
}

func (e MissingAttributeError) Error() string {
	return fmt.Sprintf("missing attribute `%s`", e.AttributeName)
}

// An error that returned whenever there is a two handlers (or more) with the exactly same path and method
type DuplicateMethodError struct {
	Method string
}

func (e DuplicateMethodError) Error() string {
	return fmt.Sprintf("duplicate method `%s`", e.Method)
}

// An error that returned whenever there is a handler with an invalid method
type InvalidMethodError struct {
	Method string
}

func (e InvalidMethodError) Error() string {
	return fmt.Sprintf("invalid method `%s`", e.Method)
}

// An invalid type error that is returned whenever there is a invalid primitive type found
type InvalidPrimitiveTypeError struct {
	TypeName string
}

func (e InvalidPrimitiveTypeError) Error() string {
	return fmt.Sprintf("invalid primitive type `%s`", e.TypeName)
}

// An error that returned whenever the type could not be found
type TypeNotFoundError struct {
	TypeName string
}

func (e TypeNotFoundError) Error() string {
	return fmt.Sprintf("type `%s` not found", e.TypeName)
}

// An error that returned whenever the type is not supported
type UnsupportedTypeError struct {
	ExpectedType string
	ActualType   string
	Embedded     bool
}

func (e UnsupportedTypeError) Error() string {
	embedded := ""
	if e.Embedded {
		embedded = "embedded "
	}

	return fmt.Sprintf("expected %stype `%s` but got `%s`", embedded, e.ExpectedType, e.ActualType)
}

// An error that returned whenever there is a duplicate response in a request handler.
type DuplicateResponseError struct {
	StatusCode string
}

func (e DuplicateResponseError) Error() string {
	return fmt.Sprintf("duplicate response `%s`", e.StatusCode)
}

func wrapError(err error, message string, args ...interface{}) error {
	return fmt.Errorf(message+": %w", append(args, err)...)
}
