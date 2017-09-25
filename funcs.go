package main

import (
	"text/template"
	html "html/template"
	"github.com/hairyhenderson/gomplate/data"
	"github.com/hairyhenderson/gomplate/funcs"
	"errors"

	"strings"
	"os"

)

// initFuncs - The function mappings are defined here!
func initFuncs(d *data.Data) template.FuncMap {
	f := template.FuncMap{}
	funcs.AddDataFuncs(f, d)
	funcs.AWSFuncs(f)
	funcs.AddBase64Funcs(f)
	funcs.AddNetFuncs(f)
	funcs.AddReFuncs(f)
	funcs.AddStringFuncs(f)
	funcs.AddEnvFuncs(f)
	funcs.AddConvFuncs(f)
	f["dict"] =dict
	f["slice"] = slice
	f["add"] = add
	f["minus"] = minus
	f["series"] = Series
	f["mult"] = mult
	f["div"] = div

	f["HtmlEscape"]= html.HTMLEscapeString
	f["URLEscape"]= html.URLQueryEscaper
	f["JsEscape"] = html.JSEscapeString
	f["AzEscape"] = AzEscape
	f["Env"] = Env
	return f
}

func Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}


func AzEscape(args ...interface{}) string {
	result := html.URLQueryEscaper(args...)
	// azure arm wants %20 instead of '+' for spaces
	return strings.Replace(result, "+", "%20", -1)
}
func add(a, b int) int {
	return a + b
}
func minus(a, b int) int {
	return a -b
}
func mult(a, b int) int {
	return a * b
}
func div(a, b int) int {
	return a / b
}

// Series - generate a series of numbers for without pre-allocating a slice  ./bin/gomplate -i '{{  range $idx,$ii := loop 1 4 1 }} ii:{{ $ii}} idx: {{$idx}} {{end}}'
//          ii:1 idx: 0  ii:2 idx: 1  ii:3 idx: 2  ii:4 idx: 3
func Series(start, end, interval int) (stream chan int) {
	stream = make(chan int)
	go func() {
		for i := start; i <= end; i+=interval {
			stream <- i
		}
		close(stream)
	}()
	return
}


func slice(v ...interface{}) []interface{} {
	return v
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

