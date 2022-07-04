package tags

import (
	"github.com/labstack/echo/v4"
)

// @route /example/no-tags
// @method POST
type NoTags struct {
	Body struct{}
}

func (s *NoTags) Handle(c echo.Context) error {
	return nil
}

// @route /example/one-tag
// @method POST
// @tags A
type OneTag struct {
	Body struct{}
}

func (s *OneTag) Handle(c echo.Context) error {
	return nil
}

// @route /example/one-complex-tag
// @method POST
// @tags "Hello World"
type OneComplexTag struct {
	Body struct{}
}

func (s *OneComplexTag) Handle(c echo.Context) error {
	return nil
}

// @route /example/multiple-complex-tag
// @method POST
// @tags "Hello World" How "Omri Siniver" "Is" Today
type MultipleComplexTag struct {
	Body struct{}
}

func (s *MultipleComplexTag) Handle(c echo.Context) error {
	return nil
}
