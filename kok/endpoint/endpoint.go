package endpoint

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/kok/gen"
	"github.com/RussellLuo/kok/reflector"
)

var (
	template = `
package {{.PkgName}}

import (
	"github.com/go-kit/kit/endpoint"
	{{- range .Imports}}
	"{{.}}"
	{{- end }}
)

{{- $srcPkgPrefix := .SrcPkgPrefix}}
{{- $interfaceName := .Interface.Name}}

{{- range .Interface.Methods}}
{{- $params := nonCtxParams .Params}}
{{- if $params}}
type {{.Name}}Request struct {
	{{- range $params}}
	{{title .Name}} {{.Type}} {{addTag .Name .Type}}
	{{- end}}
}
{{- end}}

type {{.Name}}Response struct {
	{{- range .Returns}}
	{{title .Name}} {{.Type}} {{addTag .Name .Type}}
	{{- end}}
}

{{- $errParamName := getErrParamName .Returns}}
{{- if $errParamName}}
// Failed implements endpoint.Failer.
func (r {{addAsterisks .Name}}Response) Failed() error { return r.{{title $errParamName}} }
{{- end}}

// MakeEndpointOf{{.Name}} creates the endpoint for s.{{.Name}}.
func MakeEndpointOf{{.Name}}(s {{$srcPkgPrefix}}{{$interfaceName}}) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		{{- if $params}}
		req := request.({{addAsterisks .Name}}Request)
		{{- end}}
		{{joinName .Returns ", "}} := s.{{.Name}}(
			ctx,
			{{- range $params}}
			req.{{title .Name}} {{- if .Variadic}}...{{end}},
			{{- end}}
		)
		return {{addAmpersand .Name}}Response{
			{{- range .Returns}}
			{{title .Name}}: {{.Name}},
			{{- end}}
		}, nil
	}
}
{{- end}}
`
)

type Options struct {
	SchemaPtr         bool
	SchemaTag         string
	TagKeyToSnakeCase bool
}

type Generator struct {
	opts Options
}

func New(opts Options) *Generator {
	return &Generator{opts: opts}
}

func (e *Generator) Generate(result *reflector.Result) ([]byte, error) {
	return gen.Generate(template, result, gen.Options{
		Funcs: map[string]interface{}{
			"title": strings.Title,
			"nonCtxParams": func(params []*reflector.Param) (out []*reflector.Param) {
				for _, p := range params {
					if p.Type != "context.Context" {
						out = append(out, p)
					}
				}
				return
			},
			"getErrParamName": func(params []*reflector.Param) string {
				for _, p := range params {
					if p.Type == "error" {
						return p.Name
					}
				}
				return ""
			},
			"joinName": func(returns []*reflector.Param, sep string) string {
				var names []string
				for _, r := range returns {
					names = append(names, r.Name)
				}
				return strings.Join(names, sep)
			},
			"addAsterisks": func(name string) string {
				if e.opts.SchemaPtr {
					return "*" + name
				}
				return name
			},
			"addAmpersand": func(name string) string {
				if e.opts.SchemaPtr {
					return "&" + name
				}
				return name
			},
			"addTag": func(name, typ string) string {
				if e.opts.SchemaTag == "" {
					return ""
				}

				if typ == "error" {
					name = "-"
				} else if e.opts.TagKeyToSnakeCase {
					name = ToSnakeCase(name)
				}

				return fmt.Sprintf("`%s:\"%s\"`", e.opts.SchemaTag, name)
			},
		},
		Formatters: []gen.Formatter{gen.Gofmt, gen.Goimports},
	})
}

var (
	// matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(s string) string {
	// snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake := matchAllCap.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
