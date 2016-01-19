package modelgen_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/raphael/goa/design"
	"github.com/raphael/goa/goagen/codegen"
)

var _ = Describe("ContextsWriter", func() {
	var writer *modelgen.ContextsWriter
	var filename string
	var workspace *codegen.Workspace

	JustBeforeEach(func() {
		var err error
		workspace, err = codegen.NewWorkspace("test")
		Ω(err).ShouldNot(HaveOccurred())
		pkg, err := workspace.NewPackage("contexts")
		Ω(err).ShouldNot(HaveOccurred())
		src := pkg.CreateSourceFile("test.go")
		filename = src.Abs()
		writer, err = modelgen.NewContextsWriter(filename)
		Ω(err).ShouldNot(HaveOccurred())
		codegen.TempCount = 0
	})

	AfterEach(func() {
		workspace.Delete()
	})

	Context("correctly configured", func() {
		var f *os.File
		BeforeEach(func() {
			f, _ = ioutil.TempFile("", "")
			filename = f.Name()
		})

		AfterEach(func() {
			os.Remove(filename)
		})

		Context("with data", func() {
			var params, headers *design.AttributeDefinition
			var payload *design.UserTypeDefinition
			var responses map[string]*design.ResponseDefinition
			var mediaTypes map[string]*design.MediaTypeDefinition

			var data *modelgen.ContextTemplateData

			BeforeEach(func() {
				params = nil
				headers = nil
				payload = nil
				responses = nil
				mediaTypes = nil
				data = nil
			})

			JustBeforeEach(func() {
				var version *design.APIVersionDefinition
				if design.Design != nil {
					version = design.Design.APIVersionDefinition
				} else {
					version = &design.APIVersionDefinition{}
				}
				data = &modelgen.ContextTemplateData{
					Name:         "ListBottleContext",
					ResourceName: "bottles",
					ActionName:   "list",
					Params:       params,
					Payload:      payload,
					Headers:      headers,
					Responses:    responses,
					API:          design.Design,
					Version:      version,
					DefaultPkg:   "",
				}
			})

			Context("with simple data", func() {
				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(emptyContext))
					Ω(written).Should(ContainSubstring(emptyContextFactory))
				})
			})

			Context("with an integer param", func() {
				BeforeEach(func() {
					intParam := &design.AttributeDefinition{Type: design.Integer}
					dataType := design.Object{
						"param": intParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(intContext))
					Ω(written).Should(ContainSubstring(intContextFactory))
				})
			})

			Context("with a string param", func() {
				BeforeEach(func() {
					strParam := &design.AttributeDefinition{Type: design.String}
					dataType := design.Object{
						"param": strParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(strContext))
					Ω(written).Should(ContainSubstring(strContextFactory))
				})
			})

			Context("with a number param", func() {
				BeforeEach(func() {
					numParam := &design.AttributeDefinition{Type: design.Number}
					dataType := design.Object{
						"param": numParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(numContext))
					Ω(written).Should(ContainSubstring(numContextFactory))
				})
			})

			Context("with a boolean param", func() {
				BeforeEach(func() {
					boolParam := &design.AttributeDefinition{Type: design.Boolean}
					dataType := design.Object{
						"param": boolParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(boolContext))
					Ω(written).Should(ContainSubstring(boolContextFactory))
				})
			})

			Context("with an array param", func() {
				BeforeEach(func() {
					str := &design.AttributeDefinition{Type: design.String}
					arrayParam := &design.AttributeDefinition{
						Type: &design.Array{ElemType: str},
					}
					dataType := design.Object{
						"param": arrayParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(arrayContext))
					Ω(written).Should(ContainSubstring(arrayContextFactory))
				})
			})

			Context("with an integer array param", func() {
				BeforeEach(func() {
					i := &design.AttributeDefinition{Type: design.Integer}
					intArrayParam := &design.AttributeDefinition{
						Type: &design.Array{ElemType: i},
					}
					dataType := design.Object{
						"param": intArrayParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(intArrayContext))
					Ω(written).Should(ContainSubstring(intArrayContextFactory))
				})
			})

			Context("with an param using a reserved keyword as name", func() {
				BeforeEach(func() {
					intParam := &design.AttributeDefinition{Type: design.Integer}
					dataType := design.Object{
						"int": intParam,
					}
					params = &design.AttributeDefinition{
						Type: dataType,
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(resContext))
					Ω(written).Should(ContainSubstring(resContextFactory))
				})
			})

			Context("with a required param", func() {
				BeforeEach(func() {
					intParam := &design.AttributeDefinition{Type: design.Integer}
					dataType := design.Object{
						"int": intParam,
					}
					required := design.RequiredValidationDefinition{
						Names: []string{"int"},
					}
					params = &design.AttributeDefinition{
						Type:        dataType,
						Validations: []design.ValidationDefinition{&required},
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(requiredContext))
					Ω(written).Should(ContainSubstring(requiredContextFactory))
				})
			})

			Context("with a simple payload", func() {
				BeforeEach(func() {
					payload = &design.UserTypeDefinition{
						AttributeDefinition: &design.AttributeDefinition{Type: design.String},
						TypeName:            "ListBottlePayload",
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(payloadContext))
					Ω(written).Should(ContainSubstring(payloadContextFactory))
				})
			})

			Context("with a object payload", func() {
				BeforeEach(func() {
					intParam := &design.AttributeDefinition{Type: design.Integer}
					strParam := &design.AttributeDefinition{Type: design.String}
					dataType := design.Object{
						"int": intParam,
						"str": strParam,
					}
					required := design.RequiredValidationDefinition{
						Names: []string{"int"},
					}
					payload = &design.UserTypeDefinition{
						AttributeDefinition: &design.AttributeDefinition{
							Type:        dataType,
							Validations: []design.ValidationDefinition{&required},
						},
						TypeName: "ListBottlePayload",
					}
				})

				It("writes the contexts code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(payloadObjContext))
				})
			})

		})
	})
})

var _ = Describe("ControllersWriter", func() {
	var writer *modelgen.ControllersWriter
	var workspace *codegen.Workspace
	var filename string

	BeforeEach(func() {
		var err error
		workspace, err = codegen.NewWorkspace("test")
		Ω(err).ShouldNot(HaveOccurred())
		pkg, err := workspace.NewPackage("controllers")
		Ω(err).ShouldNot(HaveOccurred())
		src := pkg.CreateSourceFile("test.go")
		filename = src.Abs()
	})

	JustBeforeEach(func() {
		var err error
		writer, err = modelgen.NewControllersWriter(filename)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		workspace.Delete()
	})

	Context("correctly configured", func() {
		var f *os.File
		BeforeEach(func() {
			f, _ = os.Create(filename)
		})

		Context("with data", func() {
			var actions, verbs, paths, contexts, unmarshals []string
			var payloads []*design.UserTypeDefinition

			var data []*modelgen.ControllerTemplateData

			BeforeEach(func() {
				actions = nil
				verbs = nil
				paths = nil
				contexts = nil
				unmarshals = nil
				payloads = nil
			})

			JustBeforeEach(func() {
				d := &modelgen.ControllerTemplateData{
					Resource: "Bottles",
					Version:  &design.APIVersionDefinition{},
				}
				as := make([]map[string]interface{}, len(actions))
				for i, a := range actions {
					var unmarshal string
					var payload *design.UserTypeDefinition
					if i < len(unmarshals) {
						unmarshal = unmarshals[i]
					}
					if i < len(payloads) {
						payload = payloads[i]
					}
					as[i] = map[string]interface{}{
						"Name": a,
						"Routes": []*design.RouteDefinition{
							&design.RouteDefinition{
								Verb: verbs[i],
								Path: paths[i],
							}},
						"Context":   contexts[i],
						"Unmarshal": unmarshal,
						"Payload":   payload,
					}
				}
				if len(as) > 0 {
					d.Actions = as
					data = []*modelgen.ControllerTemplateData{d}
				} else {
					data = nil
				}
			})

			Context("with missing data", func() {
				It("returns an empty string", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).Should(BeEmpty())
				})
			})

			Context("with a simple controller", func() {
				BeforeEach(func() {
					actions = []string{"list"}
					verbs = []string{"GET"}
					paths = []string{"/accounts/:accountID/bottles"}
					contexts = []string{"ListBottleContext"}
				})

				It("writes the controller code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(simpleController))
					Ω(written).Should(ContainSubstring(simpleMount))
				})
			})

			Context("with actions that take a payload", func() {
				BeforeEach(func() {
					actions = []string{"list"}
					verbs = []string{"GET"}
					paths = []string{"/accounts/:accountID/bottles"}
					contexts = []string{"ListBottleContext"}
					unmarshals = []string{"unmarshalListBottlePayload"}
					payloads = []*design.UserTypeDefinition{
						&design.UserTypeDefinition{
							TypeName: "ListBottlePayload",
							AttributeDefinition: &design.AttributeDefinition{
								Type: design.Object{
									"id": &design.AttributeDefinition{
										Type: design.String,
									},
								},
							},
						},
					}
				})

				It("writes the payload unmarshal function", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).Should(ContainSubstring(payloadObjUnmarshal))
				})
			})

			Context("with multiple controllers", func() {
				BeforeEach(func() {
					actions = []string{"list", "show"}
					verbs = []string{"GET", "GET"}
					paths = []string{"/accounts/:accountID/bottles", "/accounts/:accountID/bottles/:id"}
					contexts = []string{"ListBottleContext", "ShowBottleContext"}
				})

				It("writes the controllers code", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
					b, err := ioutil.ReadFile(filename)
					Ω(err).ShouldNot(HaveOccurred())
					written := string(b)
					Ω(written).ShouldNot(BeEmpty())
					Ω(written).Should(ContainSubstring(multiController))
					Ω(written).Should(ContainSubstring(multiMount))
				})
			})
		})
	})
})

var _ = Describe("HrefWriter", func() {
	var writer *modelgen.ResourcesWriter
	var workspace *codegen.Workspace
	var filename string

	BeforeEach(func() {
		var err error
		workspace, err = codegen.NewWorkspace("test")
		Ω(err).ShouldNot(HaveOccurred())
		pkg, err := workspace.NewPackage("controllers")
		Ω(err).ShouldNot(HaveOccurred())
		src := pkg.CreateSourceFile("test.go")
		filename = src.Abs()
	})

	JustBeforeEach(func() {
		var err error
		writer, err = modelgen.NewResourcesWriter(filename)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		workspace.Delete()
	})

	Context("correctly configured", func() {
		Context("with data", func() {
			var canoTemplate string
			var canoParams []string
			var mediaType *design.MediaTypeDefinition

			var data *modelgen.ResourceData

			BeforeEach(func() {
				mediaType = nil
				canoTemplate = ""
				canoParams = nil
				data = nil
			})

			JustBeforeEach(func() {
				data = &modelgen.ResourceData{
					Name:              "Bottle",
					Identifier:        "vnd.acme.com/resources",
					Description:       "A bottle resource",
					Type:              mediaType,
					CanonicalTemplate: canoTemplate,
					CanonicalParams:   canoParams,
				}
			})

			Context("with missing resource type definition", func() {
				It("does not return an error", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
				})
			})

			Context("with a string resource", func() {
				BeforeEach(func() {
					attDef := &design.AttributeDefinition{
						Type: design.String,
					}
					mediaType = &design.MediaTypeDefinition{
						UserTypeDefinition: &design.UserTypeDefinition{
							AttributeDefinition: attDef,
							TypeName:            "Bottle",
						},
					}
				})
				It("does not return an error", func() {
					err := writer.Execute(data)
					Ω(err).ShouldNot(HaveOccurred())
				})
			})

			Context("with a user type resource", func() {
				BeforeEach(func() {
					intParam := &design.AttributeDefinition{Type: design.Integer}
					strParam := &design.AttributeDefinition{Type: design.String}
					dataType := design.Object{
						"int": intParam,
						"str": strParam,
					}
					attDef := &design.AttributeDefinition{
						Type: dataType,
					}
					mediaType = &design.MediaTypeDefinition{
						UserTypeDefinition: &design.UserTypeDefinition{
							AttributeDefinition: attDef,
							TypeName:            "Bottle",
						},
					}
				})

				Context("and a canonical action", func() {
					BeforeEach(func() {
						canoTemplate = "/bottles/%v"
						canoParams = []string{"id"}
					})

					It("writes the href method", func() {
						err := writer.Execute(data)
						Ω(err).ShouldNot(HaveOccurred())
						b, err := ioutil.ReadFile(filename)
						Ω(err).ShouldNot(HaveOccurred())
						written := string(b)
						Ω(written).ShouldNot(BeEmpty())
						Ω(written).Should(ContainSubstring(simpleResourceHref))
					})
				})
			})
		})
	})
})

const (
	emptyContext = `
type ListBottleContext struct {
	*goa.Context
}
`

	emptyContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	return &ctx, err
}
`

	intContext = `
type ListBottleContext struct {
	*goa.Context
	Param *int
}
`

	intContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawParam := c.Get("param")
	if rawParam != "" {
		if param, err2 := strconv.Atoi(rawParam); err2 == nil {
			tmp2 := int(param)
			tmp1 := &tmp2
			ctx.Param = tmp1
		} else {
			err = goa.InvalidParamTypeError("param", rawParam, "integer", err)
		}
	}
	return &ctx, err
}
`

	strContext = `
type ListBottleContext struct {
	*goa.Context
	Param *string
}
`

	strContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawParam := c.Get("param")
	if rawParam != "" {
		ctx.Param = &rawParam
	}
	return &ctx, err
}
`

	numContext = `
type ListBottleContext struct {
	*goa.Context
	Param *float64
}
`

	numContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawParam := c.Get("param")
	if rawParam != "" {
		if param, err2 := strconv.ParseFloat(rawParam, 64); err2 == nil {
			tmp1 := &param
			ctx.Param = tmp1
		} else {
			err = goa.InvalidParamTypeError("param", rawParam, "number", err)
		}
	}
	return &ctx, err
}
`
	boolContext = `
type ListBottleContext struct {
	*goa.Context
	Param *bool
}
`

	boolContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawParam := c.Get("param")
	if rawParam != "" {
		if param, err2 := strconv.ParseBool(rawParam); err2 == nil {
			tmp1 := &param
			ctx.Param = tmp1
		} else {
			err = goa.InvalidParamTypeError("param", rawParam, "boolean", err)
		}
	}
	return &ctx, err
}
`

	arrayContext = `
type ListBottleContext struct {
	*goa.Context
	Param []string
}
`

	arrayContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawParam := c.Get("param")
	if rawParam != "" {
		elemsParam := strings.Split(rawParam, ",")
		ctx.Param = elemsParam
	}
	return &ctx, err
}
`

	intArrayContext = `
type ListBottleContext struct {
	*goa.Context
	Param []int
}
`

	intArrayContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawParam := c.Get("param")
	if rawParam != "" {
		elemsParam := strings.Split(rawParam, ",")
		elemsParam2 := make([]int, len(elemsParam))
		for i, rawElem := range elemsParam {
			if elem, err2 := strconv.Atoi(rawElem); err2 == nil {
				elemsParam2[i] = int(elem)
			} else {
				err = goa.InvalidParamTypeError("elem", rawElem, "integer", err)
			}
		}
		ctx.Param = elemsParam2
	}
	return &ctx, err
}
`

	resContext = `
type ListBottleContext struct {
	*goa.Context
	Int *int
}
`

	resContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawInt := c.Get("int")
	if rawInt != "" {
		if int_, err2 := strconv.Atoi(rawInt); err2 == nil {
			tmp2 := int(int_)
			tmp1 := &tmp2
			ctx.Int = tmp1
		} else {
			err = goa.InvalidParamTypeError("int", rawInt, "integer", err)
		}
	}
	return &ctx, err
}
`

	requiredContext = `
type ListBottleContext struct {
	*goa.Context
	Int int
}
`

	requiredContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	rawInt := c.Get("int")
	if rawInt == "" {
		err = goa.MissingParamError("int", err)
	} else {
		if int_, err2 := strconv.Atoi(rawInt); err2 == nil {
			ctx.Int = int(int_)
		} else {
			err = goa.InvalidParamTypeError("int", rawInt, "integer", err)
		}
	}
	return &ctx, err
}
`

	payloadContext = `
type ListBottleContext struct {
	*goa.Context
	Payload ListBottlePayload
}
`

	payloadContextFactory = `
func NewListBottleContext(c *goa.Context) (*ListBottleContext, error) {
	var err error
	ctx := ListBottleContext{Context: c}
	return &ctx, err
}
`
	payloadObjContext = `
type ListBottleContext struct {
	*goa.Context
	Payload *ListBottlePayload
}
`

	payloadObjUnmarshal = `
func unmarshalListBottlePayload(ctx *goa.Context) error {
	payload := &ListBottlePayload{}
	if err := ctx.Service().DecodeRequest(ctx, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		return err
	}
	ctx.SetPayload(payload)
	return nil
}
`

	simpleController = `// BottlesController is the controller interface for the Bottles actions.
type BottlesController interface {
	goa.Controller
	list(*ListBottleContext) error
}
`

	simpleMount = `func MountBottlesController(service goa.Service, ctrl BottlesController) {
	var h goa.Handler
	mux := service.ServeMux()
	h = func(c *goa.Context) error {
		ctx, err := NewListBottleContext(c)
		if err != nil {
			return goa.NewBadRequestError(err)
		}
		return ctrl.list(ctx)
	}
	mux.Handle("GET", "/accounts/:accountID/bottles", ctrl.HandleFunc("list", h, nil))
	service.Info("mount", "ctrl", "Bottles", "action", "list", "route", "GET /accounts/:accountID/bottles")
}
`

	multiController = `// BottlesController is the controller interface for the Bottles actions.
type BottlesController interface {
	goa.Controller
	list(*ListBottleContext) error
	show(*ShowBottleContext) error
}
`

	multiMount = `func MountBottlesController(service goa.Service, ctrl BottlesController) {
	var h goa.Handler
	mux := service.ServeMux()
	h = func(c *goa.Context) error {
		ctx, err := NewListBottleContext(c)
		if err != nil {
			return goa.NewBadRequestError(err)
		}
		return ctrl.list(ctx)
	}
	mux.Handle("GET", "/accounts/:accountID/bottles", ctrl.HandleFunc("list", h, nil))
	service.Info("mount", "ctrl", "Bottles", "action", "list", "route", "GET /accounts/:accountID/bottles")
	h = func(c *goa.Context) error {
		ctx, err := NewShowBottleContext(c)
		if err != nil {
			return goa.NewBadRequestError(err)
		}
		return ctrl.show(ctx)
	}
	mux.Handle("GET", "/accounts/:accountID/bottles/:id", ctrl.HandleFunc("show", h, nil))
	service.Info("mount", "ctrl", "Bottles", "action", "show", "route", "GET /accounts/:accountID/bottles/:id")
}
`

	simpleResourceHref = `func BottleHref(id interface{}) string {
	return fmt.Sprintf("/bottles/%v", id)
}
`
)
