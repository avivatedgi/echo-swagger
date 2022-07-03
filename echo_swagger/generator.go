package echo_swagger

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func GenerateOpenAPI(path string) error {
	openapi := OpenAPI{
		OpenAPI: "3.0.0",

		Info: Info{
			Title:          "Example",
			Description:    "My description",
			TermsOfService: "",
			Contact: InfoContact{
				Name:  "Aviv Atedgi",
				URL:   "https://www.github.com/avivatedgi",
				Email: "aviv.atedgi2000@gmail.com",
			},
			License: InfoLicense{
				Name: "GNU General Public License v3.0",
				URL:  "https://www.gnu.org/licenses/gpl-3.0.en.html",
			},
			Version: "1.0",
		},

		Servers: []Server{},

		Paths: map[string]*Path{
			"/": {
				Get: &Operation{
					Tags:        []string{"default"},
					Summary:     "Get the root",
					Description: "Get the root",
					OperationId: "getRoot",
					Parameters: []Parameter{
						{
							Name:        "name",
							In:          "query",
							Description: "The name of the user",
							Required:    true,
						},
					},
					Responses: map[string]Response{},
				},
			},
		},
	}

	data, err := yaml.Marshal(&openapi)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
