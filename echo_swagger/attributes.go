package echo_swagger

import (
	"strings"
)

type commentAttributes map[string]string

func (attrs *commentAttributes) FromComments(comments string) error {
	(*attrs) = make(commentAttributes)

	key := ""
	data := ""

	insertAttribute := func(attr string, value string) error {
		if attr == "" {
			return nil
		}

		if err := attrs.insertAttribute(attr, data); err != nil {
			return err
		}

		key = ""
		data = ""
		return nil
	}

	lines := strings.Split(comments, "\n")

	for lineIndex, line := range lines {
		line = strings.TrimSpace(line)

		// Insert the last attribute if the line is empty or it is the last line
		if line == "" {
			if err := insertAttribute(key, data); err != nil {
				return err
			}

			continue
		}

		if strings.HasPrefix(line, "@") {
			if err := insertAttribute(key, data); err != nil {
				return err
			}

			spaceIndex := strings.Index(line, " ")
			if spaceIndex == -1 {
				key = line[1:]
			} else {
				key = line[1:spaceIndex]
				data = line[spaceIndex+1:]
			}
		} else {
			if data != "" {
				data += " "
			}

			data += line
		}

		if lineIndex == len(lines)-1 {
			if err := insertAttribute(key, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func (attrs commentAttributes) HasKey(key string) bool {
	_, exists := attrs[key]
	return exists
}

func (attrs commentAttributes) RequiredAttributes(keys ...string) error {
	for _, key := range keys {
		if !attrs.HasKey(key) {
			return MissingAttributeError{AttributeError: AttributeError{AttributeName: key}}
		}
	}

	return nil
}

func (attrs commentAttributes) IsKeyValid(key string, f func(string) bool) bool {
	if !attrs.HasKey(key) {
		return false
	}

	return f(attrs[key])
}

func (attrs commentAttributes) GetOrDefault(key string) string {
	if value, exists := attrs[key]; exists {
		return value
	}

	return ""
}

func (attrs *commentAttributes) insertAttribute(key string, value string) error {
	if data, exists := (*attrs)[key]; exists && data != value {
		return DuplicateAttributeError{AttributeError: AttributeError{AttributeName: key}}
	}

	(*attrs)[key] = value
	return nil
}
