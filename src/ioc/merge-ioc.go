package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	ind "github.com/tnw-open-source/indicators"
)

func main() {

	inds := ind.IndicatorDefinitions{
		Description: "Trust Networks Indicators",
		Version:     "3"}

	inds.Definitions = make([]*ind.IndicatorNode, 0)

	for _, f := range os.Args[1:] {

		fmt.Fprintf(os.Stderr, "%s...\n", f)

		var part ind.IndicatorDefinitions

		raw, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}

		err = json.Unmarshal(raw, &part)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}

		for _, v := range part.Definitions {
			inds.Definitions = append(inds.Definitions, v)
		}
	}

	fmt.Fprintf(os.Stderr, "Total %v\n", len(inds.Definitions))

	j, err := json.Marshal(&inds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s", j)

}
