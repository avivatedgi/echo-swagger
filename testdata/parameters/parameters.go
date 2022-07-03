package example

import "github.com/labstack/echo/v4"

// @route /example
// @method GET
type ParametersTestCustomTags struct {
	Path struct {
		// @description The id description
		Id string `binder:"id"`
	}

	Query struct {
		// @description The page description
		// @required
		Page int `binder:"pageIndex"`

		// @description The amount description
		// @deprecated
		Amount float64 `binder:"amount"`
	}

	Header struct {
		Example string
	}
}

func (s *ParametersTestCustomTags) Handle(c echo.Context) error {
	return nil
}
