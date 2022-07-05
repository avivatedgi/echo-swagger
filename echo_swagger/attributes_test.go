package echo_swagger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttributes(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		comments      string
		expectedAttrs commentAttributes
		expectedError error
	}

	testCases := []testCase{
		{
			comments:      "",
			expectedAttrs: commentAttributes{},
			expectedError: nil,
		},
		{
			comments: `
			@route example
			@method GET
			@description long long
			description
			`,
			expectedAttrs: commentAttributes{
				"route":       "example",
				"method":      "GET",
				"description": "long long description",
			},
			expectedError: nil,
		},
		{
			comments: `
			@route example

			@method GET

			@description long long
			description`,
			expectedAttrs: commentAttributes{
				"route":       "example",
				"method":      "GET",
				"description": "long long description",
			},
			expectedError: nil,
		},
		{
			comments: `
			@route example
			@route example2
			`,
			expectedAttrs: commentAttributes{},
			expectedError: DuplicateAttributeError{AttributeError: AttributeError{AttributeName: "route"}},
		},
	}

	for _, testCase := range testCases {
		attrs := commentAttributes{}
		if err := attrs.FromComments(testCase.comments); err != nil || testCase.expectedError != nil {
			assert.True(errors.Is(err, testCase.expectedError))
			continue
		}

		if !assert.Equal(len(testCase.expectedAttrs), len(attrs)) {
			continue
		}

		for key, expected := range testCase.expectedAttrs {
			if actual, exists := attrs[key]; assert.True(exists) {
				assert.Equal(expected, actual)
			}
		}
	}
}
