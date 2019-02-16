package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	ind "github.com/tnw-open-source/indicators"
)

// Regular expressions to use when cleaning the files
var empty = regexp.MustCompile(`\,\"\w+\"\:\"{2}|\"\w+\"\:\"{2}\,`)
var null = regexp.MustCompile(`\,\"\w+\"\:null|\"\w+\"\:null\,`)

func clean(input []byte) string {

	output := string(input)
	output = empty.ReplaceAllString(output, "")
	output = null.ReplaceAllString(output, "")

	return output
}

func fetch(filename string, url string) error {

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Run does what it says on the tin
func run(name string) error {

	err := initialise()
	if err != nil {
		n := fmt.Errorf("%q initialise failed: %q", name, err)
		return n
	}

	err = process()
	if err != nil {
		n := fmt.Errorf("%q process failed: %q", name, err)
		return n
	}

	err = postProcess()
	if err != nil {
		n := fmt.Errorf("%q postProcess failed: %q", name, err)
		return n
	}

	return nil
}

func fetchInit(filename string, url string, minutes float64) error {

	info, err := os.Stat(filename)
	if err != nil && os.IsExist(err) {
		return err
	}

	if os.IsNotExist(err) || time.Since(info.ModTime()).Minutes() > minutes {
		err = fetch(filename, url)
		if err != nil && os.IsExist(err) {
			return err
		}
	}

	return nil
}

func save(definitions *ind.IndicatorDefinitions, outfile string) error {

	j, err := json.Marshal(definitions)
	if err != nil {
		return err
	}

	cleaned := clean(j)
	j = []byte(cleaned)

	err = ioutil.WriteFile(outfile, j, 0644)

	return nil
}

func createPatternNode(patternType string, value string, value2 string, match string) *ind.IndicatorNode {

	pattern := &ind.Pattern{Type: patternType, Value: value}

	if value2 != "" {
		pattern.Value2 = value2
	}

	if match != "" {
		pattern.Match = match
	}

	node := &ind.IndicatorNode{Pattern: pattern}

	return node
}
