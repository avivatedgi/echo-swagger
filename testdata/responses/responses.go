package example

import (
	"github.com/avivatedgi/echo-swagger/tests/body/another_package"
	"github.com/labstack/echo/v4"
)

type GenericResponse struct {
	Status int `json:"status"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

// @route /responses/example
// @method GET
type Responses struct {
	// @response 200
	// @description A success response
	OKResponse struct {
		GenericResponse
		User User `json:"user"`
	}

	// @response 400
	// @description A bad response
	BadRequestResponse struct {
		Error string `json:"error"`
		Data  struct {
			Location string `json:"location"`
		} `json:"data"`
	}

	// @response 500
	// @description A server error response
	ServerError struct {
		another_package.EmbeddedBody
		Error string `json:"error"`
	}
}

func (s *Responses) Handle(c echo.Context) error {
	return nil
}
