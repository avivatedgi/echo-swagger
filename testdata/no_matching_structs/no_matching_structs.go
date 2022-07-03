package example

import (
	"github.com/labstack/echo/v4"
)

// @route /doesnt/matters
type NotMatchingBecauseDoesntHaveHandleFunction struct {
	Path struct{}
}

type NotMatchingBecauseHandleFunctionDiffers struct {
	Path struct{}
}

func (m *NotMatchingBecauseHandleFunctionDiffers) Handle(c echo.Context) (string, error) {
	return "", nil
}

type NotMatchingBecauseHandleFunctionDiffers2 struct {
	Path struct{}
}

func (m *NotMatchingBecauseHandleFunctionDiffers2) Handle() error {
	return nil
}

type NotMatchingBecauseHandleFunctionDiffers3 struct {
	Path struct{}
}

func (m *NotMatchingBecauseHandleFunctionDiffers3) Handle() {}

type NotMatchingBecauseHandleFunctionDiffers4 struct {
	Path struct{}
}

func (m *NotMatchingBecauseHandleFunctionDiffers4) Handle(another echo.Context, c echo.Context) error {
	return nil
}
