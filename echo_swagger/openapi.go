package echo_swagger

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/fatih/structtag"
)

// This is the root document object of the OpenAPI document.
type OpenAPI struct {
	// REQUIRED. This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses. The openapi field SHOULD be used by tooling specifications and clients to interpret the OpenAPI document. This is not related to the API info.version string.
	OpenAPI string `yaml:"openapi,omitempty" validate:"required"`

	// REQUIRED. Provides metadata about the API. The metadata MAY be used by tooling as required.
	Info Info `yaml:"info,omitempty" validate:"required"`

	// An array of Server Objects, which provide connectivity information to a target server. If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
	Servers []Server `yaml:"servers,omitempty"`

	// REQUIRED. The available paths and operations for the API.
	Paths map[string]*Path `yaml:"paths,omitempty" validate:"required"`

	// An element to hold various schemas for the specification.
	Components Components `yaml:"components,omitempty"`

	// A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition. To make security optional, an empty security requirement ({}) can be included in the array.
	Security []SecurityRequirement `yaml:"security,omitempty"`

	// A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
	Tags []Tag `yaml:"tags,omitempty"`

	// Additional external documentation.
	ExternalDocs ExternalDocumentation `yaml:"externalDocs,omitempty"`
}

// The object provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
type Info struct {
	// REQUIRED. The title of the API.
	Title string `yaml:"title,omitempty" validate:"required"`

	// A short description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// A URL to the Terms of Service for the API. MUST be in the format of a URL.
	TermsOfService string `yaml:"termsOfService,omitempty"`

	// The contact information for the exposed API.
	Contact InfoContact `yaml:"contact,omitempty"`

	// The license information for the exposed API.
	License InfoLicense `yaml:"license,omitempty"`

	// REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `yaml:"version,omitempty"`
}

// Contact information for the exposed API.
type InfoContact struct {
	// The identifying name of the contact person/organization.
	Name string `yaml:"name,omitempty"`

	// The URL pointing to the contact information. MUST be in the format of a URL.
	URL string `yaml:"url,omitempty"`

	// The email address of the contact person/organization. MUST be in the format of an email address.
	Email string `yaml:"email,omitempty"`
}

// License information for the exposed API.
type InfoLicense struct {
	// REQUIRED. The license name used for the API.
	Name string `yaml:"name,omitempty"`

	// A URL to the license used for the API. MUST be in the format of a URL.
	URL string `yaml:"url,omitempty"`
}

// An object representing a Server.
type Server struct {
	// REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
	URL string `yaml:"url,omitempty" validate:"required"`

	// An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// A map between a variable name and its value. The value is used for substitution in the server's URL template.
	Variables map[string]ServerVariable `yaml:"variables,omitempty"`
}

// An object representing a Server Variable for server URL template substitution.
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options are from a limited set. The array SHOULD NOT be empty.
	Enum []string `yaml:"enum,omitempty"`

	// REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is not supplied. Note this behavior is different than the Schema Object's treatment of default values, because in those cases parameter values are optional. If the enum is defined, the value SHOULD exist in the enum's values.
	Default string `yaml:"default,omitempty" validate:"required"`

	// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`
}

// Describes the operations available on a single path. A Path Item MAY be empty, due to ACL constraints. The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available.
type Path struct {
	// Allows for an external definition of this path item. The referenced structure MUST be in the format of a Path Item Object. In case a Path Item Object field appears both in the defined object and the referenced object, the behavior is undefined.
	Reference string `yaml:"$ref,omitempty"`

	// An optional, string summary, intended to apply to all operations in this path.
	Summary string `yaml:"summary,omitempty"`

	// An optional, string description, intended to apply to all operations in this path. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// A definition of a GET operation on this path.
	Get *Operation `yaml:"get,omitempty"`

	// A definition of a PUT operation on this path.
	Put *Operation `yaml:"put,omitempty"`

	// A definition of a POST operation on this path.
	Post *Operation `yaml:"post,omitempty"`

	// A definition of a DELETE operation on this path.
	Delete *Operation `yaml:"delete,omitempty"`

	// A definition of a OPTIONS operation on this path.
	Options *Operation `yaml:"options,omitempty"`

	// A definition of a HEAD operation on this path.
	Head *Operation `yaml:"head,omitempty"`

	// A definition of a PATCH operation on this path.
	Patch *Operation `yaml:"patch,omitempty"`

	// A definition of a TRACE operation on this path.
	Trace *Operation `yaml:"trace,omitempty"`

	// An alternative server array to service all operations in this path.
	Servers []Server `yaml:"servers,omitempty"`

	// A list of parameters that are applicable for all the operations described under this path. These parameters can be overridden at the operation level, but cannot be removed there. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters []Parameter `yaml:"parameters,omitempty"`
}

func (p *Path) GetOperationByMethod(method string) (*Operation, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return p.Get, nil

	case "PUT":
		return p.Put, nil

	case "POST":
		return p.Post, nil

	case "DELETE":
		return p.Delete, nil

	case "OPTIONS":
		return p.Options, nil

	case "HEAD":
		return p.Head, nil

	case "PATCH":
		return p.Patch, nil

	case "TRACE":
		return p.Trace, nil

	default:
		return nil, InvalidMethodError{Method: method}
	}
}

func (p *Path) SetOperationByMethod(method string, operation *Operation) error {
	currentOperation, err := p.GetOperationByMethod(method)
	if err != nil {
		return err
	} else if currentOperation != nil {
		return DuplicateMethodError{Method: method}
	}

	switch strings.ToUpper(method) {
	case "GET":
		p.Get = operation
		return nil

	case "PUT":
		p.Put = operation
		return nil

	case "POST":
		p.Post = operation
		return nil

	case "DELETE":
		p.Delete = operation
		return nil

	case "OPTIONS":
		p.Options = operation
		return nil

	case "HEAD":
		p.Head = operation
		return nil

	case "PATCH":
		p.Patch = operation
		return nil

	case "TRACE":
		p.Trace = operation
		return nil

	default:
		return InvalidMethodError{Method: method}
	}
}

// Describes a single API operation on a path.
type Operation struct {
	// A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Tags []string `yaml:"tags,omitempty"`

	// A short summary of what the operation does.
	Summary string `yaml:"summary,omitempty"`

	// A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// Additional external documentation for this operation.
	ExternalDocs ExternalDocumentation `yaml:"externalDocs,omitempty"`

	// Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is case-sensitive. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
	OperationId string `yaml:"operationId,omitempty"`

	// A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters []Parameter `yaml:"parameters,omitempty"`

	// The request body applicable for this operation. The requestBody is only supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague, requestBody SHALL be ignored by consumers.
	RequestBody RequestBody `yaml:"requestBody,omitempty"`

	// REQUIRED. The list of possible responses as they are returned from executing this operation.
	Responses map[string]Response `yaml:"responses"`

	// A map of possible out-of band callbacks related to the parent operation. The key is a unique identifier for the Callback Object. Each value in the map is a Callback Object that describes a request that may be initiated by the API provider and the expected responses.
	Callbacks map[string]Callback `yaml:"callbacks,omitempty"`

	// Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is false.
	Deprecated bool `yaml:"deprecated,omitempty"`

	// A declaration of which security mechanisms can be used for this operation. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. To make security optional, an empty security requirement ({}) can be included in the array. This definition overrides any declared top-level security. To remove a top-level security declaration, an empty array can be used.
	Security []SecurityRequirement `yaml:"security,omitempty"`

	// An alternative server array to service this operation. If an alternative server object is specified at the Path Item Object or Root level, it will be overridden by this value.
	Servers []Server `yaml:"servers,omitempty"`
}

func (operation *Operation) AddParameter(in ParameterLocation, parameter *Parameter) error {
	parameter.SetLocation(in)
	operation.Parameters = append(operation.Parameters, *parameter)
	return nil
}

func (operation *Operation) AddResponse(code string, response *Response) error {
	if operation.Responses == nil {
		operation.Responses = make(map[string]Response)
	} else if _, ok := operation.Responses[code]; ok {
		return DuplicateResponseError{StatusCode: code}
	}

	operation.Responses[code] = *response
	return nil
}

type ParameterLocation string

const (
	ParameterLocationQuery  ParameterLocation = "query"
	ParameterLocationHeader ParameterLocation = "header"
	ParameterLocationCookie ParameterLocation = "cookie"
	ParameterLocationPath   ParameterLocation = "path"
)

// Describes a single operation parameter.
// A unique parameter is defined by a combination of a name and location.
//
// There are four possible parameter locations specified by the in field:
//
// 1. Used together with Path Templating, where the parameter value is actually part of the operation's URL. This does not include the host or base path of the API. For example, in /items/{itemId}, the path parameter is itemId.
//
// 2. query - Parameters that are appended to the URL. For example, in /items?id=###, the query parameter is id.
//
// 3. header - Custom headers that are expected as part of the request. Note that RFC7230 states header names are case insensitive.
//
// 4. cookie - Used to pass a specific cookie value to the API.
type Parameter struct {
	// A simple object to allow referencing other components in the specification, internally and externally.
	Reference string `yaml:"$ref,omitempty" validate:"required_without=Name In Description Required Deprecated AllowEmptyValue"`

	// REQUIRED. The name of the parameter. Parameter names are case sensitive.
	//
	// * If in is "path", the name field MUST correspond to a template expression occurring within the path field in the Paths Object. See Path Templating for further information.
	//
	// If in is "header" and the name field is "Accept", "Content-Type" or "Authorization", the parameter definition SHALL be ignored.
	//
	// For all other cases, the name corresponds to the parameter name used by the in property.
	Name string `yaml:"name,omitempty"`

	// REQUIRED. The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	In ParameterLocation `yaml:"in,omitempty" validate:"oneof=query header path cookie"`

	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Required bool `yaml:"required"`

	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `yaml:"deprecated,omitempty"`

	// Sets the ability to pass empty-valued parameters. This is valid only for query parameters and allows sending a parameter with an empty value. Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored. Use of this property is NOT RECOMMENDED, as it is likely to be removed in a later revision.
	AllowEmptyValue bool `yaml:"allowEmptyValue,omitempty"`

	// The schema defining the content of the request parameter.
	Schema Schema `yaml:"schema,omitempty"`
}

func (parameter *Parameter) SetLocation(location ParameterLocation) {
	parameter.In = location

	if location == ParameterLocationPath {
		parameter.Required = true
	}
}

// Allows referencing an external resource for extended documentation.
type ExternalDocumentation struct {
	// A short description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// REQUIRED. The URL for the target documentation. Value MUST be in the format of a URL.
	URL string `yaml:"url,omitempty" validate:"required"`
}

// Describes a single request body.
type RequestBody struct {
	// A simple object to allow referencing other components in the specification, internally and externally.
	Reference string `yaml:"$ref,omitempty" validate:"required_without=Content Description Required"`

	// A brief description of the request body. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// REQUIRED. The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]MediaType `yaml:"content,omitempty"`

	// Determines if the request body is required in the request. Defaults to false.
	Required bool `yaml:"required,omitempty"`
}

// Each Media Type Object provides schema and examples for the media type identified by its key.
type MediaType struct {
	// The schema defining the content of the request, response, or parameter.
	Schema Schema `yaml:"schema,omitempty"`

	// Example of the media type. The example object SHOULD be in the correct format as specified by the media type. The example field is mutually exclusive of the examples field. Furthermore, if referencing a schema which contains an example, the example value SHALL override the example provided by the schema./
	Example interface{} `yaml:"example,omitempty"`

	// Examples of the media type. Each example object SHOULD match the media type and specified schema if present. The examples field is mutually exclusive of the example field. Furthermore, if referencing a schema which contains an example, the examples value SHALL override the example provided by the schema.
	Examples map[string]Example `yaml:"examples,omitempty"`

	// A map between a property name and its encoding information. The key, being the property name, MUST exist in the schema as a property. The encoding object SHALL only apply to requestBody objects when the media type is multipart or application/x-www-form-urlencoded.
	Encoding map[string]Encoding `yaml:"encoding,omitempty"`
}

type PropertyType string

const (
	PropertyType_None    PropertyType = ""
	PropertyType_Integer PropertyType = "integer"
	PropertyType_Number  PropertyType = "number"
	PropertyType_String  PropertyType = "string"
	PropertyType_Boolean PropertyType = "boolean"
	PropertyType_Array   PropertyType = "array"
	PropertyType_Object  PropertyType = "object"
	PropertyType_Map     PropertyType = "object"
)

type PropertyFormat string

const (
	PropertyFormat_None     PropertyFormat = ""
	PropertyFormat_Int32    PropertyFormat = "int32"
	PropertyFormat_Int64    PropertyFormat = "int64"
	PropertyFormat_Float    PropertyFormat = "float"
	PropertyFormat_Double   PropertyFormat = "double"
	PropertyFormat_Byte     PropertyFormat = "byte"
	PropertyFormat_Binary   PropertyFormat = "binary"
	PropertyFormat_Date     PropertyFormat = "date"
	PropertyFormat_DateTime PropertyFormat = "date-time"
)

func typeAndFormatFromKind(kind types.BasicKind) (PropertyType, PropertyFormat) {
	switch kind {
	case types.Bool:
		return PropertyType_Boolean, PropertyFormat_None

	case types.Int:
		return PropertyType_Integer, PropertyFormat_None

	case types.Int8:
		return PropertyType_Integer, PropertyFormat_None

	case types.Int16:
		return PropertyType_Integer, PropertyFormat_None

	case types.Int32:
		return PropertyType_Integer, PropertyFormat_Int32

	case types.Int64:
		return PropertyType_Integer, PropertyFormat_Int64

	case types.Uint:
		return PropertyType_Number, PropertyFormat_None

	case types.Uint8:
		return PropertyType_Number, PropertyFormat_None

	case types.Uint16:
		return PropertyType_Number, PropertyFormat_None

	case types.Uint32:
		return PropertyType_Number, PropertyFormat_None

	case types.Uint64:
		return PropertyType_Number, PropertyFormat_None

	case types.Float32:
		return PropertyType_Number, PropertyFormat_Float

	case types.Float64:
		return PropertyType_Number, PropertyFormat_Double

	case types.String:
		return PropertyType_String, PropertyFormat_None
	}

	return PropertyType_None, PropertyFormat_None
}

type Property struct {
	// REQUIRED. The schema defining the type used for the query or form parameter.
	Type PropertyType `yaml:"type,omitempty" validate:"required"`

	// Specifies the properties of the object if the property type is "object".
	Properties map[string]Property `yaml:"properties,omitempty"`

	// Specifies the format of the type.
	Format PropertyFormat `yaml:"format,omitempty"`

	// Description about this property.
	Description string `yaml:"description,omitempty"`

	// A list of the required properties.
	RequiredProperties []string `yaml:"required,omitempty"`

	// Either the property is required or not. Used for internal use
	Required bool `yaml:"-"`

	// Specifies the items of the array if the property type is "array".
	// It's must be of type interface{} because Items is actually a property.
	Items interface{} `yaml:"items,omitempty"`

	// Used for internal use
	Name string `yaml:"-"`
}

func (p Property) String() string {
	switch p.Type {
	case PropertyType_Array:
		return fmt.Sprintf("[]%v", p.Items.(Property).String())

	case PropertyType_Object:
		data := "object ("
		for k, v := range p.Properties {
			data += fmt.Sprintf("%v: %v | ", k, v.String())
		}
		data += ")"
		return data

	default:
		return string(p.Type)
	}
}

func (p *Property) ParseTags(data string, nameTag string, fieldName string) error {
	p.Name = fieldName

	if data == "" {
		return nil
	}

	tags, err := structtag.Parse(data)
	if err != nil {
		return err
	}

	name, err := tags.Get(nameTag)
	if err == nil {
		p.Name = name.Name
	}

	required, err := tags.Get(ValidateTag)
	if err == nil && required.Name == ValidateRequiredValue {
		p.Required = true
	}

	return nil
}

// The Schema Object allows the definition of input and output data types. These types can be objects, but also primitives and arrays. This object is an extended subset of the JSON Schema Specification Wright Draft 00.
type Schema struct {
	// A simple object to allow referencing other components in the specification, internally and externally.
	Reference string `yaml:"$ref,omitempty" validate:"required_without=Nullable Discriminator ReadOnly WriteOnly ExternalDocumentation Example Deprecated"`

	// A true value adds "null" to the allowed type specified by the type keyword, only if type is explicitly defined within the same Schema Object. Other Schema Object constraints retain their defined behavior, and therefore may disallow the use of null as a value. A false value leaves the specified or default type unmodified. The default value is false.
	Nullable bool `yaml:"nullable,omitempty"`

	// Adds support for polymorphism. The discriminator is an object name that is used to differentiate between other schemas which may satisfy the payload description. See Composition and Inheritance for more details.
	Discriminator Discriminator `yaml:"discriminator,omitempty"`

	// Relevant only for Schema "properties" definitions. Declares the property as "read only". This means that it MAY be sent as part of a response but SHOULD NOT be sent as part of the request. If the property is marked as readOnly being true and is in the required list, the required will take effect on the response only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	ReadOnly bool `yaml:"readOnly,omitempty"`

	// Relevant only for Schema "properties" definitions. Declares the property as "write only". Therefore, it MAY be sent as part of a request but SHOULD NOT be sent as part of the response. If the property is marked as writeOnly being true and is in the required list, the required will take effect on the request only. A property MUST NOT be marked as both readOnly and writeOnly being true. Default value is false.
	WriteOnly bool `yaml:"writeOnly,omitempty"`

	// Additional external documentation for this schema.
	ExternalDocumentation ExternalDocumentation `yaml:"externalDocs,omitempty"`

	// A free-form property to include an example of an instance for this schema. To represent examples that cannot be naturally represented in JSON or YAML, a string value can be used to contain the example with escaping where necessary.
	Example interface{} `yaml:"example,omitempty"`

	// Specifies that a schema is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `yaml:"deprecated,omitempty"`

	// Specifies the type of the object.
	Property `yaml:",inline"`
}

// When request bodies or response payloads may be one of a number of different schemas, a discriminator object can be used to aid in serialization, deserialization, and validation. The discriminator is a specific object in a schema which is used to inform the consumer of the specification of an alternative schema based on the value associated with it.
type Discriminator struct {
	// REQUIRED. The name of the property in the payload that will hold the discriminator value.
	PropertyName string `yaml:"propertyName,omitempty" validate:"required"`

	// An object to hold mappings between payload values and schema names or references.
	Mapping map[string]string `yaml:"mapping,omitempty"`
}

// In all cases, the example value is expected to be compatible with the type schema of its associated value. Tooling implementations MAY choose to validate compatibility automatically, and reject the example value(s) if incompatible.
type Example struct {
	// A simple object to allow referencing other components in the specification, internally and externally.
	Reference string `yaml:"$ref,omitempty" validate:"required_without=Summary Description Value ExternalValue"`

	// Short description for the example.
	Summary string `yaml:"summary,omitempty"`

	// Long description for the example. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// Embedded literal example. The value field and externalValue field are mutually exclusive. To represent examples of media types that cannot naturally represented in JSON or YAML, use a string value to contain the example, escaping where necessary.
	Value interface{} `yaml:"value,omitempty"`

	// A URL that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON or YAML documents. The value field and externalValue field are mutually exclusive.
	ExternalValue string `yaml:"externalValue,omitempty"`
}

// A single encoding definition applied to a single schema property.
type Encoding struct {
	// The Content-Type for encoding a specific property. Default value depends on the property type: for string with format being binary – application/octet-stream; for other primitive types – text/plain; for object - application/json; for array – the default is defined based on the inner type. The value can be a specific media type (e.g. application/json), a wildcard media type (e.g. image/*), or a comma-separated list of the two types.
	ContentType string `yaml:"contentType,omitempty"`

	// A map allowing additional information to be provided as headers, for example Content-Disposition. Content-Type is described separately and SHALL be ignored in this section. This property SHALL be ignored if the request body media type is not a multipart.
	Headers map[string]Header `yaml:"headers,omitempty"`

	// Describes how a specific property value will be serialized depending on its type. See Parameter Object for details on the style property. The behavior follows the same values as query parameters, including default values. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	Style string `yaml:"style,omitempty"`

	// When this is true, property values of type array or object generate separate parameters for each value of the array, or key-value-pair of the map. For other types of properties this property has no effect. When style is form, the default value is true. For all other styles, the default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	Explode bool `yaml:"explode,omitempty"`

	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 :/?#[]@!$&'()*+,;= to be included without percent-encoding. The default value is false. This property SHALL be ignored if the request body media type is not application/x-www-form-urlencoded.
	AllowReserved bool `yaml:"allowReserved,omitempty"`
}

// The Header Object follows the structure of the Parameter Object with the following changes:
//
// 1. name MUST NOT be specified, it is given in the corresponding headers map.
//
// 2. in MUST NOT be specified, it is implicitly in header.
//
// 3. All traits that are affected by the location MUST be applicable to a location of header (for example, style).
type Header struct {
	// A simple object to allow referencing other components in the specification, internally and externally.
	Reference string `yaml:"$ref,omitempty" validate:"required_without=Description Required Deprecated AllowEmptyValue"`

	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// Determines whether this parameter is mandatory. If the parameter location is "path", this property is REQUIRED and its value MUST be true. Otherwise, the property MAY be included and its default value is false.
	Required bool `yaml:"required,omitempty"`

	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is false.
	Deprecated bool `yaml:"deprecated,omitempty"`

	// Sets the ability to pass empty-valued parameters. This is valid only for query parameters and allows sending a parameter with an empty value. Default value is false. If style is used, and if behavior is n/a (cannot be serialized), the value of allowEmptyValue SHALL be ignored. Use of this property is NOT RECOMMENDED, as it is likely to be removed in a later revision.
	AllowEmptyValue bool `yaml:"allowEmptyValue,omitempty"`
}

// Describes a single response from an API Operation, including design-time, static links to operations based on the response.
type Response struct {
	// REQUIRED. A short description of the response. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty" validate:"required"`

	// Maps a header name to its definition. RFC7230 states header names are case insensitive. If a response header is defined with the name "Content-Type", it SHALL be ignored.
	Headers map[string]Header `yaml:"headers,omitempty"`

	// A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content map[string]MediaType `yaml:"content,omitempty"`

	// A map of operations links that can be followed from the response. The key of the map is a short name for the link, following the naming constraints of the names for Component Objects.
	Links map[string]Link `yaml:"links,omitempty"`
}

// The Link object represents a possible design-time link for a response. The presence of a link does not guarantee the caller's ability to successfully invoke it, rather it provides a known relationship and traversal mechanism between responses and other operations.
//
// Unlike dynamic links (i.e. links provided in the response payload), the OAS linking mechanism does not require link information in the runtime response.
//
// For computing links, and providing instructions to execute them, a runtime expression is used for accessing values in an operation and using them as parameters while invoking the linked operation.
type Link struct {
	// A relative or absolute URI reference to an OAS operation. This field is mutually exclusive of the operationId field, and MUST point to an Operation Object. Relative operationRef values MAY be used to locate an existing Operation Object in the OpenAPI definition.
	OperationReference string `yaml:"operationRef,omitempty"`

	// The name of an existing, resolvable OAS operation, as defined with a unique operationId. This field is mutually exclusive of the operationRef field.
	OperationId string `yaml:"operationId,omitempty"`

	// A map representing parameters to pass to an operation as specified with operationId or identified via operationRef. The key is the parameter name to be used, whereas the value can be a constant or an expression to be evaluated and passed to the linked operation. The parameter name can be qualified using the parameter location [{in}.]{name} for operations that use the same parameter name in different locations (e.g. path.id).
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`

	// A literal value or {expression} to use as a request body when calling the target operation.
	RequestBody interface{} `yaml:"requestBody,omitempty"`

	// A description of the link. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// A server object to be used by the target operation.
	Server Server `yaml:"server,omitempty"`
}

// A map of possible out-of band callbacks related to the parent operation. Each value in the map is a Path Item Object that describes a set of requests that may be initiated by the API provider and the expected responses. The key value used to identify the path item object is an expression, evaluated at runtime, that identifies a URL to use for the callback operation.
type Callback map[string]Path

// Lists the required security schemes to execute this operation. The name used for each property MUST correspond to a security scheme declared in the Security Schemes under the Components Object.
//
// Security Requirement Objects that contain multiple schemes require that all schemes MUST be satisfied for a request to be authorized. This enables support for scenarios where multiple query parameters or HTTP headers are required to convey security information.
//
// When a list of Security Requirement Objects is defined on the OpenAPI Object or Operation Object, only one of the Security Requirement Objects in the list needs to be satisfied to authorize the request.
//
// Each name MUST correspond to a security scheme which is declared in the Security Schemes under the Components Object. If the security scheme is of type "oauth2" or "openIdConnect", then the value is a list of scope names required for the execution, and the list MAY be empty if authorization does not require a specified scope. For other security scheme types, the array MUST be empty.
type SecurityRequirement map[string][]string

// Holds a set of reusable objects for different aspects of the OAS. All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas map[string]Schema `yaml:"schemas,omitempty"`

	// An object to hold reusable Response Objects.
	Responses map[string]Response `yaml:"responses,omitempty"`

	// An object to hold reusable Parameter Objects.
	Parameters map[string]Parameter `yaml:"parameters,omitempty"`

	// An object to hold reusable Example Objects.
	Examples map[string]interface{} `yaml:"examples,omitempty"`

	// An object to hold reusable Request Body Objects.
	RequestBodies map[string]RequestBody `yaml:"requestBodies,omitempty"`

	// An object to hold reusable Header Objects.
	Headers map[string]Header `yaml:"headers,omitempty"`

	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes map[string]SecurityScheme `yaml:"securitySchemes,omitempty"`

	// An object to hold reusable Link Objects.
	Links map[string]Link `yaml:"links,omitempty"`

	// An object to hold reusable Callback Objects.
	Callbacks map[string]Callback `yaml:"callbacks,omitempty"`
}

// Defines a security scheme that can be used by the operations. Supported schemes are HTTP authentication, an API key (either as a header, a cookie parameter or as a query parameter), OAuth2's common flows (implicit, password, client credentials and authorization code) as defined in RFC6749, and OpenID Connect Discovery.
type SecurityScheme struct {
	// REQUIRED. The type of the security scheme. Valid values are "apiKey", "http", "oauth2", "openIdConnect".
	Type string `yaml:"type,omitempty" validate:"required"`

	// A short description for security scheme. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`
}

// Adds metadata to a single tag that is used by the Operation Object. It is not mandatory to have a Tag Object per tag defined in the Operation Object instances.
type Tag struct {
	// REQUIRED. The name of the tag.
	Name string `yaml:"name,omitempty" validate:"required"`

	// A short description for the tag. CommonMark syntax MAY be used for rich text representation.
	Description string `yaml:"description,omitempty"`

	// Additional external documentation for this tag.
	ExternalDocs ExternalDocumentation `yaml:"externalDocs,omitempty"`
}
