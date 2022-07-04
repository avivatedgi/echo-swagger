package echo_swagger

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/fatih/structtag"
)

type (
	attributes map[string]string

	pathMethodTupple struct {
		Path   string
		Method string
	}
)

func (attrs *attributes) GetValueOrDefault(key string) string {
	if value, exists := (*attrs)[key]; exists {
		return value
	}

	return ""
}

func (attrs *attributes) HasKey(key string) bool {
	_, exists := (*attrs)[key]
	return exists
}

func (attrs *attributes) IsKeyValid(key string, f func(string) bool) bool {
	return attrs.HasKey(key) && f((*attrs)[key])
}

func extractAttributesFromComments(location, comments string) (attributes, error) {
	attributes := attributes{}
	lines := strings.Split(strings.TrimLeft(comments, " "), "\n")

	attribute := ""
	lastInsertedAttribute := ""
	data := ""

	insertAttribute := func(name, data string) error {
		if currentData, exists := attributes[name]; exists && currentData != data {
			return errorDuplicateAttribute(location, name)
		}

		attributes[name] = data
		lastInsertedAttribute = name
		return nil
	}

	for _, line := range lines {
		// If the line is empty, reset the attribute and data.
		if line == "" {
			// Check if we have a attribute that we didn't entered
			if lastInsertedAttribute != attribute {
				if err := insertAttribute(attribute, data); err != nil {
					return attributes, err
				}
			}

			attribute = ""
			data = ""
			continue
		}

		// If it is the attribute prefix, handle it as an attribute
		if strings.HasPrefix(line, "@") {
			// If there is a previous attribute, add it to the attributes
			if attribute != "" {
				if err := insertAttribute(attribute, data); err != nil {
					return nil, err
				}
			}

			// Get the first space index in the string
			space := strings.Index(line, " ")
			if space == -1 {
				// There is no space, so the whole line is the attribute I guess (without the @)
				attribute = line[1:]
			} else {
				//
				attribute = line[1:space]
				data = line[space+1:]
			}
		} else {
			if data != "" {
				data += "\n"
			}

			data += line
		}
	}

	// If there is a previous attribute, add it to the attributes
	if lastInsertedAttribute != attribute {
		if err := insertAttribute(attribute, data); err != nil {
			return nil, err
		}
	}

	return attributes, nil
}

func getPackageAndTypeNameOfExpression(t ast.Expr) (string, string) {
	data := types.ExprString(t)

	if strings.Contains(data, ".") {
		parts := strings.Split(data, ".")
		return parts[0], parts[1]
	}

	return "", data
}

type structType int

const (
	// An embedded struct, as in Bar:
	//	type Foo struct {
	//		Example string
	//	}
	//
	//	type Bar struct {
	//		Foo
	//	}
	fieldType_Embedded structType = iota

	// An inline struct, as in Bar:
	// 	type Foo struct {
	// 		Bar struct {
	//			Example int
	//		}
	//	}
	fieldType_InlineStruct

	// A struct that is declared somewhere else (not inline or embedded), as in Foo:
	//	type Foo struct {
	//		Example string
	//	}
	//
	fieldType_Declared
)

func getFieldType(field *ast.Field) structType {
	if field.Names == nil || len(field.Names) == 0 {
		return fieldType_Embedded
	} else if strings.HasPrefix(types.ExprString(field.Type), "struct{") {
		return fieldType_InlineStruct
	}

	return fieldType_Declared
}

func getFieldName(field *ast.Field, tag string) string {
	name := ""

	// Get the serialize field name out of it actual name
	if field.Names != nil && len(field.Names) > 0 {
		name = field.Names[0].Name
	}

	// If the field has the tag, use it for the name
	if field.Tag != nil {
		tags, err := structtag.Parse(field.Tag.Value[1 : len(field.Tag.Value)-1])
		if err == nil {
			if tag, err := tags.Get(tag); err == nil {
				name = tag.Name
			}
		}
	}

	return name
}

func parseStringByQuotesAndSpaces(s string) []string {
	quoted := false

	data := strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}

		return !quoted && r == ' '
	})

	for idx, val := range data {
		data[idx] = strings.Trim(val, "\"")
	}

	return data
}
