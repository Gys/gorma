package gorma

import (
	"regexp"
	"strings"
	"text/template"

	"github.com/raphael/goa/design"
	"github.com/raphael/goa/goagen/codegen"
)

// WildcardRegex is the regex used to capture path parameters.
var WildcardRegex = regexp.MustCompile("(?:[^/]*/:([^/]+))+")

type (
	// ContextsWriter generate codes for a goa application contexts.
	ContextsWriter struct {
		*codegen.GoGenerator
		CtxTmpl        *template.Template
		CtxNewTmpl     *template.Template
		CtxRespTmpl    *template.Template
		PayloadTmpl    *template.Template
		NewPayloadTmpl *template.Template
	}

	// ResourcesWriter generate code for a goa application resources.
	// Resources are data structures initialized by the application handlers and passed to controller
	// actions.
	ResourcesWriter struct {
		*codegen.GoGenerator
		ResourceTmpl *template.Template
	}

	// MediaTypesWriter generate code for a goa application media types.
	// Media types are data structures used to render the response bodies.
	MediaTypesWriter struct {
		*codegen.GoGenerator
		MediaTypeTmpl *template.Template
	}

	// UserTypesWriter generate code for a goa application user types.
	// User types are data structures defined in the DSL with "Type".
	UserTypesWriter struct {
		*codegen.GoGenerator
		UserTypeTmpl *template.Template
	}

	// ContextTemplateData contains all the information used by the template to render the context
	// code for an action.
	ContextTemplateData struct {
		Name         string // e.g. "ListBottleContext"
		ResourceName string // e.g. "bottles"
		ActionName   string // e.g. "list"
		Params       *design.AttributeDefinition
		Payload      *design.UserTypeDefinition
		Headers      *design.AttributeDefinition
		Routes       []*design.RouteDefinition
		Responses    map[string]*design.ResponseDefinition
		API          *design.APIDefinition
		Version      *design.APIVersionDefinition
		DefaultPkg   string
	}

	// MediaTypeTemplateData contains all the information used by the template to redner the
	// media types code.
	MediaTypeTemplateData struct {
		MediaType  *design.MediaTypeDefinition
		Versioned  bool
		DefaultPkg string
	}

	// UserTypeTemplateData contains all the information used by the template to redner the
	// media types code.
	UserTypeTemplateData struct {
		UserType    *design.UserTypeDefinition
		PrimaryKeys map[string]PrimaryKey
		BelongsTo   []BelongsTo
		Many2Many   []Many2Many
		Options     ModelOptions
		Versioned   bool
		DefaultPkg  string
	}

	ModelOptions struct {
		Cached           bool
		TableName        string
		DynamicTableName bool
		NoMedia          bool
		Roler            bool
		SQLTag           string
	}
	BelongsTo struct {
		Parent        string
		DatabaseField string
	}
	Many2Many struct {
		Relation            string
		LowerRelation       string
		PluralRelation      string
		LowerPluralRelation string
		TableName           string
	}
	// ControllerTemplateData contains the information required to generate an action handler.
	ControllerTemplateData struct {
		Resource string                       // Lower case plural resource name, e.g. "bottles"
		Actions  []map[string]interface{}     // Array of actions, each action has keys "Name", "Routes" and "Context"
		Version  *design.APIVersionDefinition // Controller API version
	}

	// ResourceData contains the information required to generate the resource GoGenerator
	ResourceData struct {
		Name              string                      // Name of resource
		Identifier        string                      // Identifier of resource media type
		Description       string                      // Description of resource
		Type              *design.MediaTypeDefinition // Type of resource media type
		CanonicalTemplate string                      // CanonicalFormat represents the resource canonical path in the form of a fmt.Sprintf format.
		CanonicalParams   []string                    // CanonicalParams is the list of parameter names that appear in the resource canonical path in order.
	}
)

// Versioned returns true if the context was built from an API version.
func (c *ContextTemplateData) Versioned() bool {
	return !c.Version.IsDefault()
}

// IsPathParam returns true if the given parameter name corresponds to a path parameter for all
// the context action routes. Such parameter is required but does not need to be validated as
// httprouter takes care of that.
func (c *ContextTemplateData) IsPathParam(param string) bool {
	params := c.Params
	pp := false
	if params.Type.IsObject() {
		for _, r := range c.Routes {
			pp = false
			for _, p := range r.Params(c.Version) {
				if p == param {
					pp = true
					break
				}
			}
			if !pp {
				break
			}
		}
	}
	return pp
}

// MustValidate returns true if code that checks for the presence of the given param must be
// generated.
func (c *ContextTemplateData) MustValidate(name string) bool {
	return c.Params.IsRequired(name) && !c.IsPathParam(name)
}

// NewContextsWriter returns a contexts code writer.
// Contexts provide the glue between the underlying request data and the user controller.
func NewContextsWriter(filename string) (*ContextsWriter, error) {
	cw := codegen.NewGoGenerator(filename)
	funcMap := cw.FuncMap
	funcMap["gotyperef"] = codegen.GoTypeRef
	funcMap["gotypedef"] = codegen.GoTypeDef
	funcMap["goify"] = codegen.Goify
	funcMap["gotypename"] = codegen.GoTypeName
	funcMap["gopkgtypename"] = codegen.GoPackageTypeName
	funcMap["typeUnmarshaler"] = codegen.TypeUnmarshaler
	funcMap["userTypeUnmarshalerImpl"] = codegen.UserTypeUnmarshalerImpl
	funcMap["validationChecker"] = codegen.ValidationChecker
	funcMap["tabs"] = codegen.Tabs
	funcMap["add"] = func(a, b int) int { return a + b }
	funcMap["gopkgtyperef"] = codegen.GoPackageTypeRef
	ctxTmpl, err := template.New("context").Funcs(funcMap).Parse(ctxT)
	if err != nil {
		return nil, err
	}
	ctxNewTmpl, err := template.New("new").
		Funcs(cw.FuncMap).
		Funcs(template.FuncMap{
		"newCoerceData":  newCoerceData,
		"arrayAttribute": arrayAttribute,
		"tempvar":        codegen.Tempvar,
	}).Parse(ctxNewT)
	if err != nil {
		return nil, err
	}
	ctxRespTmpl, err := template.New("response").Funcs(cw.FuncMap).Parse(ctxRespT)
	if err != nil {
		return nil, err
	}
	payloadTmpl, err := template.New("payload").Funcs(funcMap).Parse(payloadT)
	if err != nil {
		return nil, err
	}
	newPayloadTmpl, err := template.New("newpayload").Funcs(cw.FuncMap).Parse(newPayloadT)
	if err != nil {
		return nil, err
	}
	w := ContextsWriter{
		GoGenerator:    cw,
		CtxTmpl:        ctxTmpl,
		CtxNewTmpl:     ctxNewTmpl,
		CtxRespTmpl:    ctxRespTmpl,
		PayloadTmpl:    payloadTmpl,
		NewPayloadTmpl: newPayloadTmpl,
	}
	return &w, nil
}

// Execute writes the code for the context types to the writer.
func (w *ContextsWriter) Execute(data *ContextTemplateData) error {
	if err := w.CtxTmpl.Execute(w, data); err != nil {
		return err
	}
	if err := w.CtxNewTmpl.Execute(w, data); err != nil {
		return err
	}
	if data.Payload != nil {
		if err := w.PayloadTmpl.Execute(w, data); err != nil {
			return err
		}
		if err := w.NewPayloadTmpl.Execute(w, data); err != nil {
			return err
		}
	}
	if len(data.Responses) > 0 {
		if err := w.CtxRespTmpl.Execute(w, data); err != nil {
			return err
		}
	}
	return nil
}

// NewResourcesWriter returns a contexts code writer.
// Resources provide the glue between the underlying request data and the user controller.
func NewResourcesWriter(filename string) (*ResourcesWriter, error) {
	cw := codegen.NewGoGenerator(filename)
	funcMap := cw.FuncMap
	funcMap["join"] = strings.Join
	resourceTmpl, err := template.New("resource").Funcs(cw.FuncMap).Parse(resourceT)
	if err != nil {
		return nil, err
	}
	w := ResourcesWriter{
		GoGenerator:  cw,
		ResourceTmpl: resourceTmpl,
	}
	return &w, nil
}

// Execute writes the code for the context types to the writer.
func (w *ResourcesWriter) Execute(data *ResourceData) error {
	return w.ResourceTmpl.Execute(w, data)
}

// NewMediaTypesWriter returns a contexts code writer.
// Media types contain the data used to render response bodies.
func NewMediaTypesWriter(filename string) (*MediaTypesWriter, error) {
	cw := codegen.NewGoGenerator(filename)
	funcMap := cw.FuncMap
	funcMap["gotypedef"] = codegen.GoTypeDef
	funcMap["gotyperef"] = codegen.GoTypeRef
	funcMap["goify"] = codegen.Goify
	funcMap["gotypename"] = codegen.GoTypeName
	funcMap["gonative"] = codegen.GoNativeType
	funcMap["typeUnmarshaler"] = codegen.TypeUnmarshaler
	funcMap["typeMarshaler"] = codegen.MediaTypeMarshaler
	funcMap["recursiveValidate"] = codegen.RecursiveChecker
	funcMap["tempvar"] = codegen.Tempvar
	funcMap["newDumpData"] = newDumpData
	funcMap["userTypeUnmarshalerImpl"] = codegen.UserTypeUnmarshalerImpl
	funcMap["mediaTypeMarshalerImpl"] = codegen.MediaTypeMarshalerImpl
	mediaTypeTmpl, err := template.New("media type").Funcs(funcMap).Parse(mediaTypeT)
	if err != nil {
		return nil, err
	}
	w := MediaTypesWriter{
		GoGenerator:   cw,
		MediaTypeTmpl: mediaTypeTmpl,
	}
	return &w, nil
}

// Execute writes the code for the context types to the writer.
func (w *MediaTypesWriter) Execute(data *MediaTypeTemplateData) error {
	return w.MediaTypeTmpl.Execute(w, data)
}

// NewUserTypesWriter returns a contexts code writer.
// User types contain custom data structured defined in the DSL with "Type".
func NewUserTypesWriter(filename string) (*UserTypesWriter, error) {
	cw := codegen.NewGoGenerator(filename)
	funcMap := cw.FuncMap
	funcMap["gotypedef"] = codegen.GoTypeDef
	funcMap["gotyperef"] = codegen.GoTypeRef
	funcMap["goify"] = codegen.Goify
	funcMap["gotypename"] = codegen.GoTypeName
	funcMap["recursiveValidate"] = codegen.RecursiveChecker
	funcMap["userTypeUnmarshalerImpl"] = codegen.UserTypeUnmarshalerImpl
	funcMap["userTypeMarshalerImpl"] = codegen.UserTypeMarshalerImpl
	funcMap["pkattributes"] = pkAttributes
	funcMap["pkwhere"] = pkWhere
	funcMap["pkwherefields"] = pkWhereFields
	funcMap["pkupdatefields"] = pkUpdateFields
	funcMap["lower"] = lower
	funcMap["storagedef"] = storageDef
	userTypeTmpl, err := template.New("user type").Funcs(funcMap).Parse(userTypeT)
	if err != nil {
		return nil, err
	}
	w := UserTypesWriter{
		GoGenerator:  cw,
		UserTypeTmpl: userTypeTmpl,
	}
	return &w, nil
}

// Execute writes the code for the context types to the writer.
func (w *UserTypesWriter) Execute(data *UserTypeTemplateData) error {
	return w.UserTypeTmpl.Execute(w, data)
}

// newCoerceData is a helper function that creates a map that can be given to the "Coerce" template.
func newCoerceData(name string, att *design.AttributeDefinition, pointer bool, pkg string, depth int) map[string]interface{} {
	return map[string]interface{}{
		"Name":      name,
		"VarName":   codegen.Goify(name, false),
		"Pointer":   pointer,
		"Attribute": att,
		"Pkg":       pkg,
		"Depth":     depth,
	}
}

// newDumpData is a helper function that creates a map that can be given to the "Dump" template.
func newDumpData(mt *design.MediaTypeDefinition, versioned bool, defaultPkg, context, source, target, view string) map[string]interface{} {
	return map[string]interface{}{
		"MediaType":  mt,
		"Context":    context,
		"Source":     source,
		"Target":     target,
		"View":       view,
		"Versioned":  versioned,
		"DefaultPkg": defaultPkg,
	}
}

// arrayAttribute returns the array element attribute definition.
func arrayAttribute(a *design.AttributeDefinition) *design.AttributeDefinition {
	return a.Type.(*design.Array).ElemType
}

const (
	// ctxT generates the code for the context data type.
	// template input: *ContextTemplateData
	ctxT = `// {{.Name}} provides the {{.ResourceName}} {{.ActionName}} action context.
type {{.Name}} struct {
	*goa.Context
{{if .Params}}{{$ctx := .}}{{range $name, $att := .Params.Type.ToObject}}{{/*
*/}}	{{goify $name true}} {{if and $att.Type.IsPrimitive ($ctx.Params.IsPrimitivePointer $name)}}*{{end}}{{gotyperef .Type nil 0}}
{{end}}{{end}}{{if .Payload}}	Payload {{gotyperef .Payload nil 0}}
{{end}}}
`
	// coerceT generates the code that coerces the generic deserialized
	// data to the actual type.
	// template input: map[string]interface{} as returned by newCoerceData
	coerceT = `{{if eq .Attribute.Type.Kind 1}}{{/*

*/}}{{/* BooleanType */}}{{/*
*/}}{{$varName := or (and (not .Pointer) .VarName) tempvar}}{{/*
*/}}{{tabs .Depth}}if {{.VarName}}, err2 := strconv.ParseBool(raw{{goify .Name true}}); err2 == nil {
{{if .Pointer}}{{tabs .Depth}}	{{$varName}} := &{{.VarName}}
{{end}}{{tabs .Depth}}	{{.Pkg}} = {{$varName}}
{{tabs .Depth}}} else {
{{tabs .Depth}}	err = goa.InvalidParamTypeError("{{.Name}}", raw{{goify .Name true}}, "boolean", err)
{{tabs .Depth}}}
{{end}}{{if eq .Attribute.Type.Kind 2}}{{/*

*/}}{{/* IntegerType */}}{{/*
*/}}{{$tmp := tempvar}}{{/*
*/}}{{tabs .Depth}}if {{.VarName}}, err2 := strconv.Atoi(raw{{goify .Name true}}); err2 == nil {
{{if .Pointer}}{{$tmp2 := tempvar}}{{tabs .Depth}}	{{$tmp2}} := int({{.VarName}})
{{tabs .Depth}}	{{$tmp}} := &{{$tmp2}}
{{tabs .Depth}}	{{.Pkg}} = {{$tmp}}
{{else}}{{tabs .Depth}}	{{.Pkg}} = int({{.VarName}})
{{end}}{{tabs .Depth}}} else {
{{tabs .Depth}}	err = goa.InvalidParamTypeError("{{.Name}}", raw{{goify .Name true}}, "integer", err)
{{tabs .Depth}}}
{{end}}{{if eq .Attribute.Type.Kind 3}}{{/*

*/}}{{/* NumberType */}}{{/*
*/}}{{$varName := or (and (not .Pointer) .VarName) tempvar}}{{/*
*/}}{{tabs .Depth}}if {{.VarName}}, err2 := strconv.ParseFloat(raw{{goify .Name true}}, 64); err2 == nil {
{{if .Pointer}}{{tabs .Depth}}	{{$varName}} := &{{.VarName}}
{{end}}{{tabs .Depth}}	{{.Pkg}} = {{$varName}}
{{tabs .Depth}}} else {
{{tabs .Depth}}	err = goa.InvalidParamTypeError("{{.Name}}", raw{{goify .Name true}}, "number", err)
{{tabs .Depth}}}
{{end}}{{if eq .Attribute.Type.Kind 4}}{{/*

*/}}{{/* StringType */}}{{/*
*/}}{{tabs .Depth}}{{.Pkg}} = {{if .Pointer}}&{{end}}raw{{goify .Name true}}
{{end}}{{if eq .Attribute.Type.Kind 5}}{{/*

*/}}{{/* AnyType */}}{{/*
*/}}{{tabs .Depth}}{{.Pkg}} = {{if .Pointer}}&{{end}}raw{{goify .Name true}}
{{end}}{{if eq .Attribute.Type.Kind 6}}{{/*

*/}}{{/* ArrayType */}}{{/*
*/}}{{tabs .Depth}}elems{{goify .Name true}} := strings.Split(raw{{goify .Name true}}, ",")
{{if eq (arrayAttribute .Attribute).Type.Kind 4}}{{tabs .Depth}}{{.Pkg}} = elems{{goify .Name true}}
{{else}}{{tabs .Depth}}elems{{goify .Name true}}2 := make({{gotyperef .Attribute.Type nil .Depth}}, len(elems{{goify .Name true}}))
{{tabs .Depth}}for i, rawElem := range elems{{goify .Name true}} {
{{template "Coerce" (newCoerceData "elem" (arrayAttribute .Attribute) false (printf "elems%s2[i]" (goify .Name true)) (add .Depth 1))}}{{tabs .Depth}}}
{{tabs .Depth}}{{.Pkg}} = elems{{goify .Name true}}2
{{end}}{{end}}`

	// ctxNewT generates the code for the context factory method.
	// template input: *ContextTemplateData
	ctxNewT = `{{define "Coerce"}}` + coerceT + `{{end}}` + `
// New{{goify .Name true}} parses the incoming request URL and body, performs validations and creates the
// context used by the {{.ResourceName}} controller {{.ActionName}} action.
func New{{.Name}}(c *goa.Context) (*{{.Name}}, error) {
	var err error
	ctx := {{.Name}}{Context: c}
{{if .Headers}}{{$headers := .Headers}}{{range $name, $_ := $headers.Type.ToObject}}{{if ($headers.IsRequired $name)}}	if c.Request().Header.Get("{{$name}}") == "" {
		err = goa.MissingHeaderError("{{$name}}", err)
	}{{end}}{{end}}
{{end}}{{if.Params}}{{$ctx := .}}{{range $name, $att := .Params.Type.ToObject}}	raw{{goify $name true}} := c.Get("{{$name}}")
{{$mustValidate := $ctx.MustValidate $name}}{{if $mustValidate}}	if raw{{goify $name true}} == "" {
		err = goa.MissingParamError("{{$name}}", err)
	} else {
{{else}}	if raw{{goify $name true}} != "" {
{{end}}{{template "Coerce" (newCoerceData $name $att ($ctx.Params.IsPrimitivePointer $name) (printf "ctx.%s" (goify $name true)) 2)}}{{/*
*/}}{{$validation := validationChecker $att ($ctx.Params.IsNonZero $name) ($ctx.Params.IsRequired $name) (printf "ctx.%s" (goify $name true)) $name 2}}{{/*
*/}}{{if $validation}}{{$validation}}
{{end}}	}
{{end}}{{end}}{{/* if .Params */}}{{if .Payload}}	p, err := New{{gotypename .Payload nil 0}}(c.Payload())
	if err != nil {
		return nil, err
	}
	ctx.Payload = p
{{end}}	return &ctx, err
}

`
	// ctxRespT generates response helper methods GoGenerator
	// template input: *ContextTemplateData
	ctxRespT = `{{$ctx := .}}{{range .Responses}}{{$mt := $ctx.API.MediaTypeWithIdentifier .MediaType}}{{/*
*/}}// {{goify .Name true}} sends a HTTP response with status code {{.Status}}.
func (ctx *{{$ctx.Name}}) {{goify .Name true}}({{/*
*/}}{{if $mt}}resp {{gopkgtyperef $mt $mt.AllRequired $ctx.Versioned $ctx.DefaultPkg 0}}{{if gt (len $mt.ComputeViews) 1}}, view {{gopkgtypename $mt $mt.AllRequired $ctx.Versioned $ctx.DefaultPkg 0}}ViewEnum{{end}}{{/*
*/}}{{else if .MediaType}}resp []byte{{end}}) error {
{{if $mt}}	r, err := resp.Dump({{if gt (len $mt.ComputeViews) 1}}view{{end}})
	if err != nil {
		return fmt.Errorf("invalid response: %s", err)
	}
	ctx.Header().Set("Content-Type", "{{$mt.Identifier}}; charset=utf-8")
	return ctx.JSON({{.Status}}, r){{else}}	return ctx.Respond({{.Status}}, {{if and (not $mt) .MediaType}}resp{{else}}nil{{end}}){{end}}
}

{{end}}`

	// payloadT generates the payload type definition GoGenerator
	// template input: *ContextTemplateData
	payloadT = `{{$payload := .Payload}}// {{gotypename .Payload nil 0}} is the {{.ResourceName}} {{.ActionName}} action payload.
type {{gotypename .Payload nil 1}} {{gotypedef .Payload .Versioned .DefaultPkg 0 false}}
`
	// newPayloadT generates the code for the payload factory method.
	// template input: *ContextTemplateData
	newPayloadT = `{{$typeName := gotypename .Payload nil 0}}// New{{$typeName}} instantiates a {{$typeName}} from a raw request body.
// It validates each field and returns an error if any validation fails.
func New{{$typeName}}(raw interface{}) (p {{gotyperef .Payload nil 0}}, err error) {

{{typeUnmarshaler .Payload false "" "payload" "raw" "p"}}
	return
}{{if (not .Payload.IsPrimitive)}}

{{userTypeUnmarshalerImpl .Payload .Versioned .DefaultPkg "payload"}}{{end}}
`

	// ctrlT generates the controller interface for a given resource.
	// template input: *ControllerTemplateData
	ctrlT = `// {{.Resource}}Controller is the controller interface for the {{.Resource}} actions.
type {{.Resource}}Controller interface {
	goa.Controller
{{range .Actions}}	{{.Name}}(*{{.Context}}) error
{{end}}}
`

	// mountT generates the code for a resource "Mount" function.
	// template input: *ControllerTemplateData
	mountT = `
// Mount{{.Resource}}Controller "mounts" a {{.Resource}} resource controller on the given service.
func Mount{{.Resource}}Controller(service goa.Service, ctrl {{.Resource}}Controller) {
	var h goa.Handler
	mux := service.ServeMux(){{if not .Version.IsDefault}}.Version("{{.Version.Version}}"){{end}}
{{$res := .Resource}}{{$ver := .Version}}{{range .Actions}}{{$action := .}}	h = func(c *goa.Context) error {
		ctx, err := New{{.Context}}(c)
		if err != nil {
			return goa.NewBadRequestError(err)
		}
		return ctrl.{{.Name}}(ctx)
	}
{{range .Routes}}	mux.Handle("{{.Verb}}", "{{.FullPath $ver}}", ctrl.HandleFunc("{{$action.Name}}", h))
	service.Info("mount", "ctrl", "{{$res}}",{{if not $ver.IsDefault}} "version", "{{$ver.Version}}",{{end}} "action", "{{$action.Name}}", "route", "{{.Verb}} {{.FullPath $ver}}")
{{end}}{{end}}}
`

	// resourceT generates the code for a resource.
	// template input: *ResourceData
	resourceT = `{{if .CanonicalTemplate}}// {{.Name}}Href returns the resource href.
func {{.Name}}Href({{if .CanonicalParams}}{{join .CanonicalParams ", "}} interface{}{{end}}) string {
	return fmt.Sprintf("{{.CanonicalTemplate}}", {{join .CanonicalParams ", "}})
}
{{end}}`

	// mediaTypeT generates the code for a media type.
	// template input: MediaTypeTemplateData
	mediaTypeT = `{{define "Dump"}}` + dumpT + `{{end}}` + `// {{if .MediaType.Description}}{{.MediaType.Description}}{{else}}{{gotypename .MediaType .MediaType.AllRequired 0}} media type{{end}}
// Identifier: {{.MediaType.Identifier}}{{$typeName := gotypename .MediaType .MediaType.AllRequired 0}}
type {{$typeName}} {{gotypedef .MediaType .Versioned .DefaultPkg 0 false}}{{$computedViews := .MediaType.ComputeViews}}{{if gt (len $computedViews) 1}}

// {{$typeName}} views
type {{$typeName}}ViewEnum string


const (
{{range $name, $view := $computedViews}}// {{if .Description}}{{.Description}}{{else}}{{$typeName}} {{.Name}} view{{end}}
       {{$typeName}}{{goify .Name true}}View {{$typeName}}ViewEnum = "{{.Name}}"

{{end}}){{end}}
// Load{{$typeName}} loads raw data into an instance of {{$typeName}}

// into a variable of type interface{}. See https://golang.org/pkg/encoding/json/#Unmarshal for the
// complete list of supported data types.
func Load{{$typeName}}(raw interface{}) (res {{gotyperef .MediaType .MediaType.AllRequired 1}}, err error) {

	{{typeUnmarshaler .MediaType .Versioned .DefaultPkg "" "raw" "res"}}
	return
}

// Dump produces raw data from an instance of {{$typeName}} running all the
// validations. See Load{{$typeName}} for the definition of raw data.

func (mt {{gotyperef .MediaType .MediaType.AllRequired 0}}) Dump({{if gt (len $computedViews) 1}}view {{$typeName}}ViewEnum{{end}}) (res {{gonative .MediaType}}, err error) {
{{$mt := .MediaType}}{{$ctx := .}}{{if gt (len $computedViews) 1}}{{range $computedViews}}	if view == {{gotypename $mt $mt.AllRequired 0}}{{goify .Name true}}View {
		{{template "Dump" (newDumpData $mt $ctx.Versioned $ctx.DefaultPkg (printf "%s view" .Name) "mt" "res" .Name)}}
	}
{{end}}{{else}}{{range $mt.ComputeViews}}{{template "Dump" (newDumpData $mt $ctx.Versioned $ctx.DefaultPkg (printf "%s view" .Name) "mt" "res" .Name)}}{{/* ranges over the one element */}}
{{end}}{{end}}	return
}

{{$validation := recursiveValidate .MediaType.AttributeDefinition false false "mt" "response" 1}}{{if $validation}}// Validate validates the media type instance.
func (mt {{gotyperef .MediaType .MediaType.AllRequired 0}}) Validate() (err error) {
{{$validation}}
	return
}
{{end}}{{range $computedViews}}
{{mediaTypeMarshalerImpl $mt $ctx.Versioned $ctx.DefaultPkg .Name}}
{{end}}
{{userTypeUnmarshalerImpl .MediaType.UserTypeDefinition .Versioned .DefaultPkg "load"}}
`

	// dumpT generates the code for dumping a media type or media type collection element.
	dumpT = `{{if .MediaType.IsArray}}	{{.Target}} = make({{gonative .MediaType}}, len({{.Source}}))
{{$tmp := tempvar}}	for i, {{$tmp}} := range {{.Source}} {
{{$tmpel := tempvar}}		var {{$tmpel}} {{gonative .MediaType.ToArray.ElemType.Type}}
		{{template "Dump" (newDumpData .MediaType.ToArray.ElemType.Type .Versioned .DefaultPkg (printf "%s[*]" .Context) $tmp $tmpel .View)}}
		{{.Target}}[i] = {{$tmpel}}
	}{{else}}{{typeMarshaler .MediaType .Versioned .DefaultPkg .Context .Source .Target .View}}{{end}}`

	// userTypeT generates the code for a user type.
	// template input: UserTypeTemplateData
	userTypeT = `// {{if .UserType.Description}}{{.UserType.Description}}{{else}}{{gotypename .UserType .UserType.AllRequired 0}} type{{end}}
type {{gotypename .UserType .UserType.AllRequired 0}} {{gotypedef .UserType .Versioned .DefaultPkg 0 false}}
{{ if ne .Options.TableName "" }}
func (m {{.UserType.TypeName}}) TableName() string {
	return "{{ .Options.TableName}}"
}{{end}}
type {{.UserType.TypeName}}DB struct {
	Db gorm.DB
	{{ if .Options.Cached }}cache *cache.Cache{{end}}
}
func New{{.UserType.TypeName}}DB(db gorm.DB) *{{.UserType.TypeName}}DB {
	{{ if .Options.Cached }}return &{{.UserType.TypeName}}DB{
		Db: db,
		cache: cache.New(5*time.Minute, 30*time.Second),
	}
	{{ else  }}return &{{.UserType.TypeName}}DB{Db: db}{{ end  }}
}

func (m *{{.UserType.TypeName}}DB) DB() interface{} {
	return &m.Db
}
{{ if .Options.Roler }}
func (m {{.UserType.TypeName}}) GetRole() string {
	return *m.Role
}
{{end}}

{{$pks := .PrimaryKeys }}type {{.UserType.TypeName}}Storage interface {
	DB() interface{}
	List(ctx context.Context{{ if .Options.DynamicTableName}}, tableName string{{ end }}) []{{.UserType.TypeName}}
	One(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, {{ pkattributes $pks  }}) ({{.UserType.TypeName}}, error)
	Add(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, o {{.UserType.TypeName}}) ({{.UserType.TypeName}}, error)
	Update(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, o {{.UserType.TypeName}}) (error)
	Delete(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, {{ pkattributes $pks }}) (error)
	{{$typename:= .UserType.TypeName}}{{$options:=.Options}}{{ range $idx, $bt := .BelongsTo}}ListBy{{$bt.Parent}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, parentid int) []{{$typename}}
	OneBy{{$bt.Parent}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, parentid, id int) ({{$typename}}, error)
	{{end}}{{storagedef .UserType}}
}
func (m *{{$typename}}DB) One(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, {{pkattributes $pks}}) ({{$typename}}, error) {
	{{ if .Options.Cached }}//first attempt to retrieve from cache
	o,found := m.cache.Get(strconv.Itoa(id))
	if found {
		return o.({{$typename}}), nil
	}
	// fallback to database if not found{{ end }}
	var obj {{$typename}}
	{{ $l := len $pks }}
	{{ if eq $l 1 }}
	err := m.Db{{ if .Options.DynamicTableName }}.Table(tableName){{ end }}.Find(&obj, id).Error
	{{ else  }}
	err := m.Db{{ if .Options.DynamicTableName }}.Table(tableName){{ end }}.Find(&obj).Where("{{pkwhere $pks}}", {{pkwherefields $pks}}).Error
	{{ end }}
	{{ if .Options.Cached }} go m.cache.Set(strconv.Itoa(id), obj, cache.DefaultExpiration) {{ end }}
	return obj, err
}
func (m *{{$typename}}DB) Add(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, model {{$typename}}) ({{$typename}}, error) {
	err := m.Db{{ if .Options.DynamicTableName }}.Table(tableName){{ end }}.Create(&model).Error
	{{ if .Options.Cached }} go m.cache.Set(strconv.Itoa(model.ID), model, cache.DefaultExpiration) {{ end }}
	return model, err
}
func (m *{{$typename}}DB) Update(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, model {{$typename}}) error {
	obj, err := m.One(ctx{{ if .Options.DynamicTableName }}, tableName{{ end }}, {{pkupdatefields $pks}})
	if err != nil {
		return  err
	}
	err = m.Db{{ if .Options.DynamicTableName }}.Table(tableName){{ end }}.Model(&obj).Updates(model).Error
	{{ if .Options.Cached }}
	go func(){
	obj, err := m.One(ctx, model.ID)
	if err == nil {
		m.cache.Set(strconv.Itoa(model.ID), obj, cache.DefaultExpiration)
	}
	}()
	{{ end }}
	return err
}
func (m *{{$typename}}DB) Delete(ctx context.Context{{ if .Options.DynamicTableName }}, tableName string{{ end }}, {{pkattributes $pks}})  error {
	var obj {{$typename}}
	{{ $l := len $pks }}
	{{ if eq $l 1 }}
	err := m.Db{{ if .Options.DynamicTableName }}.Table(tableName){{ end }}.Delete(&obj, id).Error
	{{ else  }}
	err := m.Db{{ if .Options.DynamicTableName }}.Table(tableName){{ end }}.Delete(&obj).Where("{{pkwhere $pks}}", {{pkwherefields $pks}}).Error
	{{ end }}
	if err != nil {
		return  err
	}
	{{ if .Options.Cached }} go m.cache.Delete(strconv.Itoa(id)) {{ end }}
	return  nil
}
{{$options := .Options}}{{$typename := .UserType.TypeName}}{{ range $idx, $bt := .BelongsTo}}
// Belongs To Relationships
func {{$typename}}FilterBy{{$bt.Parent}}(parentid int, originaldb *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if parentid > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("{{ $bt.DatabaseField }} = ?", parentid)
		}
	} else {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
}
func (m *{{$typename}}DB) ListBy{{$bt.Parent}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, parentid int) []{{$typename}} {
	var objs []{{$typename}}
	m.Db{{ if $options.DynamicTableName }}.Table(tableName){{ end }}.Scopes({{$typename}}FilterBy{{$bt.Parent}}(parentid, &m.Db)).Find(&objs)
	return objs
}
func (m *{{$typename}}DB) OneBy{{$bt.Parent}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, parentid, {{ pkattributes $pks }}) ({{$typename}}, error) {
	{{ if $options.Cached }}//first attempt to retrieve from cache
	o,found := m.cache.Get(strconv.Itoa(id))
	if found {
		return o.({{$typename}}), nil
	}
	// fallback to database if not found{{ end }}
	var obj {{$typename}}
	err := m.Db{{ if $options.DynamicTableName }}.Table(tableName){{ end }}.Scopes({{$typename}}FilterBy{{$bt.Parent}}(parentid, &m.Db)).Find(&obj, id).Error
	{{ if $options.Cached }} go m.cache.Set(strconv.Itoa(id), obj, cache.DefaultExpiration) {{ end }}
	return obj, err
}
{{end}}

{{$options := .Options}}{{$typeName := .UserType.TypeName}}{{ range $idx, $bt := .Many2Many}}
// Many To Many Relationships
func (m *{{$typeName}}DB) Delete{{$bt.Relation}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, {{lower $typeName}}ID,  {{$bt.LowerRelation}}ID int)  error {
	var obj {{$typeName}}
	obj.ID = {{lower $typeName}}ID
	var assoc {{$bt.LowerRelation}}.{{$bt.Relation}}
	var err error
	assoc.ID = {{$bt.LowerRelation}}ID
	if err != nil {
		return err
	}
	err = m.Db{{ if $options.DynamicTableName }}.Table(tableName){{ end }}.Model(&obj).Association("{{$bt.PluralRelation}}").Delete(assoc).Error
	if err != nil {
		return  err
	}
	return  nil
}
func (m *{{$typeName}}DB) Add{{$bt.Relation}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, {{lower $typeName}}ID, {{$bt.LowerRelation}}ID int) error {
	var {{lower $typeName}} {{$typeName}}
	{{lower $typeName}}.ID = {{lower $typeName}}ID
	var assoc {{$bt.LowerRelation}}.{{$bt.Relation}}
	assoc.ID = {{$bt.LowerRelation}}ID
	err := m.Db{{ if $options.DynamicTableName }}.Table(tableName){{ end }}.Model(&{{lower $typeName}}).Association("{{$bt.PluralRelation}}").Append(assoc).Error
	if err != nil {
		return  err
	}
	return  nil
}
func (m *{{$typeName}}DB) List{{$bt.PluralRelation}}(ctx context.Context{{ if $options.DynamicTableName }}, tableName string{{ end }}, {{lower $typeName}}ID int)  []{{$bt.LowerRelation}}.{{$bt.Relation}} {
	var list []{{$bt.LowerRelation}}.{{$bt.Relation}}
	var obj {{$typeName}}
	obj.ID = {{lower $typeName}}ID
	m.Db{{ if $options.DynamicTableName }}.Table(tableName){{ end }}.Model(&obj).Association("{{$bt.PluralRelation}}").Find(&list)
	return  list
}
{{end}}
{{ range $idx, $bt := .BelongsTo}}
func Filter{{$typename}}By{{$bt.Parent}}(parent *int, list []{{$typename}}) []{{$typename}} {
	var filtered []{{$typename}}
	for _,o := range list {
		if o.{{$bt.Parent}}ID == int(*parent) {
			filtered = append(filtered,o)
		}
	}
	return filtered
}
{{end}}

`
)

/*


 */
