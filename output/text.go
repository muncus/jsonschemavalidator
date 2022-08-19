package output

import (
	"io"
	"text/template"

	"github.com/xeipuuv/gojsonschema"
)

// Text-based output of validation.
// maybe allow format customization in the future.

var textOutputTemplate = template.Must(template.New("textoutput").Parse(`
{{- .status }} {{ .source }}
{{- range $e := .errors }}
  - {{ $e }}
{{- end }}
`))

var githubOutputTemplate = template.Must(template.New("githuboutput").Parse(`
::group::{{- .status }} {{ .source }}
{{- range $e := .errors }}
::error {{ $e }}
{{- end }}
::endgroup::
`))

// StatusRune returns a rune to indicate the state of this validation.
func StatusRune(r *gojsonschema.Result) string {
	if r.Valid() {
		return "🎉"
	}
	return "❌"
}

// TODO: wrap these in our own Result object.
func templateData(source string, r *gojsonschema.Result) map[string]interface{} {
	return map[string]interface{}{
		"status": StatusRune(r),
		"source": source,
		"errors": r.Errors(),
	}
}

// TextOutput emits human readable text status messages to the provided Writer.
func TextOutput(w io.Writer, source string, r *gojsonschema.Result) error {
	err := textOutputTemplate.Execute(w, templateData(source, r))
	if err != nil {
		return err
	}
	return nil
}

func GithubOutput(w io.Writer, source string, r *gojsonschema.Result) error {
	err := githubOutputTemplate.Execute(w, templateData(source, r))
	if err != nil {
		return err
	}
	return nil
}
