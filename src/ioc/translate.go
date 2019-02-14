package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	flags "github.com/jessevdk/go-flags"
	dt "github.com/trustnetworks/analytics-common/datatypes"
	ind "github.com/trustnetworks/indicators"
)

var options struct {
	Input  string `short:"i" long:"input-file" description:"Name of the JSON file to use as input" required:"true"`
	Output string `short:"o" long:"output-file" description:"Name of the JSON file to use as output (if empty, writes to stdout)"`
}

var work struct {
	indicators *dt.Indicators
}

func openDefinitions() {

	data, err := ioutil.ReadFile(options.Input)
	if err != nil {
		println(err.Error())
		return
	}

	var ind dt.Indicators
	err = json.Unmarshal(data, &ind)
	if err != nil {
		println(err.Error())
		return
	}

	work.indicators = &ind
}

// TranslateIndicators does what it says on the tin from the old format to the new one
func TranslateIndicators() {

	// Parse the command line
	_, err := flags.Parse(&options)
	if err != nil {
		return
	}

	openDefinitions()

	indicators := make([]*ind.IndicatorNode, 0)

	for _, i := range work.indicators.Indicators {

		if i != nil {
			pattern := &ind.TypeValuePair{
				Type:  i.Type,
				Value: i.Value}
			node := &ind.IndicatorNode{
				ID:        i.Id,
				Indicator: i,
				Pattern:   pattern}

			indicators = append(indicators, node)
		}
	}

	collection := ind.IndicatorDefinitions{
		Description: "Trust Networks Indicators",
		Version:     "3",
		Definitions: indicators}

	j, err := json.Marshal(collection)

	//j = clean(j)

	err = ioutil.WriteFile(options.Output, j, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func clean(data []byte) []byte {

	temp := string(data)
	replace1 := strings.NewReplacer(",\"comment\":\"\",\"and\":null,\"or\":null,\"not\":null", "")
	replace2 := strings.NewReplacer(",\"service\":null,\"tags\":null", "")
	temp = replace1.Replace(temp)
	temp = replace2.Replace(temp)
	data = []byte(temp)

	return data
}

func main() {
	TranslateIndicators()
}
