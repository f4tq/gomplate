package main

import (
	"io"
	"log"
	"text/template"

	"github.com/hairyhenderson/gomplate/data"
	"fmt"
	"os"
	"reflect"
)

func (g *Gomplate) createTemplate() *template.Template {

	if g.Template == nil {
		g.Template = template.New("template")
	}

	return g.Template.Funcs(g.funcMap).Option("missingkey=error")
}

// Gomplate -
type Gomplate struct {
	*template.Template
	funcMap    template.FuncMap
	leftDelim  string
	rightDelim string
	verbose    bool
	data       *data.Data
}
func (g *Gomplate) populateContext(datasource string, resolved interface{}, target map[string]interface{}) error {
	if g.verbose {
		fmt.Fprintf(os.Stderr, "make top-level %s val: %+v type: %+v kind: %+v\n", datasource, resolved, reflect.TypeOf(resolved), reflect.TypeOf(resolved).Kind())
	}
	// make it so if -d .=foo.yaml, then all values are moved to top-level
	var kind reflect.Kind
	var val reflect.Value
	if reflect.TypeOf(resolved) == reflect.TypeOf(reflect.Value{}) {
		val = resolved.(reflect.Value)
		if val.Type().Kind() == reflect.Interface{
			val = val.Elem()
		}
		ty := val.Type()
		kind = ty.Kind()
		if g.verbose {
			log.Printf( "value is reflect.Value.  val: %+v type: %+v kind: %+v\n", val, ty, kind)
		}
	} else {
		kind = reflect.TypeOf(resolved).Kind()
		val = reflect.ValueOf(resolved)
	}

	switch kind  {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			if key.Kind() == reflect.String {
				target[key.String()] = val.MapIndex(key).Interface()
				if g.verbose {
					log.Printf( "key %s\n", key)
				}

			} else if key.Type() == reflect.TypeOf(reflect.Value{}) {
				kyV := key.Addr()
				if g.verbose {
					log.Printf( "key is Value type: %+v kind: %+v\n", kyV.Type(), kyV.Kind())
				}
				target[kyV.String()] = val.MapIndex(key).Interface()

			} else if key.Kind()== reflect.Interface {
				kvY := key.Elem()
				if kvY.Kind() == reflect.String {
					target[kvY.String()] = val.MapIndex(key).Interface()
				}
			} else {

				log.Fatalf("unknown key type key %s value=%+v\n", key.Type(), val.MapIndex(key).Interface())
			}

		}
	case reflect.Struct:
		if g.verbose {
			log.Printf( "struct: %+v orig val: %v\n", val, resolved)
		}
		for i := 0; i < val.NumField(); i += 1 {
			if val.Field(i).CanInterface() {
				if g.verbose {
					log.Printf( "adding struct field key %s\n", val.Type().Field(i).Name, val.Field(i).Interface())
				}
				target[val.Type().Field(i).Name] = val.Field(i).Interface()
			} else {
				log.Fatalf("%s SKipping private field %s %s\n", datasource, val.Type().Field(i).Name, val.Type().String())
			}
		}
	default:
		target[datasource] = resolved
	}
	return nil

}

// RunTemplate -
func (g *Gomplate) RunTemplate(text string, out io.Writer) {
	tmpl, err := g.createTemplate().Delims(g.leftDelim, g.rightDelim).Parse(text)
	if err != nil {
		log.Fatalf("Line %q: %v\n", text, err)
	}
	ctxt := make(map[string]interface{}, 0)
	if g.data != nil && g.data.Sources != nil {
		for k, ss := range g.data.Sources {
			resolved := g.data.Datasource(k)
			//if strings.TrimSpace(k) == "." {
			if ss.Type == "application/json" || ss.Type == "application/yaml" {

				if k !="." {
					target  := make(map[string]interface{},0)
					g.populateContext(k,resolved,target)
					ctxt[k] = target
				} else{
					g.populateContext(k,resolved,ctxt)
				}

			} else {
				ctxt[k] = resolved
				if g.verbose {
					fmt.Fprintf(os.Stderr, "Injecting context name: .%s\n", k)
				}
			}
		}
	}

	if err := tmpl.Execute(out, ctxt); err != nil {
		log.Fatalf("template failed because %s",err)

	}
}

// NewGomplate -
func NewGomplate(d *data.Data, leftDelim, rightDelim string) *Gomplate {
	return &Gomplate{
		leftDelim:  leftDelim,
		rightDelim: rightDelim,
		funcMap:    initFuncs(d),
		data: d,
	}
}

func runTemplate(o *GomplateOpts) error {
	defer runCleanupHooks()
	d := data.NewData(o.dataSources, o.dataSourceHeaders)
	addCleanupHook(d.Cleanup)

	g := NewGomplate(d, o.lDelim, o.rDelim)
	if o.verbose {
		g.verbose = true
	}
	if o.partialDir != "" {
		if err := partialInputDir(o.partialDir, o.flattenPartial, g); err != nil {
			return err
		}
	}

	if o.inputDir != "" {
		return processInputDir(o.inputDir, o.outputDir, g)
	}

	return processInputFiles(o.input, o.inputFiles, o.outputFiles, g)
}

// Called from process.go ...
func renderTemplate(g *Gomplate, inString string, outPath string) error {
	outFile, err := openOutFile(outPath)
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer outFile.Close()
	g.RunTemplate(inString, outFile)
	return nil
}
