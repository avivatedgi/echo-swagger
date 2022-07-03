package example

import (
	"github.com/avivatedgi/echo-swagger/tests/body/another_package"
	"github.com/labstack/echo/v4"
)

type Nested struct {
	Nested []float64 `json:"nested"`
}

// @route /influencers/{id}
// @method POST
// @summary This is a summary of the endpoint.
// @description This is a very very very long description of the endpoint.
// It also includes several lines.
// Bli bli
// bla bla
// @tags A B C
type MyExample struct {
	// @description My body description
	Body struct {
		another_package.EmbeddedBody

		// @description the name
		Name string `json:"name" validate:"required"`

		// @description a nested example
		// @required
		Nested Nested `json:"nested"`

		// @description a map example
		Mapper map[string]string `json:"mapper"`

		// @description a boolean example
		Boolean bool `json:"boolExample"`
	}
}

func (s *MyExample) Handle(c echo.Context) error {
	return nil
}
