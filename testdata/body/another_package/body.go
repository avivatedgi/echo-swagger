package another_package

type EmbeddedBody struct {
	A string `json:"a"`
	B string `json:"b"`
	C struct {
		CNested []float32 `json:"nested"`
	} `json:"c"`
}

type X int
