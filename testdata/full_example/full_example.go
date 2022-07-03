package full_example

import "github.com/labstack/echo/v4"

type CommonHeader struct {
	JWT            string `binder:"jwt"`
	AcceptLanguage string `binder:"Accept-Language"`
}

type CommonResponse struct {
	StatusCode int `json:"status_code"`
}

// @route /users/{id}
// @method PUT
// @summary A little summary about this route - update user by id.
// @description Some description about this route.
// A description can also be multi-line.
// This route update the user by it's id.
// @operationId update-user-by-id
// @tags Users
type UpdateUserById struct {
	// @description The body description
	Body struct {
		// @description The username description, will be also required because of
		// the go-playground/validator validate required tag.
		Username string `json:"username" validate:"required"`
	}

	Header struct {
		CommonHeader
	}

	Path struct {
		// @description The id path param, required by default.
		Id string `binder:"id"`
	}

	Query struct {
		// @description The age of the user
		// @required
		Age int `binder:"age"`
	}

	// @response 200
	// @description A success response
	OKResponse struct {
		CommonResponse
	}

	// @response 400
	// @description A bad response
	BadRequestResponse struct {
		CommonResponse
	}
}

func (s *UpdateUserById) Handle(c echo.Context) error {
	return nil
}
