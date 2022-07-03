package example

import "github.com/labstack/echo/v4"

// @route /duplicate
// @method GET
// @route /duplicate2
type DuplicateAttributeError struct {
	Path struct{}
}

func (s *DuplicateAttributeError) Handle(c echo.Context) error {
	return nil
}
