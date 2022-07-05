# Echo Swagger

A OpenAPI generator for the [Echo](https://echo.labstack.com/) web framework that works **ONLY** with the [echo-binder](https://github.com/avivatedgi/echo-binder). The reason that it only works with the echo-binder is because it matches the binder work pattern.

This generator is a tool that parse Golang files from given directory, check for the echo-binder structures pattern, parse the structures and their documentation and returns an `OpenAPI` object with the full data needed.

## Usage

### Installation

Download [echo-swagger](https://github.com/avivatedgi/echo-swagger) by using:

```bash
go install github.com/avivatedgi/echo-swagger@latest
```

### How To Use

Example:

```bash
echo-swagger --info my_info.yaml --dir package/ --patern ./... --out openapi.yaml
```

Options:

* `--dir` Is the directory you want to parse from the OpenAPI specifications
* `--pattern` The pattern to scan with the packages (default: `./...`)
* `--out` Is the output file to write into the OpenAPI scheme (default: `stdout`)
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

The structure format is exactly as described in the [echo-binder](https://github.com/avivatedgi/echo-binder) documentation, but it has an extra thing: documentation attributes (starting with `@`).

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
// @tags "First Tag" SecondTag
type UpdateUserById struct {
    Body struct {
        // will be required because of the go-playground/validator validate required tag.
        Username string `json:"username" validate:"required"`
    }

    Header struct {
        // will embed the `CommonHeader` fields
        CommonHeader
    }

    Path struct {
        Id string `binder:"id"`
    }

    Query struct {
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

Parameters meaning is all the parameters that are related to the `Body`, `Path`, `Query` and `Header`. Currently, the only supported attribute is `required` and it is only valid through the `validate:"required"` tag.

<details>
  <summary>Example</summary>

  ```go
  // @route /example
  // @method PUT
  type Example struct {
    Query struct {
        // This field will be required because of the validate:"required" tag.
        Page int `binder:"page" validate:"required"`
    }
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

## TODO

* [ ] Add support for adding attributes for fields
