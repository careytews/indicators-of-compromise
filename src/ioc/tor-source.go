package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const name = "Tor"
const infile = "ioc-sources/tor/tors.html"
const outfile = "output/tors.json"
const entryCategory = "tor.entry"
const exitCategory = "tor.exit"
const email = "me@dan.me.uk"
const url = "https://www.dan.me.uk/tornodes"
const start = "<!-- __BEGIN_TOR_NODE_LIST__ //-->"
const end = "<!-- __END_TOR_NODE_LIST__ //-->"
const srcExitPortRange = "210b48ef-9d4a-486e-86c1-bf4b89431102"
const destExitPortRange = "7847a7b1-fc10-498c-a421-d2795806491b"

// TorNode is all of the information about a Tor Node
// included in the list at https://www.dan.me.uk/tornodes
type TorNode struct {
	ID    string
	IP    string
	Ports []string
	Exit  bool
}

var data struct {
	nodes            map[string]*TorNode
	entryDefinitions []*IndicatorNode
	exitDefinitions  []*IndicatorNode
}

func initialise() error {

	err := fetchInit(infile, url, 30)
	if err != nil {
		return err
	}

	data.nodes = make(map[string]*TorNode, 0)
	data.entryDefinitions = make([]*IndicatorNode, 0)
	data.exitDefinitions = make([]*IndicatorNode, 0)

	return nil
}

func process() error {

	b, err := ioutil.ReadFile(infile)
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

		if strings.Contains(line, start) {
			break
		}
	}

	for true {
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil
		}

		if !strings.Contains(line, end) {
			err = processLine(line)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func processLine(line string) error {

	replacer := strings.NewReplacer("&lt;", "(", "&gt;", ")", "<br />", "", "\n", "")
	line = replacer.Replace(line)
	fields := strings.Split(line, "|")

	id := sha1.New()
	io.WriteString(id, line)

	ip := fields[0]
	port := fields[2]
	exit := false

	if strings.ContainsAny(fields[4], "EX") {
		exit = true
	}

	node := data.nodes[ip]

	if node != nil {
		node.Ports = append(node.Ports, port)
		node.Exit = exit
	} else {
		newNode := &TorNode{
			ID:   fmt.Sprintf("%x", id.Sum(nil)),
			IP:   ip,
			Exit: exit}

		newNode.Ports = append(newNode.Ports, port)
		data.nodes[ip] = newNode
	}

	return nil
}

func postProcess() error {

	for _, node := range data.nodes {
		addEntryNode(node)
		if node.Exit {
			addExitNode(node)
		}
	}

	collection := &IndicatorDefinitions{
		Description: "Tor nodes",
		Version:     "3"}

	entryNodes := &IndicatorNode{Operator: "OR"}
	entryNodes.Indicator = &Indicator{
		Id:          "TorEntryNodes-133c7a77-81dd-4d05-babb-fb2e9344d9cc",
		Description: "TOR entry node",
		Category:    entryCategory,
		Author:      email,
		Source:      url}
	entryNodes.Children = data.entryDefinitions

	collection.Definitions = append(collection.Definitions, entryNodes)

	exitNodes := &IndicatorNode{Operator: "OR"}
	exitNodes.Indicator = &Indicator{
		Id:          "TorExitNodes-133c7a77-81dd-4d05-babb-fb2e9344d9cc",
		Description: "TOR exit node",
		Category:    exitCategory,
		Author:      email,
		Source:      url}
	exitNodes.Children = data.exitDefinitions

	collection.Definitions = append(collection.Definitions, exitNodes)

	err := save(collection, outfile)
	if err != nil {
		return err
	}

	return nil
}

func addEntryNode(tor *TorNode) {

	id := fmt.Sprintf("TorEntryNode-%s", tor.ID)

	// Create top level node
	definition := &IndicatorNode{Operator: "OR"}

	if tor.Exit {
		definition.ID = id
	}

	srcType := "src.ipv4"
	destType := "dest.ipv4"
	if strings.Contains(tor.IP, ":") {
		srcType = "src.ipv6"
		destType = "dest.ipv6"
	}

	// Create the definition's children
	srcChild := &IndicatorNode{Operator: "AND"}
	srcChild.Children = append(srcChild.Children, createPatternNode(srcType, tor.IP, "", ""))
	if len(tor.Ports) > 1 {
		orChild := &IndicatorNode{Operator: "OR"}
		for _, child := range tor.Ports {
			orChild.Children = append(orChild.Children, createPatternNode("src.tcp", child, "", "int"))
		}
		srcChild.Children = append(srcChild.Children, orChild)
	} else {
		srcChild.Children = append(srcChild.Children, createPatternNode("src.tcp", tor.Ports[0], "", "int"))
	}

	destChild := &IndicatorNode{Operator: "AND"}
	destChild.Children = append(destChild.Children, createPatternNode(destType, tor.IP, "", ""))
	if len(tor.Ports) > 1 {
		orChild := &IndicatorNode{Operator: "OR"}
		for _, child := range tor.Ports {
			orChild.Children = append(orChild.Children, createPatternNode("dest.tcp", child, "", "int"))
		}
		destChild.Children = append(destChild.Children, orChild)
	} else {
		destChild.Children = append(destChild.Children, createPatternNode("dest.tcp", tor.Ports[0], "", "int"))
	}

	// Append the definition's children
	definition.Children = append(definition.Children, srcChild)
	definition.Children = append(definition.Children, destChild)

	data.entryDefinitions = append(data.entryDefinitions, definition)
}

func addExitNode(tor *TorNode) {

	entryID := fmt.Sprintf("TorEntryNode-%s", tor.ID)

	// Create the top level node
	definition := &IndicatorNode{Operator: "AND"}

	srcType := "src.ipv4"
	destType := "dest.ipv4"
	if strings.Contains(tor.IP, ":") {
		srcType = "src.ipv6"
		destType = "dest.ipv6"
	}

	// Create the definition's children
	orChild := &IndicatorNode{Operator: "OR"}

	srcChild := &IndicatorNode{Operator: "AND"}
	srcChild.Children = append(srcChild.Children, createPatternNode(srcType, tor.IP, "", ""))
	srcChild.Children = append(srcChild.Children, &IndicatorNode{Ref: srcExitPortRange})

	destChild := &IndicatorNode{Operator: "AND"}
	destChild.Children = append(destChild.Children, createPatternNode(destType, tor.IP, "", ""))
	destChild.Children = append(destChild.Children, &IndicatorNode{Ref: destExitPortRange})

	orChild.Children = append(orChild.Children, srcChild)
	orChild.Children = append(orChild.Children, destChild)

	notChild := &IndicatorNode{Operator: "NOT"}
	notChild.Children = append(notChild.Children, &IndicatorNode{Ref: entryID})

	// Append the definition's children
	definition.Children = append(definition.Children, orChild)
	definition.Children = append(definition.Children, notChild)

	data.exitDefinitions = append(data.exitDefinitions, definition)
}

func main() {
	run(name)
}
