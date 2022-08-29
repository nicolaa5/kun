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
	"fmt"
	"reflect"
	"github.com/RussellLuo/validating/v2"
	"github.com/go-kit/kit/endpoint"
	"github.com/RussellLuo/kun/pkg/httpoption"
	"github.com/RussellLuo/kun/pkg/httpcodec"

	{{- range .Data.Imports}}
	{{.ImportString}}
	{{- end }}
)

var GetType = func(name string) (interface{}, error) {
	return nil, fmt.Errorf("implement GetType function for marshaling custom types")
}

{{- range .DocMethods}}

{{- $params := nonCtxParams .Params .Op.Request.Params}}
{{- $hasCtxParam := hasCtxParam .Params}}

{{ interfaceWrapper .Params .Returns}}

{{- if $params}}
type {{.Name}}Request struct {
	{{- range $params}}
	{{title .Name}} {{getType .Param}} {{addTag .Alias .TypeString}}
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
	{{title .Name}} {{getType .}} {{addTag .Name .TypeString}}
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
			{{joinRequestParam .Param "req"}} {{- if .Variadic}}...{{end}},
			{{- end}}
		)
		return {{addAmpersand .Name}}Response{
			{{- range .Returns}}
			{{title .Name}}: {{joinResponseParam .}},
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

func marshal(x interface{}) ([]byte, error) {
	if x == nil {
		return []byte("null"), nil 
	}

	//store the type information in the wrapper struct
	var w struct {
		Raw json.RawMessage
		Type string
	}
	var err error 

	w.Raw, err = json.Marshal(x)
	if err != nil {
		return nil, err
	}

	typeOf := reflect.TypeOf(x)
	sub := strings.Split(typeOf.String(), ".")
	w.Type = sub[len(sub)-1]

	return json.Marshal(w)
}

func unmarshal(data []byte, wrapper interface{}) error {
	s := string(data)
	if s == "null" {
		return nil
	}

	var x struct {
		Raw json.RawMessage
		Type string
	}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	t, err := GetType(x.Type) 
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
			"getType": func(param *ifacetool.Param) string {
				if param.TypeString == "context.Context" ||
					param.TypeString == "error" {
					return param.TypeString
				}

				switch v := param.Type.Underlying().(type) {
				case *types.Interface:
					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("w%s", name)

				case *types.Slice:
					if _, ok := v.Elem().Underlying().(*types.Interface); !ok {
						return param.TypeString
					}

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("[]w%s", name)

				case *types.Map:
					if _, ok := v.Elem().Underlying().(*types.Interface); !ok {
						return param.TypeString
					}

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("map[%s]w%s", v.Key().String(), name)
				default:
					return param.TypeString
				}
			},
			"joinRequestParam": func(param *ifacetool.Param, parent string) string {
				basic := fmt.Sprintf("%s.%s", parent, strings.Title(param.Name))

				switch v := param.Type.Underlying().(type) {
				case *types.Interface:
					return fmt.Sprintf("%s.%s.W", parent, strings.Title(param.Name))
				case *types.Slice:
					if _, ok := v.Elem().Underlying().(*types.Interface); !ok {
						return basic
					}

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("%sList(%s.%s)", name, parent, strings.Title(param.Name))
				case *types.Map:
					if _, ok := v.Elem().Underlying().(*types.Interface); !ok {
						return basic
					}

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("%sMap(%s.%s)", name, parent, strings.Title(param.Name))
				default:
					return basic
				}

			},
			"joinResponseParam": func(param *ifacetool.Param) string {
				if param.TypeString == "error" {
					return param.Name
				}

				switch v := param.Type.Underlying().(type) {
				case *types.Interface:

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("w%s{%s}", name, param.Name)
				case *types.Slice:
					if _, ok := v.Elem().Underlying().(*types.Interface); !ok {
						return param.Name
					}

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("w%sList(%s)", name, param.Name)
				case *types.Map:
					if _, ok := v.Elem().Underlying().(*types.Interface); !ok {
						return param.Name
					}

					s := strings.Split(param.TypeString, ".")
					name := s[len(s)-1]

					return fmt.Sprintf("w%sMap(%s)", name, param.Name)
				default:
					return param.Name
				}
			},
			"interfaceWrapper": func(params []*ifacetool.Param, returns []*ifacetool.Param) string {
				var results []string
				typesMap := make(map[string]*ifacetool.Param)

				for _, p := range params {
					if _, ok := typesMap[p.Name]; !ok {
						typesMap[p.Name] = p
					}
				}
				for _, r := range returns {
					if _, ok := typesMap[r.Name]; !ok {
						typesMap[r.Name] = r
					}
				}

				for _, v := range typesMap {
					t := v.Type.Underlying()
					if _, ok := t.(*types.Interface); !ok || 
						v.TypeString == "context.Context" || 
						v.TypeString == "error" {
						continue
					}

					s := strings.Split(v.TypeString, ".")
					name := s[len(s)-1]

					results = append(results, fmt.Sprintf("type w%s struct { W %s }", name, v.TypeString))
					results = append(results, fmt.Sprintf("func (w w%s) MarshalJSON() ([]byte, error) { return marshal(w.W); }", name))
					results = append(results, fmt.Sprintf("func (w *w%s) UnmarshalJSON(raw []byte) error { return unmarshal(raw, w); }", name))
					results = append(results, fmt.Sprintf("func %sList(list []w%s) (result []%s) { for _, item := range list { result = append(result, item.W) }; return result }", name, name, v.TypeString))
					results = append(results, fmt.Sprintf("func w%sList(list []%s) (result []w%s) { for _, item := range list { result = append(result, w%s{item}) }; return result }", name, v.TypeString, name, name))

					results = append(results, fmt.Sprintf("func %sMap(m map[string]w%s) (result map[string]%s) { mm := make(map[string]%s); for k, item := range m { mm[k] = item.W }; return mm }", name, name, v.TypeString, v.TypeString))
					results = append(results, fmt.Sprintf("func w%sMap(m map[string]%s) (result map[string]w%s) { mm := make(map[string]w%s); for k, item := range m { mm[k] = w%s{item} }; return mm }", name, v.TypeString, name, name, name))
				}

				result := strings.Join(results, "\n\n")
				return result
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "endpoint.go",
	})
}
