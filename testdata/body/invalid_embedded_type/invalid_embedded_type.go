package invalid_embedded_type

import "github.com/labstack/echo/v4"

// @route /invalid_embedded_type/example
// @method GET
type Example struct {
	Body struct {
		string
	}
}

func (example *Example) Handle(c echo.Context) error {
	return nil
}
