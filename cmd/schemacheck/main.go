package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/muncus/jsonschemavalidator/output"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

type ArrayFlag []string

func (f *ArrayFlag) String() string {
	return strings.Join(*f, ",")
}
func (f *ArrayFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var schemafile = flag.String("s", "", "Directory to load json schemas from")
var outputformat = flag.String("output", "github", "output format: 'text' or 'github'")
var docs ArrayFlag

func loaderForFile(fname string) (gojsonschema.JSONLoader, error) {
	absfile, err := filepath.Abs(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve abs path %s: %v", fname, err)
	}
	dl := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", absfile))
	_, err = dl.LoadJSON()
	if err == nil {
		// successfully loaded json, return this loader.
		return dl, nil
	}
	// Otherwise, try to parse it as yaml.
	ydata, err := os.ReadFile(absfile)
	yobj := make(map[string]interface{})
	err = yaml.Unmarshal(ydata, &yobj)
	if err != nil {
		return nil, err
	}
	yl := gojsonschema.NewGoLoader(yobj)
	return yl, nil

}

func main() {
	log.SetFlags(0)
	flag.Var(&docs, "d", "JSON or YAML document to validate")
	flag.Parse()

	// TODO: allow limited use of http urls? (e.g. schemastore.org)
	var schemasource string
	// canonicalize the provided schema path.
	if _, err := url.Parse(*schemafile); err == nil {
		// url is valid. use it.
		schemasource = *schemafile
	} else {
		abspath, _ := filepath.Abs(*schemafile)
		schemasource = fmt.Sprintf("file://%s", abspath)
	}
	sl := gojsonschema.NewReferenceLoader(schemasource)
	schema, err := gojsonschema.NewSchema(sl)
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}
	for _, input := range flag.Args() {
		// TODO: consider something different for outputs. (e.g. TAP/junit)
		l, err := loaderForFile(input)
		if err != nil {
			log.Printf("failed to load %s: %v\n", input, err)
			continue
		}
		result, err := schema.Validate(l)
		if err != nil {
			log.Printf("error during validation: %v\n", err)
			continue
		}
		if *outputformat == "github" {
			output.GithubOutput(os.Stdout, input, result)
		} else {
			output.TextOutput(os.Stdout, input, result)
		}
	}
}
