package main

// some code copy pasted from https://tpaschalis.github.io/golang-graphviz/

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type edge struct {
	node  string
	label string
}
type graph struct {
	nodes  map[string][]edge
	topics []string
}

func newGraph() *graph {
	return &graph{nodes: make(map[string][]edge)}
}

func (g *graph) addEdge(from, to, label string) {
	g.nodes[from] = append(g.nodes[from], edge{node: to, label: label})
}

func (g *graph) addTopic(topic string) {
	g.topics = append(g.topics, topic)
}

func (g *graph) getEdges(node string) []edge {
	return g.nodes[node]
}

func (e *edge) String() string {
	return fmt.Sprintf("%v", e.node)
}

func (g *graph) String() string {
	out := `
digraph G {
	rankdir=LR;
	layout = dot
`
	for _, t := range g.topics {
		out += fmt.Sprintf(`
	"%s" [ shape = box,color=lightgrey,style=filled ];`, t)
	}

	for k := range g.nodes {
		for _, v := range g.getEdges(k) {
			out += fmt.Sprintf(`
	"%s" -> "%s" [ label = "%s" ];`, k, v.node, v.label)
		}
	}
	out += "\n}"
	return out
}

func main() {
	// csv format: topic,operation,app,lang,version,transport,environment
	csvfile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	g := newGraph()

	// add different shapes for the topics
	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		g.addTopic(record[0])
	}

	// re-read file and add nodes
	csvfile.Seek(0, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// all topics are consumed by mm, remove noise
		if record[2] == "mirrormaker" {
			continue
		}
		// if record[2] == "shopify" {
		// 	continue
		// }

		// use arrows for how the data flows
		if record[1] == "CONSUMER" {
			g.addEdge(record[0], record[2], record[1])
		} else {
			g.addEdge(record[2], record[0], record[1])
		}
	}
	fmt.Println(g)
}
