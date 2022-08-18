package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

var schemadir = flag.String("s", "", "Directory to load json schemas from")
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

	sdir, err := filepath.Abs(*schemadir)
	if err != nil {
		log.Fatalf("failed to resolve schema dir: %v", err)
	}
	// TODO: allow limited use of http urls? (e.g. schemastore.org)
	sl := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", sdir))
	schema, err := gojsonschema.NewSchema(sl)
	if err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	for _, input := range docs {
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
		if result.Valid() {
			log.Printf("%s: 🎉 Success!\n", input)
			continue
		} else {
			log.Printf("%s: ❌ Validation Failed:\n", input)
			for _, e := range result.Errors() {
				log.Printf("- %v\n", e)
			}
		}
	}
}
