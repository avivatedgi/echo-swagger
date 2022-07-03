package missing_description

import "github.com/labstack/echo/v4"

// @route /missing_description
// @method GET
type MissingDescription struct {
	Body struct{}

	// @response 200
	OKResponse struct{}
}

func (s *MissingDescription) Handle(c echo.Context) error {
	return nil
}
