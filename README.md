# Echo Swagger

A OpenAPI generator for the [Echo](https://echo.labstack.com/) web framework that works **ONLY** with the [echo-binder](https://github.com/avivatedgi/echo-binder). The reason that it only works with the echo-binder is because it matches the binder work pattern.

This generator is a tool that parse Golang files from given directories, check for the echo-binder structures pattern, parse the structures and their documentation and returns an `OpenAPI` object with the full data needed.

## Usage

### Installation

Download [echo-swagger](https://github.com/avivatedgi/echo-swagger) by using:

```bash
get get -u github.com/avivatedgi/echo-swagger
```

### How To Use

Example:

```bash
./echo-swagger --info my_info.yaml --dir package_a/ --dir package_b/ --dir package_c/ --out openapi.yaml
```

Options:

* `--dir` Are the directories you want to parse from the OpenAPI handlers
* `--out` Is the output file to write into the OpenAPI scheme
* `--info` Is the path to the OpenAPI info file

### Info File Example

```yaml
title: Example
description: My Description
termsOfService: Example
contact:
    name: Aviv Atedgi
    url: https://www.github.com/avivatedgi
    email: aviv.atedgi2000@gmail.com
license:
    name: GNU General Public License v3.0
    url: https://www.gnu.org/licenses/gpl-3.0.en.html
version: "1.0"
```

### Format

The structure format is exactly as described in the [echo-binder](https://github.com/avivatedgi/echo-binder) documentation, but it has an extra thing: A `Handle` method and documentation attributes (starting with `@`).

<details>
  <summary>Full Example</summary>

```go
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
```

</details>

The echo swagger knows to parse the structure documentation into attributes, and formation into swagger-scheme and with both of them to build an OpenAPI scheme.

#### Struct Handler Attributes

A struct handler is a struct that implements the `Handle(c echo.Context)error` function.

* `@route` - **REQUIRED** The route to display on the OpenAPI scheme
* `@method` - **REQUIRED** The operation method of the route, can be one of those: `GET`,  `PUT`,  `POST`,  `DELETE`,  `OPTIONS`,  `HEAD`,  `PATCH`,  `TRACE`
* `@summary` - A short summary of what the operation does.
* `@description` - A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
* `@operationId` - Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
* `@deprecated` - Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is false.
* `@tags` - A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.

<details>
  <summary>Example</summary>

```go
// @route /users/{id}
// @method PUT
// @summary A little summary about this route - update user by id.
// @description Some description about this route.
// A description can also be multi-line.
// This route update the user by it's id.
// @operationId update-user-by-id
// @tags Users Updates
type UpdateUserById struct {}
```

</details>

#### Parameter Attributes

Parameters meaning is all the parameters that are related to the `Body`, `Path`, `Query` and `Header`.

* `@description` - A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
* `@required` - Determines whether this parameter is mandatory. If the parameter location is "path", this property is automatically sets to true. Also, if the field has a `validate:"required"` tag (as in the [validator](https://github.com/go-playground/validator) package) it automatically sets to true.

<details>
  <summary>Example</summary>

  ```go
  // @route /example/{id}
  // @method PUT
  type Example struct {
    Path struct {
        // This field will be required by default because it is a path field.
        Id string `binder:"id"`
    }

    Header struct {
        // This field will be required because of the required tag.
        // @required
        AcceptLanguage string `binder:"Accept-Language"`
    }

    Query struct {
        // This field will be required because of the validate:"required" tag.
        Page int `binder:"page" validate:"required"`
    }
  }

  func (e *Example) Handle(c echo.Context) error {
    return nil
  }
  ```

</details>

#### Response Attributes

A response is any struct in a request handler struct that ends with `Response` and has an `@response` attribute.

* `@response` - **REQUIRED** The matching HTTP response status code
* `@description` - **REQUIRED** A short description of the response. CommonMark syntax MAY be used for rich text representation.

<details>
  <summary>Example</summary>

  ```go
    // @response 200
    // @description A success response
    OKResponse struct {
        Status int `json:"status"`
    }

    // @response 400
    // @description A bad response
    BadRequestResponse struct {
        Error string `json:"error"`
        Data  struct {
            Location string `json:"location"`
        } `json:"data"`
    }
  ```

</details>

## Notes

* Make sure you pass the packages which in you declare the types you are using in your routes before the routes package. For example:

In the file `one-package/types.go` we declare the `CommonHeader` and `CommonQuery` structures:

```go
package one_package

type CommonHeader struct {
    AcceptVersion   string `binder:"Accept-Version"`
    AcceptLanguage  string `binder:"Accept-Language"`
}

type CommonQuery struct {
    Page    int `binder:"page"`
    Amount  int `binder:"amount"`
}
```

And in the file `another-package/routes.go` we declare a `MyRequest` struct who uses both `one_package.CommonHeader` and `one_package.CommonQuery`:

```go
package another_package

// All the needed documentation
type MyRequest struct {
    Header struct {
        one_package.CommonHeader
    }

    Query struct {
        one_package.CommonQuery
    }
}
```

So for `another_package` to identify those structures, you must first pass `one-package/` directory and only after `another-package/`.
