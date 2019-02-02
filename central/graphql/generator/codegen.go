package generator

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"
)

var timestampType = reflect.TypeOf((*types.Timestamp)(nil))
var templates = map[string]string{
	"codegen": `// Code generated by graphql-generator. DO NOT EDIT.

package resolvers

import (
	"reflect"

	"github.com/graph-gophers/graphql-go"
	"github.com/stackrox/rox/central/graphql/generator"
{{- range $i, $_ := .Imports }}{{if $i}}
	"{{$i}}"
{{- end}}{{end}} // end range imports
)

func registerGeneratedTypes(builder generator.SchemaBuilder) {
	builder.AddType("Label", []string{"key: String!", "value: String!"})
{{- range $td := .Entries }}
{{- if isEnum $td.Data.Type }}
	generator.RegisterProtoEnum(builder, reflect.TypeOf({{ importedName $td.Data.Type }}(0)))
{{- else }}
	builder.AddType("{{ $td.Data.Name }}", []string{
{{- range $td.Data.FieldData }}{{ if schemaType . }}
		"{{ lower .Name}}: {{ schemaType .}}",
{{- end }}{{ end }}
{{- range $td.Data.UnionData }}
		"{{lower .Name }}: {{ $td.Data.Name }}{{.Name}}",
{{- end }}
	})
{{- range $ud := $td.Data.UnionData }}
	builder.AddUnionType("{{ $td.Data.Name }}{{ $ud.Name }}", []string{
{{- range $ud.Entries }}
		"{{ schemaType . }}",
{{- end }}
	})
{{- end }}
{{- end }}
{{- end }}
}
{{range $td := .Entries}}{{if isEnum $td.Data.Type}}
func to{{.Data.Name}}(value *string) {{importedName .Data.Type}} {
	if value != nil {
		return {{importedName .Data.Type}}({{importedName .Data.Type}}_value[*value])
	}
	return {{importedName .Data.Type}}(0)
}

func to{{plural .Data.Name}}(values *[]string) []{{importedName .Data.Type}} {
	if values == nil {
		return nil
	}
	output := make([]{{importedName .Data.Type}}, len(*values))
	for i, v := range *values {
		output[i] = to{{.Data.Name}}(&v)
	}
	return output
}
{{else}}
type {{lower .Data.Name}}Resolver struct {
	root *Resolver
	data *{{importedName .Data.Type}}
{{- if .ListData}}
	list *{{listName .Data}}
{{- end}}
}

func (resolver *Resolver) wrap{{.Data.Name}}(value *{{importedName .Data.Type}}, ok bool, err error) (*{{lower .Data.Name}}Resolver, error) {
	if !ok || err != nil || value == nil {
		return nil, err
	}
	return &{{lower .Data.Name}}Resolver{resolver, value{{if .ListData}}, nil{{end}}}, nil
}

func (resolver *Resolver) wrap{{plural .Data.Name}}(values []*{{importedName .Data.Type}}, err error) ([]*{{lower .Data.Name}}Resolver, error) {
	if err != nil || len(values) == 0 {
		return nil, err
	}
	output := make([]*{{lower .Data.Name}}Resolver, len(values))
	for i, v := range values {
		output[i] = &{{lower .Data.Name}}Resolver{resolver, v{{if .ListData}}, nil{{end}}}
	}
	return output, nil
}
{{if .ListData}}
func (resolver *Resolver) wrapList{{plural .Data.Name}}(values []*{{listName .Data}}, err error) ([]*{{lower .Data.Name}}Resolver, error) {
	if err != nil || values == nil {
		return nil, err
	}
	output := make([]*{{lower .Data.Name}}Resolver, len(values))
	for i, v := range values {
		output[i] = &{{lower .Data.Name}}Resolver{resolver, nil, v}
	}
	return output, nil
}

func (resolver *{{lower .Data.Name}}Resolver) ensureData() {
	if resolver.data == nil {
		resolver.data = resolver.root.get{{.Data.Name}}(resolver.list.GetId())
	}
}
{{end -}}
{{range $_, $fd := .Data.FieldData -}}
{{$vt := valueType $fd}}{{if $vt}}
func (resolver *{{lower $td.Data.Name}}Resolver) {{ $fd.Name }}() {{ $vt }} {
{{- if $td.ListData}}{{if nonListField $td $fd}}
	resolver.ensureData()
{{- end}}{{end}}
	value := resolver.data.Get{{$fd.Name}}()
{{- if listField $td $fd}}
	if resolver.data == nil {
		value = resolver.list.Get{{$fd.Name}}()
	}
{{- end}}
	return {{translator $fd }}
}
{{end -}}
{{end}}
{{- range $ud := $td.Data.UnionData}}
type {{lower $td.Data.Name}}{{$ud.Name}}Resolver struct {
	resolver *{{lower $td.Data.Name}}Resolver
}

func (resolver *{{lower $td.Data.Name}}Resolver) {{$ud.Name}}() *{{lower $td.Data.Name}}{{$ud.Name}}Resolver {
	return &{{lower $td.Data.Name}}{{$ud.Name}}Resolver{resolver}
}
{{range $ut := $ud.Entries}}
func (resolver *{{lower $td.Data.Name}}{{$ud.Name}}Resolver) To{{$ut.Type.Elem.Name}}() (*{{lower $ut.Type.Elem.Name}}Resolver, bool) {
	value := resolver.resolver.data.Get{{$ut.Name}}()
	if value != nil {
		return &{{lower $ut.Type.Elem.Name}}Resolver{resolver.resolver.root, value}, true
	}
	return nil, false
}
{{end}}{{end}}{{end}}{{end}}`,
	"enum":         `value.String()`,
	"enumslice":    `stringSlice(value)`,
	"float":        `float64(value)`,
	"id":           `graphql.ID(value)`,
	"int":          `int32(value)`,
	"label":        `labelsResolver(value)`,
	"pointer":      `resolver.root.wrap{{.Type.Elem.Name}}(value, true, nil)`,
	"pointerslice": `resolver.root.wrap{{plural .Type.Elem.Elem.Name}}(value, nil)`,
	"raw":          `value`,
	"rawslice":     `value`,
	"time":         `timestamp(value)`,
}

func importedName(p reflect.Type) string {
	split := strings.Split(p.PkgPath(), "/")
	return fmt.Sprintf("%s.%s", split[len(split)-1], p.Name())
}

func listName(td typeData) string {
	split := strings.Split(td.Package, "/")
	return fmt.Sprintf("%s.List%s", split[len(split)-1], td.Name)
}

func isEnum(p reflect.Type) bool {
	if p == nil {
		return false
	}
	return proto.EnumValueMap(importedName(p)) != nil
}

func getFieldTransform(fd fieldData) (templateName string, returnType string) {
	switch fd.Type.Kind() {
	case reflect.String:
		if fd.Name == "Id" {
			return "id", "graphql.ID"
		}
		return "raw", "string"
	case reflect.Int32:
		if isEnum(fd.Type) {
			return "enum", "string"
		}
		return "raw", "int32"
	case reflect.Uint32:
		return "int", "int32"
	case reflect.Int64:
		return "int", "int32"
	case reflect.Float32:
		return "float", "float64"
	case reflect.Float64:
		return "raw", "float64"
	case reflect.Bool:
		return "raw", "bool"
	case reflect.Map:
		if fd.Type.Elem().Kind() == reflect.String && fd.Type.Elem().Kind() == reflect.String {
			return "label", "labels"
		}
	case reflect.Ptr:
		if fd.Type == timestampType {
			return "time", "(*graphql.Time, error)"
		}
		if fd.Type.Implements(messageType) {
			if isListType(fd.Type) {
				// if a field returns a list type, we don't automatically handle this for now.
				return "", ""
			}
			return "pointer", fmt.Sprintf("(*%sResolver, error)", lower(fd.Type.Elem().Name()))
		}
	case reflect.Slice:
		template, ret := getFieldTransform(fieldData{Name: fd.Name, Type: fd.Type.Elem()})
		if len(ret) > 0 && ret[0] == '(' {
			// this converts (*fooResolver, error) into ([]*fooResolver, error)
			return template + "slice", ret[0:1] + "[]" + ret[1:]
		}
		return template + "slice", "[]" + ret
	}
	return "", ""
}

func translator(t *template.Template) func(fieldData) string {
	return func(fd fieldData) string {
		tmplName, _ := getFieldTransform(fd)
		b := &bytes.Buffer{}
		err := t.ExecuteTemplate(b, tmplName, fd)
		if err != nil {
			panic(err)
		}
		return b.String()
	}
}

func valueType(fd fieldData) string {
	_, returnType := getFieldTransform(fd)
	return returnType
}

func listField(td schemaEntry, field fieldData) bool {
	t, ok := td.ListData[field.Name]
	return ok && t == field.Type
}

// GenerateResolvers produces go code for resolvers for all the types found by the typewalk.
func GenerateResolvers(parameters TypeWalkParameters, writer io.Writer) {
	data := typeWalk(
		parameters.IncludedTypes,
		[]reflect.Type{
			reflect.TypeOf((*types.Timestamp)(nil)),
		},
	)
	rootTemplate := template.New("codegen")
	rootTemplate.Funcs(template.FuncMap{
		"importedName": importedName,
		"isEnum":       isEnum,
		"listField":    listField,
		"listName":     listName,
		"lower":        lower,
		"nonListField": func(td schemaEntry, field fieldData) bool { return !listField(td, field) },
		"plural":       plural,
		"translator":   translator(rootTemplate),
		"valueType":    valueType,
		"schemaType":   schemaType,
	})
	for name, text := range templates {
		thisTemplate := rootTemplate
		if rootTemplate.Name() != name {
			thisTemplate = rootTemplate.New(name)
		}
		_, err := thisTemplate.Parse(text)
		if err != nil {
			panic(fmt.Sprintf("Template %q: %s", name, err))
		}
	}
	imports := make(map[string]bool)
	for _, td := range data {
		imports[td.Package] = true
	}
	entries := makeSchemaEntries(data, nil)
	err := rootTemplate.Execute(writer, struct {
		Entries []schemaEntry
		Imports map[string]bool
	}{entries, imports})
	if err != nil {
		panic(err)
	}
}
