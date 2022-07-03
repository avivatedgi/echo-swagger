package echo_swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	assert := assert.New(t)
	parser := New()

	openapi, err := parser.ParseDirectories([]string{"../testdata/parameters"})
	assert.NoError(err)
	_, exists := openapi.Paths["/example"]
	if assert.True(exists) {
		operation, err := openapi.Paths["/example"].GetOperationByMethod("GET")
		assert.NoError(err)
		assert.NotEqual(operation, nil)

		params := map[string]Parameter{
			"id":        {Name: "id", In: "path", Required: true, Description: "The id description", Schema: Schema{Property: Property{PropertyData: PropertyData{Type: PropertyType_String}}}},
			"pageIndex": {Name: "pageIndex", In: "query", Required: true, Description: "The page description", Schema: Schema{Property: Property{PropertyData: PropertyData{Type: PropertyType_Integer}}}},
			"amount":    {Name: "amount", In: "query", Deprecated: true, Description: "The amount description", Schema: Schema{Property: Property{PropertyData: PropertyData{Type: PropertyType_Number, Format: PropertyFormat_Double}}}},
			"Example":   {Name: "Example", In: "header", Required: false, Description: "", Schema: Schema{Property: Property{PropertyData: PropertyData{Type: PropertyType_String}}}},
		}

		for _, param := range operation.Parameters {
			expected, exists := params[param.Name]
			assert.True(exists)
			assert.Equal(expected.In, param.In)
			assert.Equal(expected.Required, param.Required)
			assert.Equal(expected.Deprecated, param.Deprecated)
			assert.Equal(expected.Description, param.Description)
			assert.Equal(expected.Schema.PropertyData.Type, param.Schema.PropertyData.Type)
			assert.Equal(expected.Schema.PropertyData.Format, param.Schema.PropertyData.Format)
		}
	}
}

func itemExists[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

func TestBody(t *testing.T) {
	assert := assert.New(t)
	parser := New()

	_, err := parser.ParseDirectories([]string{"../testdata/body", "../testdata/body/another_package"})
	assert.Error(err)
	assert.Equal(errorWithPackageFile("example", "example", errorWithLocation("MyExample", errorUnfoundPackage("another_package", "Body"))), err)

	_, err = parser.ParseDirectories([]string{"../testdata/body/invalid_embedded_type"})
	assert.Error(err)
	assert.Equal(errorWithPackageFile("invalid_embedded_type", "invalid_embedded_type", errorWithLocation("Example", errorInvalidEmbeddedType("Body"))), err)

	openapi, err := parser.ParseDirectories([]string{"../testdata/body/another_package", "../testdata/body"})
	if !assert.NoError(err) {
		return
	}

	_, exists := openapi.Paths["/influencers/{id}"]
	if !assert.True(exists) {
		return
	}

	operation, err := openapi.Paths["/influencers/{id}"].GetOperationByMethod("POST")
	assert.NoError(err)
	assert.NotNil(operation)

	assert.NotNil(operation.RequestBody.Content)
	assert.True(itemExists(operation.RequestBody.Content, "application/json"))
	schema := operation.RequestBody.Content["application/json"].Schema
	assert.Equal(schema.Property.Type, PropertyType_Object)

	assert.True(itemExists(schema.Property.Properties, "name"))
	assert.Equal(schema.Property.Properties["name"].PropertyData.Type, PropertyType_String)

	assert.True(itemExists(schema.Property.Properties, "a"))
	assert.Equal(schema.Property.Properties["a"].PropertyData.Type, PropertyType_String)

	assert.True(itemExists(schema.Property.Properties, "b"))
	assert.Equal(schema.Property.Properties["b"].PropertyData.Type, PropertyType_String)

	assert.True(itemExists(schema.Property.Properties, "c"))
	assert.Equal(schema.Property.Properties["c"].PropertyData.Type, PropertyType_Object)

	assert.True(itemExists(schema.Property.Properties["c"].Properties, "nested"))
	assert.Equal(schema.Property.Properties["c"].Properties["nested"].PropertyData.Type, PropertyType_Array)
	assert.Equal(schema.Property.Properties["c"].Properties["nested"].Items.Type, PropertyType_Number)
	assert.Equal(schema.Property.Properties["c"].Properties["nested"].Items.Format, PropertyFormat_Float)

	assert.True(itemExists(schema.Property.Properties, "nested"))
	assert.True(itemExists(schema.Property.Properties["nested"].Properties, "nested"))
	assert.Equal(schema.Property.Properties["nested"].Properties["nested"].PropertyData.Type, PropertyType_Array)
	assert.Equal(schema.Property.Properties["nested"].Properties["nested"].Items.Type, PropertyType_Number)
	assert.Equal(schema.Property.Properties["nested"].Properties["nested"].Items.Format, PropertyFormat_Double)

	assert.True(itemExists(schema.Property.Properties, "mapper"))
	assert.Equal(schema.Property.Properties["mapper"].Type, PropertyType_Object)

	assert.True(itemExists(schema.Property.Properties, "boolExample"))
	assert.Equal(schema.Property.Properties["boolExample"].Type, PropertyType_Boolean)
}

func TestResponses(t *testing.T) {
	assert := assert.New(t)
	parser := New()

	_, err := parser.ParseDirectories([]string{"../testdata/responses/missing_description"})
	assert.Error(err)
	assert.Equal(errorWithPackageFile("missing_description", "missing_description", errorWithLocation("MissingDescription.OKResponse", errorMissingDescription("MissingDescription"))), err)

	openapi, err := parser.ParseDirectories([]string{"../testdata/body/another_package", "../testdata/responses"})
	assert.NoError(err)

	_, exists := openapi.Paths["/responses/example"]
	assert.True(exists)

	operation, err := openapi.Paths["/responses/example"].GetOperationByMethod("GET")
	assert.NoError(err)
	assert.NotNil(operation)
	assert.NotNil(operation.Responses)

	assert.True(itemExists(operation.Responses, "200"))
	assert.Equal(operation.Responses["200"].Description, "A success response")
	assert.True(itemExists(operation.Responses["200"].Content, "application/json"))

	schema := operation.Responses["200"].Content["application/json"].Schema
	assert.Equal(schema.Property.Type, PropertyType_Object)

	assert.True(itemExists(schema.Property.Properties, "status"))
	assert.Equal(schema.Property.Properties["status"].PropertyData.Type, PropertyType_Integer)

	assert.True(itemExists(schema.Property.Properties, "user"))
	assert.Equal(schema.Property.Properties["user"].PropertyData.Type, PropertyType_Object)

	assert.True(itemExists(schema.Property.Properties["user"].Properties, "id"))
	assert.Equal(schema.Property.Properties["user"].Properties["id"].Type, PropertyType_String)

	assert.True(itemExists(schema.Property.Properties["user"].Properties, "username"))
	assert.Equal(schema.Property.Properties["user"].Properties["username"].Type, PropertyType_String)
}

func TestErrors(t *testing.T) {
	assert := assert.New(t)
	parser := New()

	_, err := parser.ParseDirectories([]string{"../testdata/errors/duplicate-attribute-error"})
	assert.EqualError(err, errorWithPackageFile("example", "example", errorDuplicateAttribute("DuplicateAttributeError", "route")).Error())

	_, err = parser.ParseDirectories([]string{"../testdata/errors/missing-route-attribute"})
	assert.EqualError(err, errorWithPackageFile("example", "example", errorMissingAttribute("MissingRouteAttribute", "route")).Error())

	_, err = parser.ParseDirectories([]string{"../testdata/errors/missing-method-attribute"})
	assert.EqualError(err, errorWithPackageFile("example", "example", errorMissingAttribute("MissingRouteAttribute", "method")).Error())
}

func TestNoMatchingStructs(t *testing.T) {
	assert := assert.New(t)
	parser := New()

	openapi, err := parser.ParseDirectories([]string{"../testdata/no_matching_structs"})
	assert.NoError(err)
	assert.Equal(len(openapi.Paths), 0)
}
