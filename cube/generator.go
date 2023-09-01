package cube

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	edges         map[string]*Node
	outboundEdges []string
	id            string
}

type Generator struct {
	TransformStack []int // transform 0 applies to node 0 -> 1
	nodesStack     []*Node
	depth          int
	transformNum   int
}

func (generator *Generator) Next() string {
	res := generator.GetCurrentString()
	generator.transformNum += 1
	for i := len(generator.nodesStack) - 1; i >= 0; i-- {
		nextIndex := generator.TransformStack[i] + 1
		if nextIndex != len(generator.nodesStack[i].outboundEdges) { // if current index can be incremented
			generator.TransformStack[i] = nextIndex
			generator.restructure(i + 1)
			break
		} else if i == 0 {
			generator.depth += 1
			generator.nodesStack = append(generator.nodesStack, nil)
			generator.TransformStack = append(generator.TransformStack, 0)
			generator.restructure(0) // reform entire tree of [0, 0, 0, 0, ...]
		}
	}
	return res
}

func (generator *Generator) restructure(s int) {
	if s == 0 {
		generator.TransformStack[0] = 0
		s = 1
	}
	for i := s; i < generator.depth; i++ {
		prevNode := generator.nodesStack[i-1]
		prevNodeTransform := prevNode.outboundEdges[generator.TransformStack[i-1]]
		generator.nodesStack[i] = prevNode.edges[prevNodeTransform]
		generator.TransformStack[i] = 0
	}
}

func (generator *Generator) GetCurrentString() string {
	res := strings.Builder{}
	for i := 0; i < len(generator.TransformStack); i++ {
		res.WriteString(generator.nodesStack[i].outboundEdges[generator.TransformStack[i]])
	}
	return res.String()
}

func (generator *Generator) GetStackEncoded() string {
	b := make([]string, len(generator.TransformStack))
	for i, stackElement := range generator.TransformStack {
		b[i] = strconv.Itoa(stackElement)
	}
	return strings.Join(b, ",")
}

func (generator *Generator) GetCurrentDepth() int {
	return generator.depth
}

func (generator *Generator) GetCurrentTransformNum() int {
	return generator.transformNum - 1
}

func CreateNewGenerator(stack []int, transformNo int, file string) Generator {
	g := new(Generator)
	g.TransformStack = stack
	g.depth = len(stack)
	g.transformNum = transformNo
	startNode := createGraphFromFile(file)
	nodeStack := []*Node{&startNode}
	for i := 0; i < len(stack)-1; i++ {
		currentNode := nodeStack[i]
		transform := currentNode.outboundEdges[stack[i]]
		nextNode := currentNode.edges[transform]
		nodeStack = append(nodeStack, nextNode)
	}
	g.nodesStack = nodeStack
	return *g
}

func createGraphFromFile(file string) Node {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Unable to read input file graph.csv\n", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("Couldn't close file graph.csv")
		}
	}(f)

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(records) < 1 || len(records[0]) != len(records) {
		log.Fatal("Graph matrix must be square and have length >= 1")
	}

	nodeList := make([]Node, len(records))
	for i := 0; i < len(nodeList); i++ {
		nodeList[i] = Node{id: fmt.Sprintf("id:%d", i)}
	}

	for fromNodeIndex, nodeEdges := range records {
		fromNode := nodeList[fromNodeIndex]
		fromNode.edges = make(map[string]*Node)
		for toNodeIndex, edge := range nodeEdges {
			if edge != "_" {
				fromNode.edges[edge] = &nodeList[toNodeIndex]
			}
		}
		fromNode.outboundEdges = make([]string, len(fromNode.edges))
		i := 0
		for k := range fromNode.edges {
			fromNode.outboundEdges[i] = k
			i++
		}
		sort.Strings(fromNode.outboundEdges)
		nodeList[fromNodeIndex] = fromNode
	}

	return nodeList[0]
}
