package echo_swagger

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type Context struct {
	OpenAPI        *OpenAPI
	directory      string
	packagesConfig *packages.Config
	pkg            *packages.Package
	file           *ast.File
}

func NewContext() *Context {
	return &Context{
		OpenAPI: &OpenAPI{
			OpenAPI:  OpenApiVersion,
			Servers:  []Server{},
			Paths:    map[string]*Path{},
			Security: []SecurityRequirement{},
			Tags:     []Tag{},
		},
	}
}

func (context *Context) ParseDirectory(directory string, pattern string) error {
	fset := token.NewFileSet()

	context.directory = directory

	context.packagesConfig = &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,

		Fset: fset,

		Dir: directory,
	}

	pkgs, err := packages.Load(context.packagesConfig, pattern)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if err := context.parseTypesFromPackage(pkg); err != nil {
			return err
		}
	}

	return nil
}

func (context *Context) parseTypesFromPackage(pkg *packages.Package) error {
	context.pkg = pkg

	for _, file := range context.pkg.Syntax {
		context.file = file

		ast.Inspect(context.file, func(n ast.Node) bool {
			node, ok := n.(*ast.GenDecl)
			if !ok || node.Tok != token.TYPE {
				// Check if the node is a type declaration, if no let's continue the inspection.
				return true
			}

			documentation := ""
			if len(node.Specs) == 1 && node.Doc != nil {
				documentation = node.Doc.Text()
			}

			for _, spec := range node.Specs {
				spec, ok := spec.(*ast.TypeSpec)
				if !ok {
					// Check if the spec is a type spec, if no let's continue the inspection.
					continue
				}

				structType, ok := spec.Type.(*ast.StructType)
				if !ok {
					// Check if the spec is a struct type, if no let's continue the inspection.
					continue
				}

				structName := spec.Name.Name

				if !strings.HasSuffix(strings.ToLower(structName), "request") {
					// We want only structures that end with "Request".
					continue
				} else if spec.Doc == nil && documentation == "" {
					// We want only structures that have a comments to parse their attributes.
					log.Debug("Handler", structName, "found without attributes, skipping...")
					continue
				} else if spec.Doc != nil {
					documentation = spec.Doc.Text()
				}

				attributes := make(commentAttributes)
				if err := attributes.FromComments(documentation); err != nil {
					log.Warning("directory: ", context.directory, ", package: ", context.pkg.Name, ", file: ", context.file.Name.Name, " request: ", structName, " - failed to extract attributes: ", err)
					continue
				} else if err := attributes.RequiredAttributes(RouteAttribute, MethodAttribute); err != nil {
					log.Warning("directory: ", context.directory, ", package: ", context.pkg.Name, ", file: ", context.file.Name.Name, " request: ", structName, " - ", err)
					continue
				}

				if err := context.parseStruct(structName, attributes, structType); err != nil {
					log.Warning("directory: ", context.directory, ", package: ", context.pkg.Name, ", file: ", context.file.Name.Name, " request: ", structName, " - ", err)
					continue
				}
			}

			return true
		})
	}

	return nil
}

func (context *Context) parseStruct(name string, attributes commentAttributes, structType *ast.StructType) error {
	route := attributes[RouteAttribute]
	method := attributes[MethodAttribute]

	if _, exists := context.OpenAPI.Paths[route]; !exists {
		context.OpenAPI.Paths[route] = &Path{}
	}

	operation := &Operation{
		Summary:     attributes.GetOrDefault(SummaryAttribute),
		Description: attributes.GetOrDefault(DescriptionAttribute),
		OperationId: attributes.GetOrDefault(OperationIdAttribute),
		Deprecated:  attributes.HasKey(DeprecatedAttribute),
		Tags:        parseStringByQuotesAndSpaces(attributes.GetOrDefault(TagsAttribute)),
	}

	for _, field := range structType.Fields.List {
		if field.Names == nil {
			// We are searching for the Body/Header/Query/Pat attributes.
			continue
		}

		fieldName := field.Names[0].Name

		var err error

		switch fieldName {
		case BodyField:
			err = context.parseBody(operation, field)

		case HeaderField:
			err = context.parseParameter(operation, "header", field)

		case QueryField:
			err = context.parseParameter(operation, "query", field)

		case PathField:
			err = context.parseParameter(operation, "path", field)
		}

		if strings.HasSuffix(fieldName, ResponseFieldSuffix) {
			err = context.parseResponse(operation, field)
		}

		if err != nil {
			return wrapError(err, "failed to parse `%s`", fieldName)
		}
	}

	if err := context.OpenAPI.Paths[route].SetOperationByMethod(method, operation); err != nil {
		return wrapError(err, "route `%s`", route)
	}

	return nil
}

func (context *Context) parseBody(operation *Operation, field *ast.Field) error {
	t := context.pkg.TypesInfo.TypeOf(field.Type)
	if t == nil {
		return TypeNotFoundError{TypeName: types.ExprString(field.Type)}
	}

	property, err := context.parseProperty(t, JsonTag)
	if err != nil {
		return err
	}

	attributes := make(commentAttributes)
	if err := attributes.FromComments(field.Doc.Text()); err != nil {
		return wrapError(err, "failed to extract attributes")
	}

	property.Description = attributes.GetOrDefault(DescriptionAttribute)
	property.Required = attributes.HasKey(RequiredAttribute)

	operation.RequestBody = RequestBody{
		Content: map[string]MediaType{
			ContentTypeJson: {
				Schema: Schema{
					Property: *property,
				},
			},
		},
	}

	return nil
}

func (context *Context) parseParameter(operation *Operation, in ParameterLocation, field *ast.Field) error {
	t := context.pkg.TypesInfo.TypeOf(field.Type)
	if t == nil {
		return TypeNotFoundError{TypeName: types.ExprString(field.Type)}
	}

	structType, ok := t.Underlying().(*types.Struct)
	if !ok {
		return UnsupportedTypeError{ExpectedType: "struct", ActualType: t.Underlying().String()}
	}

	for fieldIndex := 0; fieldIndex < structType.NumFields(); fieldIndex++ {
		field := structType.Field(fieldIndex)
		fieldTag := structType.Tag(fieldIndex)

		// Check if the field is embedded.
		// Only embedded fields that are structures are allowed
		if field.Embedded() {
			_, ok := field.Type().Underlying().(*types.Struct)
			if !ok {
				return UnsupportedTypeError{ExpectedType: "struct", ActualType: t.Underlying().String(), Embedded: true}
			}

			embeddedProperty, err := context.parseProperty(field.Type(), BinderTag)
			if err != nil {
				return wrapError(err, "failed to parse embedded field", field.Type().Underlying().String())
			}

			for _, property := range embeddedProperty.Properties {
				if property.IgnoreProperty() {
					// Ignore the property
					continue
				}

				operation.AddParameter(in, &Parameter{
					Name:     property.Name,
					Required: property.Required,
					Schema: Schema{
						Property: property,
					},
				})
			}

			continue
		}

		fieldProperty, err := context.parseProperty(field.Type(), BinderTag)
		if err != nil {
			return wrapError(err, "failed to parse field `%s`", field.Name())
		}

		if err := fieldProperty.ParseTags(fieldTag, BinderTag, field.Name()); err != nil {
			return wrapError(err, "failed to parse field `%s` tags", field.Name())
		} else if fieldProperty.IgnoreProperty() {
			// Ignore the property
			continue
		}

		if fieldProperty.Type == PropertyType_None || fieldProperty.Type == PropertyType_Map || fieldProperty.Type == PropertyType_Object {
			return UnsupportedTypeError{ExpectedType: "primitive/slice of primitives", ActualType: field.Type().String()}
		} else if in != "query" && fieldProperty.Type == PropertyType_Array {
			return UnsupportedTypeError{ExpectedType: "primitives", ActualType: field.Type().String()}
		}

		if in == "path" {
			fieldProperty.Required = true
		}

		// Add the parameter to the operation
		operation.AddParameter(in, &Parameter{
			Name:     fieldProperty.Name,
			Required: fieldProperty.Required,
			Schema: Schema{
				Property: *fieldProperty,
			},
		})
	}

	return nil
}

func (context *Context) parseResponse(operation *Operation, field *ast.Field) error {
	t := context.pkg.TypesInfo.TypeOf(field.Type)
	if t == nil {
		return TypeNotFoundError{TypeName: types.ExprString(field.Type)}
	}

	attributes := make(commentAttributes)
	if err := attributes.FromComments(field.Doc.Text()); err != nil {
		return wrapError(err, "failed to extract attributes")
	} else if err := attributes.RequiredAttributes(ResponseAttribute, DescriptionAttribute); err != nil {
		return err
	}

	property, err := context.parseProperty(t, JsonTag)
	if err != nil {
		return err
	} else if property == nil {
		return nil
	}

	property.Description = attributes[DescriptionAttribute]
	response := Response{
		Description: attributes[DescriptionAttribute],
		Content: map[string]MediaType{
			ContentTypeJson: {
				Schema: Schema{
					Property: *property,
				},
			},
		},
	}

	if err = operation.AddResponse(attributes[ResponseAttribute], &response); err != nil {
		return err
	}

	return nil
}

func (context *Context) parseProperty(t types.Type, tag string) (*Property, error) {
	property := Property{}

	switch t := t.Underlying().(type) {
	case *types.Basic:
		property.Type, property.Format = typeAndFormatFromKind(t.Kind())

	case *types.Slice:
		property.Type = PropertyType_Array

		items, err := context.parseProperty(t.Elem(), tag)
		if err != nil {
			return nil, err
		} else if items == nil {
			return nil, nil
		}

		property.Items = *items

	case *types.Array:
		property.Type = PropertyType_Array

		items, err := context.parseProperty(t.Elem(), tag)
		if err != nil {
			return nil, err
		} else if items == nil {
			return nil, nil
		}

		property.Items = *items

	case *types.Struct:
		property.Type = PropertyType_Object
		property.Properties = make(map[string]Property)

		for fieldIndex := 0; fieldIndex < t.NumFields(); fieldIndex++ {
			field := t.Field(fieldIndex)
			fieldTag := t.Tag(fieldIndex)

			fieldProperty, err := context.parseProperty(field.Type(), tag)
			if err != nil {
				return nil, wrapError(err, "failed to parse field `%s`", field.Name())
			} else if fieldProperty == nil {
				continue
			} else if err := fieldProperty.ParseTags(fieldTag, tag, field.Name()); err != nil {
				return nil, wrapError(err, "failed to parse field `%s` tags", field.Name())
			} else if fieldProperty.IgnoreProperty() {
				continue
			}

			if field.Anonymous() {
				for name, fieldProperty := range fieldProperty.Properties {
					property.Properties[name] = fieldProperty

					if fieldProperty.Required {
						property.RequiredProperties = append(property.RequiredProperties, name)
					}
				}
			} else {
				property.Properties[fieldProperty.Name] = *fieldProperty

				if fieldProperty.Required {
					property.RequiredProperties = append(property.RequiredProperties, fieldProperty.Name)
				}
			}
		}

	case *types.Map:
		property.Type = PropertyType_Map

		vprop, err := context.parseProperty(t.Elem(), tag)
		if err != nil {
			return nil, err
		}

		property.AdditionalProperties = *vprop

	case *types.Pointer:
		return context.parseProperty(t.Elem(), tag)

	default:
		// Invalid type, return no property but also no error.
		return nil, nil
	}

	if property.IgnoreProperty() {
		return nil, nil
	}

	return &property, nil
}
