package echo_swagger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generalInfo() Info {
	return Info{
		Title:          "Example",
		Description:    "My Description",
		TermsOfService: "Example",
		Contact: InfoContact{
			Name:  "Example",
			URL:   "http://example.com",
			Email: "example@example.com",
		},
		License: InfoLicense{
			Name: "Example",
			URL:  "http://example.com",
		},
		Version: "1,0",
	}
}

func emptyOpenapi() OpenAPI {
	openapi := NewContext().OpenAPI
	openapi.Info = generalInfo()
	return *openapi
}

func validOpenapi() OpenAPI {
	openapi := emptyOpenapi()
	openapi.Paths["/example/{id}"] = &Path{
		Post: &Operation{
			Summary:     "This is a summary test",
			Description: "This is a description test",
			OperationId: "operation-id-test",
			Tags:        []string{"Test Number 1", "Test Number 2", "TestNumber3"},
			Parameters: []Parameter{
				{
					Name:     "Accept-Language",
					In:       ParameterLocationHeader,
					Required: true,
					Schema: Schema{
						Property: Property{
							Type:     PropertyType_String,
							Format:   PropertyFormat_None,
							Required: true,
						},
					},
				},
				{
					Name:     "Version",
					In:       ParameterLocationHeader,
					Required: false,
					Schema: Schema{
						Property: Property{
							Type:     PropertyType_String,
							Format:   PropertyFormat_None,
							Required: false,
						},
					},
				},
				{
					Name:     "id",
					In:       ParameterLocationPath,
					Required: true,
					Schema: Schema{
						Property: Property{
							Type:     PropertyType_String,
							Format:   PropertyFormat_None,
							Required: true,
						},
					},
				},
				{
					Name:     "page",
					In:       ParameterLocationQuery,
					Required: false,
					Schema: Schema{
						Property: Property{
							Type:     PropertyType_Integer,
							Format:   PropertyFormat_None,
							Required: false,
						},
					},
				},
				{
					Name:     "amount",
					In:       ParameterLocationQuery,
					Required: false,
					Schema: Schema{
						Property: Property{
							Type:     PropertyType_Integer,
							Format:   PropertyFormat_None,
							Required: false,
						},
					},
				},
				{
					Name:     "types",
					In:       ParameterLocationQuery,
					Required: true,
					Schema: Schema{
						Property: Property{
							Type: PropertyType_Array,
							Items: Property{
								Type: PropertyType_String,
							},
							Required: true,
						},
					},
				},
			},
			RequestBody: RequestBody{
				Content: map[string]MediaType{
					ContentTypeJson: {
						Schema: Schema{
							Property: Property{
								Type: PropertyType_Object,
								RequiredProperties: []string{
									"users",
								},
								Properties: map[string]Property{
									"username": {
										Type:   PropertyType_String,
										Format: PropertyFormat_None,
									},

									"users": {
										Type:     PropertyType_Array,
										Required: true,
										Items: Property{
											Type: PropertyType_Object,
											Properties: map[string]Property{
												"id": {
													Type:   PropertyType_String,
													Format: PropertyFormat_None,
												},
												"username": {
													Type:   PropertyType_String,
													Format: PropertyFormat_None,
												},
												"age": {
													Type:   PropertyType_Integer,
													Format: PropertyFormat_None,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "A valid response",
					Content: map[string]MediaType{
						ContentTypeJson: {
							Schema: Schema{
								Property: Property{
									Type: PropertyType_Object,
									Properties: map[string]Property{
										"id": {
											Type:   PropertyType_String,
											Format: PropertyFormat_None,
										},
									},
								},
							},
						},
					},
				},

				"400": {
					Description: "A bad request response",
					Content: map[string]MediaType{
						ContentTypeJson: {
							Schema: Schema{
								Property: Property{
									Type: PropertyType_Object,
									Properties: map[string]Property{
										"error": {
											Type:   PropertyType_String,
											Format: PropertyFormat_None,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return openapi
}

func testProperty(assert *assert.Assertions, expected *Property, actual *Property) bool {
	assert.Equal(expected.Type, actual.Type)
	assert.Equal(expected.Format, actual.Format)
	// assert.Equal(expected.Required, actual.Required)

	if expected.Type == PropertyType_Array {
		expectedItem := expected.Items.(Property)
		actualItem := actual.Items.(Property)
		if !testProperty(assert, &expectedItem, &actualItem) {
			return false
		}
	} else if expected.Type == PropertyType_Object {
		assert.Equal(len(expected.Properties), len(actual.Properties))

		for key, expectedValue := range expected.Properties {
			actualValue, ok := actual.Properties[key]
			assert.True(ok)
			if !testProperty(assert, &expectedValue, &actualValue) {
				return false
			}
		}
	}

	return true
}

func testOperation(assert *assert.Assertions, expected *Operation, actual *Operation) bool {
	if expected == nil {
		return assert.Nil(actual)
	}

	assert.Equal(expected.Summary, actual.Summary)
	assert.Equal(expected.Description, actual.Description)
	assert.Equal(expected.OperationId, actual.OperationId)
	assert.Equal(expected.Tags, actual.Tags)

	// Test the parameters
	expectedParameters := map[string]Parameter{}

	for _, expectedParameter := range expected.Parameters {
		expectedParameters[expectedParameter.Name] = expectedParameter
	}

	for _, actualParameter := range actual.Parameters {
		expectedParameter, exists := expectedParameters[actualParameter.Name]
		if !assert.True(exists) {
			continue
		}

		assert.Equal(expectedParameter.Deprecated, actualParameter.Deprecated)
		assert.Equal(expectedParameter.Description, actualParameter.Description)
		assert.Equal(expectedParameter.Name, actualParameter.Name)
		assert.Equal(expectedParameter.Required, actualParameter.Required)
		testProperty(assert, &expectedParameter.Schema.Property, &actualParameter.Schema.Property)
	}

	// Test the responses
	assert.Equal(len(expected.Responses), len(actual.Responses))

	for responseName, expectedResponse := range expected.Responses {
		actualResponse, exists := actual.Responses[responseName]
		if !assert.True(exists) {
			return false
		}

		assert.Equal(expectedResponse.Description, actualResponse.Description)

		for contentType, expectedContent := range expectedResponse.Content {
			actualContent, exists := actualResponse.Content[contentType]
			if !assert.True(exists) {
				return false
			}

			testProperty(assert, &expectedContent.Schema.Property, &actualContent.Schema.Property)
		}
	}

	// Test the body
	assert.Equal(len(expected.RequestBody.Content), len(actual.RequestBody.Content))

	for contentType, expectedContent := range expected.RequestBody.Content {
		actualContent, exists := actual.RequestBody.Content[contentType]
		if !assert.True(exists) {
			return false
		}

		testProperty(assert, &expectedContent.Schema.Property, &actualContent.Schema.Property)
	}

	return true
}

func TestParser(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		name          string
		directory     string
		pattern       string
		operation     string
		expectedError error
		openapi       OpenAPI
	}

	testCases := []testCase{
		{
			name:          "errors with recursive pattern",
			directory:     "../testdata/errors",
			pattern:       "./...",
			operation:     "",
			expectedError: nil,
			openapi:       emptyOpenapi(),
		},
		{
			name:          "errors with normal pattern",
			directory:     "../testdata/errors",
			pattern:       ".",
			operation:     "",
			expectedError: nil,
			openapi:       emptyOpenapi(),
		},
		{
			name:          "valid",
			directory:     "../testdata/valid",
			pattern:       "./...",
			expectedError: nil,
			openapi:       validOpenapi(),
		},
	}

	for _, testCase := range testCases {
		t.Log()
		t.Log("==============================")
		t.Log("\t", testCase.name)
		t.Log("==============================")

		context := NewContext()
		context.OpenAPI.Info = generalInfo()

		err := context.ParseDirectory(testCase.directory, testCase.pattern)

		if err != nil || testCase.expectedError != nil {
			assert.True(errors.Is(err, testCase.expectedError))
			continue
		}

		assert.Equal(len(testCase.openapi.Paths), len(context.OpenAPI.Paths))
		assert.Equal(testCase.openapi.Info, context.OpenAPI.Info)

		for pathName, expectedPath := range testCase.openapi.Paths {
			actualPath, exists := context.OpenAPI.Paths[pathName]
			assert.True(exists)

			for _, operation := range []string{"get", "post", "put", "delete", "options", "head", "patch", "trace"} {
				expectedOperation, expectedErr := expectedPath.GetOperationByMethod(operation)
				actualOperation, actualErr := actualPath.GetOperationByMethod(operation)

				if expectedErr != nil {
					assert.True(errors.Is(actualErr, expectedErr))
					continue
				}

				testOperation(assert, expectedOperation, actualOperation)
			}
		}
	}
}
