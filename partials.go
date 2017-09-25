package main

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"fmt"
	"text/template"
	"path"
	"strings"
)
var (
	supported  map[string]bool
)
func init(){
	supported = map[string]bool{".tpl": true, ".html": true, ".tmpl": true, ".txt": true,".t": true}
}

// == Recursive input dir processing ======================================

func partialTemplate(filename string , flattenName bool,g *Gomplate) error{
	inFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s\n%v", filename, err)
	}
	// nolint: errcheck
	defer inFile.Close()
	bytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		fmt.Errorf("read failed for %s\n%v", filename, err)
		return err
	}
	// template name == filename
	// {{ template "resources/foo.t" }}
	name:= filename
	if flattenName {
		name = path.Base(filename)
	}
	nest :=g.Template.New(name)
	nest.Delims(g.leftDelim,g.rightDelim)
	if _,err = nest.Parse(string(bytes)) ; err != nil {
		keys:= make([]string,0)
		for key  := range map[string]interface{} (g.funcMap){
			keys = append(keys,key)
		}
		return fmt.Errorf("Error parsing template %s - %s. Available  functions: %+v", filename,err.Error(), keys )
	}
	if g.verbose {
		fmt.Fprintf(os.Stderr, "Added template name: %s  path: %s \n", name, filename)
	}
	// added to template's parse tree by Parse
	return nil
}
func partialInputDir(inputIn string,  flattenName bool, g *Gomplate) error {
	if g.Template == nil {
		if g.verbose {
			fmt.Fprintf(os.Stderr, "Creating toplevel path: %s \n", inputIn )
		}
		g.Template = template.New("template")
		g.Template.Funcs(g.funcMap)
	}
	for _,input := range strings.Split(inputIn, ":") {
		input = filepath.Clean(input)

		// assert tha input path exists
		st, err := os.Stat(input)
		if err != nil {
			return err
		}
		if !st.IsDir() {
			return fmt.Errorf("%s is not a directory", input)
		}
		// read directory
		entries, err := ioutil.ReadDir(input)
		if err != nil {
			return err
		}

		// process or dive in again
		for _, entry := range entries {
			nextInPath := filepath.Join(input, entry.Name())
			extension := filepath.Ext(nextInPath)
			if g.verbose {
				fmt.Fprintf(os.Stderr, "asset path(%s)  ext: %s\n", nextInPath, extension)

			}

			if _, ok := supported[extension]; !ok {
				if g.verbose {
					fmt.Fprintf(os.Stderr, "Skipping template %s - %s\n", nextInPath, extension)
				}
				continue
			}

			if entry.IsDir() {
				err := partialInputDir(nextInPath, flattenName, g)
				if err != nil {
					return err
				}
			} else {
				if err = partialTemplate(nextInPath, flattenName, g); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
