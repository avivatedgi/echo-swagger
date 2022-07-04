package echo_swagger

import "fmt"

func errorDuplicateAttribute(location, attribute string) error {
	return fmt.Errorf("duplicate attribute `%s` in `%s`", attribute, location)
}

func errorInvalidAttributeValue(location, attribute, value string) error {
	return fmt.Errorf("invalid attribute value for `%s` in `%s`: `%s`", attribute, location, value)
}

func errorMissingAttribute(location, attribute string) error {
	return fmt.Errorf("missing attribute `%s` in `%s`", attribute, location)
}

func errorDuplicateOperation(location, operation string) error {
	return fmt.Errorf("duplicate method `%s` in `%s`", operation, location)
}

func errorInvalidMethod(method string) error {
	return fmt.Errorf("invalid method `%s`", method)
}

func errorInvalidType(location, attribute string) error {
	return fmt.Errorf("the `%s` attribute in `%s` must be of type `struct`", attribute, location)
}

func errorWithLocation(location string, err error) error {
	return fmt.Errorf("%s: %s", location, err)
}

func errorWithPackageFile(pkg string, err error) error {
	return fmt.Errorf("error in package %s: %s", pkg, err.Error())
}

func errorUnfoundPackage(pkg string) error {
	return fmt.Errorf("unfound package `%s`", pkg)
}

func errorUnfoundStructInPackage(pkg string, structure string) error {
	return fmt.Errorf("unfound struct `%s` in package `%s`", structure, pkg)
}

func errorInvalidEmbeddedType(location string) error {
	return fmt.Errorf("can not embed a non-struct type in a struct (%s)", location)
}

func errorMissingDescription(location string) error {
	return fmt.Errorf("missing description in `%s`", location)
}
