package valid

type CommonBody struct {
	Username string `json:"username"`
}

type CommonHeader struct {
	AcceptLanguage string `binder:"Accept-Language" validate:"required"`
}

type CommonPath struct {
	Id string `binder:"id"`
}

type CommonQuery struct {
	Page   int `binder:"page"`
	Amount int `binder:"amount"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

type Users []User

// @route /example/{id}
// @method POST
// @description This is a description test
// @summary This is a summary test
// @operationId operation-id-test
// @tags "Test Number 1" "Test Number 2" TestNumber3
type ExampleRequest struct {
	Body struct {
		CommonBody
		Users Users `json:"users" validate:"required"`
	}

	Path struct {
		CommonPath
	}

	Query struct {
		CommonQuery
		Types []string `binder:"types" validate:"required"`
	}

	Header struct {
		CommonHeader
		Version string `binder:"Version"`
	}
}
