package echo_swagger

// This interface is used to implement the pattern for echo request handles.
// The echo swagger generator will use this interface to determine if a function
// is a echo request handler and will use to generate for it the swagger.
//
// Also, the echo swagger searches for a `Response` struct under the main structure,
// and it uses it to generate the specific swagger response.
// type RequestHandler interface {
// 	Handle(c echo.Context) error
// }

// Check whether a function implements our RequestHandler interface.
// func isRequestHandler(fn *ast.FuncDecl) bool {
// 	if fn.Recv == nil || fn.Name.Name != "Handle" {
// 		// Not a receiver method or not named "Handle"
// 		return false
// 	} else if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 {
// 		// There is not exactly one parameter
// 		return false
// 	} else if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
// 		// There is not exactly one return value
// 		return false
// 	}

// 	// Check the param type and the return type
// 	paramType := fmt.Sprintf("%v", fn.Type.Params.List[0].Type)
// 	returnType := fmt.Sprintf("%v", fn.Type.Results.List[0].Type)

// 	return paramType == "&{echo Context}" && returnType == "error"
// }
