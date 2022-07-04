package echo_swagger

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
)

type (
	ParserContext struct {
		openapi    *OpenAPI
		pkg        *ast.Package
		structures map[string]map[string]*ast.StructType
	}
)

func New() *ParserContext {
	return &ParserContext{
		structures: make(map[string]map[string]*ast.StructType),
	}
}

func (context *ParserContext) ParseDirectories(dirs []string) (*OpenAPI, error) {
	context.openapi = &OpenAPI{
		OpenAPI:  OpenApiVersion,
		Servers:  []Server{},
		Paths:    map[string]*Path{},
		Security: []SecurityRequirement{},
		Tags:     []Tag{},
	}

	for _, dir := range dirs {
		fileSet := token.NewFileSet()
		packages, err := parser.ParseDir(fileSet, dir, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		for name, pkg := range packages {
			// Update the current package
			context.pkg = pkg

			if err := context.parsePackage(); err != nil {
				return nil, errorWithPackageFile(name, err)
			}
		}
	}

	return context.openapi, nil
}

func (context *ParserContext) parsePackage() error {
	documentation := doc.New(context.pkg, "./", doc.AllDecls)

	// This map is used to store the struct names to the path and method
	// because parsing the struct documentation and struct fields are done in
	// two different places
	routesData := map[string]pathMethodTupple{}

	for _, t := range documentation.Types {
		for _, v := range t.Methods {
			if isRequestHandler(v.Decl) {
				path, method, err := context.parseHandlerStructComments(t.Name, t.Doc)
				if err != nil {
					return err
				}

				operation, err := context.openapi.Paths[path].GetOperationByMethod(method)
				if operation == nil || err != nil {
					return err
				}

				routesData[t.Name] = pathMethodTupple{Path: path, Method: method}

				// We found our Handle method, no need to continue run on this structure methods.
				break
			}
		}
	}

	var rootError error = nil

	ast.Inspect(context.pkg, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					name := spec.Name.Name

					switch spec := spec.Type.(type) {
					case *ast.StructType:
						data, ok := routesData[name]
						if !ok {
							if _, ok := context.structures[context.pkg.Name]; !ok {
								context.structures[context.pkg.Name] = make(map[string]*ast.StructType)
							}

							context.structures[context.pkg.Name][name] = spec
							continue
						}

						operation, err := context.openapi.Paths[data.Path].GetOperationByMethod(data.Method)
						if err != nil {
							continue
						}

						if err := context.parseHandlerStructFields(name, data, operation, spec); err != nil {
							rootError = err
							return false
						}

						continue
					}
				}
			}
		}

		return true
	})

	return rootError
}

func (context *ParserContext) parseHandlerStructFields(structName string, data pathMethodTupple, operation *Operation, structType *ast.StructType) error {
	for _, field := range structType.Fields.List {
		if field.Names == nil || len(field.Names) == 0 {
			// Probably an embedded field or something
			continue
		}

		fieldName := field.Names[0].Name

		attributes, err := extractAttributesFromComments(fieldName, field.Doc.Text())
		if err != nil {
			return errorWithLocation(structName, err)
		}

		fieldAsStruct, ok := field.Type.(*ast.StructType)

		switch fieldName {
		case BodyParam:
			if err := context.parseBody(structName, data, operation, field, attributes); err != nil {
				return errorWithLocation(structName, err)
			}

		case PathParam:
			fallthrough

		case QueryParam:
			fallthrough

		case HeaderParam:
			if !ok {
				return errorInvalidType(structName, fieldName)
			}

			if err := context.parseParameters(structName, fieldName, data, operation, fieldAsStruct, attributes); err != nil {
				return errorWithLocation(structName+"."+fieldName, err)
			}

		default:
			f := func(data string) bool {
				_, err := strconv.Atoi(data)
				return err == nil
			}

			if strings.HasSuffix(fieldName, "Response") && attributes.IsKeyValid(ResponseAttribute, f) {
				if err := context.parseResponses(structName, data, operation, field, attributes); err != nil {
					return errorWithLocation(structName+"."+fieldName, err)
				}
			}
		}

	}

	return nil
}

func (context *ParserContext) parseHandlerStructComments(name string, comments string) (string, string, error) {
	attributes, err := extractAttributesFromComments(name, comments)
	if err != nil {
		return "", "", err
	}

	route, exists := attributes[RouteAttribute]
	if !exists {
		return "", "", errorMissingAttribute(name, RouteAttribute)
	} else if route == "" {
		return "", "", errorInvalidAttributeValue(name, RouteAttribute, route)
	}

	method, exists := attributes[MethodAttribute]
	if !exists {
		return "", "", errorMissingAttribute(name, MethodAttribute)
	}

	summary := attributes.GetValueOrDefault(SummaryAttribute)
	description := attributes.GetValueOrDefault(DescriptionAttribute)
	operationId := attributes.GetValueOrDefault(OperationIdAttribute)
	deprecated := attributes.HasKey(DeprecatedAttribute)

	tagsValue := attributes.GetValueOrDefault(TagsAttribute)
	tags := []string{}
	if tagsValue != "" {
		tags = strings.Split(tagsValue, " ")
	}

	path, exists := context.openapi.Paths[route]
	if exists {
		operation, err := path.GetOperationByMethod(method)
		if err != nil {
			return "", "", err
		} else if operation != nil {
			return "", "", errorDuplicateOperation(name, method)
		}
	} else {
		path = &Path{}
	}

	operation := Operation{
		Summary:     summary,
		Description: description,
		OperationId: operationId,
		Deprecated:  deprecated,
		Tags:        tags,
		Responses:   map[string]Response{},
	}

	if err := path.SetOperationByMethod(method, &operation); err != nil {
		return "", "", err
	}

	context.openapi.Paths[route] = path
	return route, method, nil
}

func (context *ParserContext) parseBody(structName string, data pathMethodTupple, operation *Operation, field *ast.Field, attributes attributes) error {
	property, err := context.extractPropertyFromField(field, JsonTag)
	if err != nil {
		return err
	}

	operation.RequestBody.Required = attributes.HasKey(RequiredAttribute)
	operation.RequestBody.Description = attributes.GetValueOrDefault(DescriptionAttribute)
	operation.RequestBody.Content = map[string]MediaType{
		ContentTypeJson: {
			Schema: Schema{
				Property:   property,
				Nullable:   attributes.HasKey(NullableAttribute),
				Deprecated: attributes.HasKey(DeprecatedAttribute),
			},
		},
	}

	return nil
}

func (context *ParserContext) parseParameters(structName string, paramName string, data pathMethodTupple, operation *Operation, structType *ast.StructType, attributes attributes) error {
	for _, field := range structType.Fields.List {
		fieldName := ""
		if field.Names != nil && len(field.Names) > 0 {
			fieldName = field.Names[0].Name
		}

		serializeName := fieldName

		if field.Tag != nil {
			tags, err := structtag.Parse(field.Tag.Value[1 : len(field.Tag.Value)-1])
			if err != nil {
				return errorWithLocation(structName+"."+paramName+"."+fieldName, err)
			}

			if tag, err := tags.Get(BinderTag); err == nil {
				serializeName = tag.Name
			}
		}

		if serializeName == "" {
			// The field doesn't have a valid name
			continue
		}

		nestedFieldAttributes, err := extractAttributesFromComments(fieldName, field.Doc.Text())
		if err != nil {
			return errorWithLocation(structName+"."+fieldName, err)
		}

		if paramName == PathParam {
			nestedFieldAttributes[RequiredAttribute] = "true"
		}

		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        serializeName,
			Description: nestedFieldAttributes.GetValueOrDefault(DescriptionAttribute),
			Required:    nestedFieldAttributes.HasKey(RequiredAttribute),
			Deprecated:  nestedFieldAttributes.HasKey(DeprecatedAttribute),
			In:          strings.ToLower(paramName),
			Schema: Schema{
				// There is no need to use extractPropertyFromField because params can not be objects,
				// only primitives and arrays of primitives.
				Property: propertyFromLiteralType(types.ExprString(field.Type)),
			},
		})

	}

	return nil
}

func (context *ParserContext) parseResponses(structName string, data pathMethodTupple, operation *Operation, field *ast.Field, attributes attributes) error {
	property, err := context.extractPropertyFromField(field, JsonTag)
	if err != nil {
		return err
	}

	description, exists := attributes[DescriptionAttribute]
	if !exists {
		return errorMissingDescription(structName)
	}

	status := attributes.GetValueOrDefault(ResponseAttribute)
	operation.Responses[status] = Response{
		Description: description,
		Content: map[string]MediaType{
			ContentTypeJson: {
				Schema: Schema{
					Property: property,
				},
			},
		},
	}
	return nil
}

func (context *ParserContext) extractPropertyFromStruct(property *Property, structName string, structType *ast.StructType, tag string) error {
	property.Properties = make(map[string]Property)

	for _, field := range structType.Fields.List {
		fieldProperty, err := context.extractPropertyFromField(field, tag)
		if err != nil {
			return err
		}

		fieldName := getFieldName(field, tag)
		fieldType := getFieldType(field)

		if fieldProperty.Type != PropertyType_Object && fieldType == fieldType_Embedded {
			// If it is not an object and it is embedded, the type is invalid
			return errorInvalidEmbeddedType(structName)
		} else if fieldProperty.Type != PropertyType_Object && fieldType != fieldType_Embedded {
			// If it is not an object, and it is not an embedded struct, then let's add it to the properties
			property.Properties[fieldName] = fieldProperty.fixType()

			if fieldProperty.Required {
				property.RequiredProperties = append(property.RequiredProperties, fieldName)
			}

			continue
		} else if fieldType == fieldType_InlineStruct {
			if err := context.extractPropertyFromStruct(&fieldProperty, field.Names[0].Name, field.Type.(*ast.StructType), tag); err != nil {
				return errorWithLocation(structName+"."+fieldName, err)
			}

			property.Properties[fieldName] = fieldProperty.fixType()

			if fieldProperty.Required {
				property.RequiredProperties = append(property.RequiredProperties, fieldName)
			}

			continue
		}

		// If it is a declared or embedded struct, then we need to extract it from the structs map
		pkg, name := getPackageAndTypeNameOfExpression(field.Type)
		if pkg == "" {
			pkg = context.pkg.Name
		}

		if _, ok := context.structures[pkg]; !ok {
			return errorUnfoundPackage(pkg, structName)
		}

		structType, ok := context.structures[pkg][name]
		if !ok {
			return errorUnfoundStructInPackage(pkg, name, structName)
		}

		if err := context.extractPropertyFromStruct(&fieldProperty, fieldName, structType, tag); err != nil {
			return err
		}

		if fieldType == fieldType_Embedded {
			for key, value := range fieldProperty.Properties {
				property.Properties[key] = value.fixType()

				if value.Required {
					property.RequiredProperties = append(property.RequiredProperties, key)
				}
			}
		} else if fieldType == fieldType_Declared {
			property.Properties[fieldName] = fieldProperty.fixType()

			if fieldProperty.Required {
				property.RequiredProperties = append(property.RequiredProperties, fieldName)
			}
		}
	}

	return nil
}

func (context *ParserContext) extractPropertyFromField(field *ast.Field, tag string) (Property, error) {
	property := propertyFromLiteralType(types.ExprString(field.Type))
	fieldName := ""
	if field.Names != nil && len(field.Names) > 0 {
		fieldName = field.Names[0].Name
	}

	fieldAttributes, err := extractAttributesFromComments(fieldName, field.Doc.Text())
	if err != nil {
		return property, err
	}

	property.Description = fieldAttributes.GetValueOrDefault(DescriptionAttribute)
	property.Required = fieldAttributes.HasKey(RequiredAttribute)

	if field.Tag != nil && field.Tag.Value != "" {
		tags, err := structtag.Parse(field.Tag.Value[1 : len(field.Tag.Value)-1])
		if err != nil {
			return property, errorWithLocation(fieldName, err)
		}

		actualTag, err := tags.Get(ValidateTag)
		if err == nil && actualTag.Name == ValidateRequiredValue {
			property.Required = true
		}
	}

	structType, ok := field.Type.(*ast.StructType)
	if !ok || structType == nil || structType.Fields == nil || structType.Fields.List == nil {
		return property, err
	}

	if err := context.extractPropertyFromStruct(&property, fieldName, structType, tag); err != nil {
		return property, err
	}

	property.fixType()
	return property, nil
}
