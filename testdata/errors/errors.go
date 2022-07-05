package test_errors

// @route /duplicate/route
// @method GET
// @route /duplicate/route2
type DuplicateRouteRequest struct{}

// @route /duplicate/method
// @method GET
// @method POST
type DuplicateMethodRequest struct{}

// @method GET
// @description Missing @route attribute
type MissingRouteAttribute struct{}

// @route /missing/method
// @description Missing @method attribute
type MissingMethodAttribute struct{}
