package endpoint

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/gen/util/openapi"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

var (
	template = annotation.FileHeader + `
package {{.PkgInfo.CurrentPkgName}}

import (
	"github.com/RussellLuo/validating/v2"
	"github.com/go-kit/kit/endpoint"
	"github.com/RussellLuo/kun/pkg/httpoption"

	{{- range .Data.Imports}}
	{{.ImportString}}
	{{- end }}
)

{{- range .DocMethods}}

{{- $params := nonCtxParams .Params .Op.Request.Params}}
{{- $hasCtxParam := hasCtxParam .Params}}

{{ interfaceWrapper .Params}}
 
{{- if $params}}
type {{.Name}}Request struct {
	{{- range $params}}
	{{title .Name}} {{.TypeString}} {{addTag .Alias .TypeString}}
	{{- end}}
}

// Validate{{.Name}}Request creates a validator for {{.Name}}Request.
func Validate{{.Name}}Request(newSchema func({{addAsterisks .Name}}Request) validating.Schema) httpoption.Validator {
	return httpoption.FuncValidator(func(value interface{}) error {
		req := value.({{addAsterisks .Name}}Request)
		return httpoption.Validate(newSchema(req))
	})
}
{{- end}}

{{if .Returns -}}

type {{.Name}}Response struct {
	{{- range .Returns}}
	{{title .Name}} {{.TypeString}} {{addTag .Name .TypeString}}
	{{- end}}
}

{{- $respBodyField := .Op.SuccessResponse.BodyField}}
{{- if $respBodyField}}
func (r {{addAsterisks .Name}}Response) Body() interface{} { return &r.{{title $respBodyField}} }
{{- else}}
func (r {{addAsterisks .Name}}Response) Body() interface{} { return r }
{{- end}}

{{- end}} {{/* if .Returns */}}

{{- $errParamName := getErrParamName .Returns}}
{{- if $errParamName}}
// Failed implements endpoint.Failer.
func (r {{addAsterisks .Name}}Response) Failed() error { return r.{{title $errParamName}} }
{{- end}}

// MakeEndpointOf{{.Name}} creates the endpoint for s.{{.Name}}.
func MakeEndpointOf{{.Name}}(s {{$.Data.SrcPkgQualifier}}{{$.Data.InterfaceName}}) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		{{- if $params}}
		req := request.({{addAsterisks .Name}}Request)
		{{- end}}

		{{- if .Returns}}
		{{joinName .Returns ", "}} := s.{{.Name}}(
			{{- if $hasCtxParam}}
			ctx,
			{{- end}}
			{{- range $params}}
			req.{{title .Name}} {{- if .Variadic}}...{{end}},
			{{- end}}
		)
		return {{addAmpersand .Name}}Response{
			{{- range .Returns}}
			{{title .Name}}: {{.Name}},
			{{- end}}
		}, nil
		{{- else}}
		s.{{.Name}}(
			{{- if $hasCtxParam}}
			ctx,
			{{- end}}
			{{- range $params}}
			req.{{title .Name}} {{- if .Variadic}}...{{end}},
			{{- end}}
		)
		return nil, nil
		{{- end}} {{/* End of if .Returns */}}
	}
}

{{- end}} {{/* End of range .DocMethods */}}

func unmarshal(data []byte, wrapper interface{}) error {
	s := string(data)
	if s == "null" || s == "" {
		return nil
	}

	var x struct {
		Raw json.RawMessage
		Type string
	}
	t, err := codec.GetType(x.Type) 
	if err != nil {
		return err
	}

	value := reflect.ValueOf(wrapper).Elem()
	field := value.FieldByName("W")

	defer field.Set(reflect.ValueOf(t))

	if len(x.Raw) == 0 {
		return nil 
	}
	return json.Unmarshal(x.Raw, t)
}
`
)

type Options struct {
	SchemaPtr bool
	SchemaTag string
	Formatted bool
	SnakeCase bool
}

type Generator struct {
	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, ifaceData *ifacetool.Data, spec *openapi.Specification) (*generator.File, error) {
	operationMap := make(map[string]*openapi.Operation)
	for _, op := range spec.Operations {
		operationMap[op.GoMethodName] = op
	}

	type MethodWithOp struct {
		*ifacetool.Method
		Op *openapi.Operation
	}

	type ParamWithAlias struct {
		*ifacetool.Param
		Alias string
	}

	var docMethods []MethodWithOp
	for _, m := range ifaceData.Methods {
		if op, ok := operationMap[m.Name]; ok {
			docMethods = append(docMethods, MethodWithOp{
				Method: m,
				Op:     op,
			})
		}
	}

	data := struct {
		PkgInfo    *generator.PkgInfo
		Data       *ifacetool.Data
		DocMethods []MethodWithOp
	}{
		PkgInfo:    pkgInfo,
		Data:       ifaceData,
		DocMethods: docMethods,
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"title": strings.Title,
			"nonCtxParams": func(params []*ifacetool.Param, reqParams []*openapi.Param) (out []ParamWithAlias) {
				nameToAlias := make(map[string]string)
				for _, p := range reqParams {
					if p.In == openapi.InBody {
						// Only parameters in body are supported for changing tag-name by alias.
						nameToAlias[p.Name] = p.Alias
					}
				}

				for _, p := range params {
					if p.TypeString != "context.Context" {
						out = append(out, ParamWithAlias{
							Param: p,
							Alias: nameToAlias[p.Name],
						})
					}
				}
				return
			},
			"hasCtxParam": func(params []*ifacetool.Param) bool {
				for _, p := range params {
					if p.TypeString == "context.Context" {
						return true
					}
				}
				return false
			},
			"getErrParamName": func(params []*ifacetool.Param) string {
				for _, p := range params {
					if p.TypeString == "error" {
						return p.Name
					}
				}
				return ""
			},
			"joinName": func(returns []*ifacetool.Param, sep string) string {
				var names []string
				for _, r := range returns {
					names = append(names, r.Name)
				}
				return strings.Join(names, sep)
			},
			"addAsterisks": func(name string) string {
				if g.opts.SchemaPtr {
					return "*" + name
				}
				return name
			},
			"addAmpersand": func(name string) string {
				if g.opts.SchemaPtr {
					return "&" + name
				}
				return name
			},
			"addTag": func(name, typ string) string {
				if g.opts.SchemaTag == "" {
					return ""
				}

				if name == "" || typ == "error" {
					name = "-"
				} else {
					// Only useful for adding correct tags for Response fields.
					if g.opts.SnakeCase {
						name = caseconv.ToSnakeCase(name)
					} else {
						name = caseconv.ToLowerCamelCase(name)
					}
				}

				return fmt.Sprintf("`%s:\"%s\"`", g.opts.SchemaTag, name)
			},
			"interfaceWrapper": func(params []*ifacetool.Param) string {
				var results []string

				for _, p := range params {
					t := p.Type.Underlying()
					if _, ok :=  t.(*types.Interface); !ok || p.TypeString == "context.Context" {
						continue
					}
					
					s := strings.Split(p.TypeString, ".")
					name := s[len(s) -1]

					results = append(results, fmt.Sprintf("type w%s struct { W %s }", name, p.TypeString))
					results = append(results, fmt.Sprintf("func (w *w%s) UnmarshalJSON(raw []byte) error { return unmarshal(w, raw) }", name))
				}

				result := strings.Join(results, "\n")
				fmt.Printf("\nresult: %#v\n", result)

				return result
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "endpoint.go",
	})
}
