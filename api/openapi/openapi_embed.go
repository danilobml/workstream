package openapi

import _ "embed"

//go:embed workstream-gateway.yaml
var specYAML []byte

//go:embed swagger-ui.html
var swaggerHTML []byte

func SpecYAML() []byte     { return specYAML }
func SwaggerHTML() []byte  { return swaggerHTML }
