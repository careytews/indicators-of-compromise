package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"

	"github.com/google/uuid"
)

const name = "TargetedThreats"
const infile = "targetedthreats/targetedthreats.csv"
const outfile = "output/tt.json"

// TargetedThreat is all of the information about a targeted threat
// in the repository at https://github.com/botherder/targetedthreats
type TargetedThreat struct {
	Address   string `json:"address"`
	Family    string `json:"family"`
	Country   string `json:"country"`
	Reference string `json:"reference"`
}

var data struct {
	threats     []*TargetedThreat
	definitions []*IndicatorNode
}

func initialise() error {

	data.definitions = make([]*IndicatorNode, 0)

	return nil
}

func process() error {

	data.threats = make([]*TargetedThreat, 0)

	in, err := os.Open(infile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(bufio.NewReader(in))

	// Ignore the first line
	line, err := reader.Read()
	if err == io.EOF {
		return nil
	}

	for true {

		line, err = reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil
		}

		processLine(line)
	}

	return nil
}

func processLine(record []string) error {

	u, err := url.Parse(record[3])
	if err != nil {
		return err
	}

	ref := u.Scheme + "://" + u.Host
	threat := &TargetedThreat{
		Address:   record[0],
		Family:    record[1],
		Country:   record[2],
		Reference: ref}

	data.threats = append(data.threats, threat)

	return nil
}

func postProcess() error {

	for _, threat := range data.threats {
		addIndicator(threat)
	}

	collection := &IndicatorDefinitions{
		Description: "Botherder Targeted Threats",
		Version:     "3",
		Definitions: data.definitions}

	fmt.Printf("Targeted Threats: %v\n", len(data.definitions))

	err := save(collection, outfile)
	if err != nil {
		return err
	}

	return nil
}

func addIndicator(threat *TargetedThreat) {

	ipAddress := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`)

	if ipAddress.MatchString(threat.Address) {
		addIPNode(threat)
	} else {
		addHostnameNode(threat)
	}
}

func addIPNode(threat *TargetedThreat) {

	// Create top level node
	definition := &IndicatorNode{Operator: "OR"}

	// Create Indicator
	definition.Indicator = &Indicator{
		Id:          uuid.New().String(),
		Type:        "ipv4",
		Value:       threat.Address,
		Description: threat.Family + ":" + threat.Country,
		Category:    "apt",
		Author:      "botherder@github",
		Source:      threat.Reference}

	// Create the definition's children
	definition.Children = append(definition.Children, createPatternNode("src.ipv4", threat.Address, "", ""))
	definition.Children = append(definition.Children, createPatternNode("dest.ipv4", threat.Address, "", ""))

	data.definitions = append(data.definitions, definition)
}

func addHostnameNode(threat *TargetedThreat) {

	// Create the top level node
	definition := &IndicatorNode{}

	// Create indicator
	definition.Indicator = &Indicator{
		Id:          uuid.New().String(),
		Type:        "hostname",
		Value:       threat.Address,
		Description: threat.Family + ":" + threat.Country,
		Category:    "apt",
		Author:      "botherder@github",
		Source:      threat.Reference}

	definition.Pattern = &Pattern{
		Type:  "hostname",
		Value: threat.Address,
		Match: "dns"}

	data.definitions = append(data.definitions, definition)
}

func main() {
	run(name)
}
