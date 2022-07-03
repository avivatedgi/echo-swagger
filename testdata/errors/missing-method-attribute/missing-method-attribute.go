package example

import "github.com/labstack/echo/v4"

// @route /duplicate
type MissingRouteAttribute struct {
	Path struct{}
}

func (s *MissingRouteAttribute) Handle(c echo.Context) error {
	return nil
}
