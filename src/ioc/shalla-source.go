package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	dt "github.com/trustnetworks/analytics-common/datatypes"
	ind "github.com/trustnetworks/indicators"
)

var categories = []string{
	"anonvpn",
	"dynamic",
	"hacking",
	"redirector",
	"spyware",
	"warez"}

const name = "Shallalist"
const indir = "ioc-shallalist"
const domains = "domains"
const urls = "urls"
const url = "http://www.shallalist.de/"
const email = "info@shallalist.de"

var data struct {
	definitions map[string][]*ind.IndicatorNode
}

func initialise() error {

	_, err := os.Stat(indir)
	if err != nil {
		return err
	}

	data.definitions = make(map[string][]*ind.IndicatorNode, 0)

	return nil
}

func process() error {

	for _, category := range categories {

		directory := indir + "/" + category

		_, err := os.Stat(directory)
		if err != nil {
			panic(err)
		}

		err = processCategory(directory, category)
		if err != nil {
			return err
		}
	}

	return nil
}

func processCategory(directory string, category string) error {

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, file := range files {

		err = processFile(category, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func processFile(category string, file os.FileInfo) error {

	filename := indir + "/" + category + "/" + file.Name()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	reader := bytes.NewBuffer(b)
	var line string

	for true {
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil
		}

		line = strings.Trim(line, "\n")

		addIndicator(category, file.Name(), line)
	}

	return nil
}

func addIndicator(category string, file string, line string) {

	var tp string
	var match string

	if file == domains {
		tp = "hostname"
		match = "dns"
	} else {
		tp = "url"
		match = ""
	}

	node := createPatternNode(tp, line, "", match)

	data.definitions[category] = append(data.definitions[category], node)
}

func postProcess() error {

	for _, category := range categories {
		collection := &ind.IndicatorDefinitions{
			Description: name + " " + category,
			Version:     "3"}

		orList := &ind.IndicatorNode{Operator: "OR"}
		orList.Indicator = &dt.Indicator{
			Id:          name + "-" + category + "-67599039-b461-4312-96dd-d01e2bcfc380",
			Description: name + " " + category,
			Category:    category,
			Author:      email,
			Source:      url}
		orList.Children = data.definitions[category]

		collection.Definitions = append(collection.Definitions, orList)

		fmt.Printf("Shallalist %s: %v\n", category, len(orList.Children))

		err := save(collection, "output/shalla-"+category+".json")
		if err != nil {
			return err
		}

	}

	return nil
}

func main() {
	run(name)
}
