package openapi

import (
	"bytes"
	"encoding/json"
	"text/template"
)

var scalarHTML = `
<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <!-- Add your own OpenAPI/Swagger spec file URL here: -->
    <!-- Note: this includes our proxy, you can remove the following line if you do not need it -->
    <!-- data-proxy-url="https://api.scalar.com/request-proxy" -->
	<script
  	id="api-reference"
  	type="application/json">
  	{{.Spec}}
	</script>

    <!-- You can also set a full configuration object like this -->
    <!-- easier for nested objects -->
    <script>
	  var configuration = {{.Config}}

      var apiReference = document.getElementById('api-reference')
      apiReference.dataset.configuration = JSON.stringify(configuration)
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`

// Scalar returns text/HTML for serving an OpenAPI spec using the scalar library.
func Scalar(openAPI *OpenAPI, configuration map[string]any) []byte {
	scalar := template.New("scalar")
	scalar, err := scalar.Parse(scalarHTML)
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	specJSON, err := json.Marshal(openAPI)
	if err != nil {
		panic(err)
	}
	configJSON, err := json.Marshal(configuration)
	if err != nil {
		panic(err)
	}

	type ScalarConfig struct {
		Spec   string
		Config string
	}

	err = scalar.Execute(buf, &ScalarConfig{
		Spec:   string(specJSON),
		Config: string(configJSON),
	})

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
