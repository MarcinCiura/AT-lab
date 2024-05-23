/*
 * Based on Sergey Makagonov's Java class ST, available at
 * https://gist.github.com/makagonov/22ab3675e3fc0031314e8535ffcbee2c
 * and on his C++ code, available at
 * https://gist.github.com/makagonov/f7ed8ce729da72621b321f0ab547debb
 * under the following license:
 *
 * Copyright (c) 2016 Sergey Makagonov
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the
 * "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use, copy, modify, merge, publish,
 * distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to
 * the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 * LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
 * OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package suffixtree

import (
	"fmt"
	"log"
	"math"
)

const (
	oo = math.MaxInt
)

type node struct {
	start, end  int
	suffixStart int // Used only in leaves
	suffixLink  int
	edges       map[byte]int
}

type Index struct {
	nodes                                                []node
	text                                                 []byte
	root, position, currentNode, needSuffixLink, backlog int
	activeNode, activeLength, activeEdge                 int
}

// New zwraca indeks łańcucha ‘text‘
func New(text []byte) *Index {
	st := &Index{
		nodes:    []node{},
		text:     make([]byte, len(text)+1), // Last byte is a sentinel
		position: -1,
	}
	copy(st.text, text)
	st.root = st.newNode(-1, -1)
	st.activeNode = st.root
	for _, c := range text {
		st.addChar(c)
	}
	st.addSuffixStartToLeaves(st.root, len(text))
	return st
}

func (st *Index) addSuffixStartToLeaves(v int, suffixStart int) {
	if st.IsLeaf(v) {
		st.nodes[v].suffixStart = suffixStart
	}
	for _, w := range st.Edges(v) {
		st.addSuffixStartToLeaves(w, suffixStart-st.edgeLength(&st.nodes[w]))
	}
}

// NumNodes zwraca liczbę węzłów drzewa sufiksów
func (st *Index) NumNodes() int {
	return len(st.nodes)
}

// Root zwraca indeks korzenia drzewa sufiksów
func (st *Index) Root() int {
	return st.root
}

// Edges zwraca mapę. Elementy tej mapy opisują te krawędzie, które
// wychodzą z tego wierzchołka drzewa sufiksów, który ma indeks
// ‘node‘. Każdy klucz tej mapy to pierwszy znak etykiety pewnej
// krawędzi. Wartość mapy, która odpowiada temu kluczowi, to indeks
// tego wierzchołka drzewa sufiksów, do którego prowadzi ta krawędź
func (st *Index) Edges(node int) map[byte]int {
	return st.nodes[node].edges
}

// EdgeLabel zwraca etykietę tej krawędzi, która prowadzi do tego
// wierzchołka drzewa sufiksów, który ma indeks ‘node‘
func (st *Index) EdgeLabel(node int) string {
	return string(st.text[st.nodes[node].start:min(st.nodes[node].end, st.position+1)])
}

// IsLeaf zwraca true, jeśli ten węzeł drzewa sufiksów, który ma
// indeks ‘node‘, jest liściem tego drzewa
func (st *Index) IsLeaf(node int) bool {
	return (len(st.Edges(node)) == 0)
}

// SuffixStart zwraca indeks początku tego sufiksu łańcucha ‘st.text‘,
// który jest równy połączonym etykietom krawędzi na ścieżce od
// korzenia drzewa sufiksów do tego liścia tego drzewa, który ma
// indeks ‘node‘
func (st *Index) SuffixStart(node int) int {
	if !st.IsLeaf(node) {
		log.Fatalf("SuffixStart(%d); not a leaf", node)
	}
	return st.nodes[node].suffixStart
}

// LookupAll zwraca indeksy wszystkich tych pozycji, od których
// zaczynają się wystąpienia wzorca ‘pat‘ w łańcuchu ‘st.text‘
func (st *Index) LookupAll(pat []byte) []int {
	// Based on Listing 6.14 from the book "Algorytmika
	// praktyczna. Nie tylko dla mistrzów" by Piotr Stańczyk,
	// Wydawnictwo Naukowe PWN, Warszawa 2009
	v := st.root
	patLen := 0
	for i := 0; i < len(pat); {
		var ok bool
		if v, ok = st.Edges(v)[pat[i]]; !ok {
			return []int{}
		}
		patLen += st.edgeLength(&st.nodes[v])
		for x := st.nodes[v].start; x < st.nodes[v].end && i < len(pat); {
			if pat[i] != st.text[x] {
				return []int{}
			}
			i++
			x++
		}
	}
	return st.dfs([]int{}, v, len(st.text)-1-patLen)
}

func (st *Index) dfs(r []int, v, pos int) []int {
	if st.IsLeaf(v) {
		r = append(r, pos)
	}
	for _, w := range st.Edges(v) {
		r = st.dfs(r, w, pos-st.edgeLength(&st.nodes[w]))
	}
	return r
}

func (st *Index) newNode(start, end int) int {
	n := node{start: start, end: end, edges: make(map[byte]int)}
	st.nodes = append(st.nodes, n)
	return len(st.nodes) - 1
}

func (st *Index) edgeLength(n *node) int {
	return min(n.end, st.position+1) - n.start
}

func (st *Index) addSuffixLink(node int) {
	if st.needSuffixLink > 0 {
		st.nodes[st.needSuffixLink].suffixLink = node
	}
	st.needSuffixLink = node
}

func (st *Index) walkDown(next int) bool {
	if st.activeLength >= st.edgeLength(&st.nodes[next]) {
		st.activeEdge += st.edgeLength(&st.nodes[next])
		st.activeLength -= st.edgeLength(&st.nodes[next])
		st.activeNode = next
		return true
	}
	return false
}

func (st *Index) actE() byte {
	return st.text[st.activeEdge]
}

func (st *Index) addChar(c byte) {
	st.position++
	st.needSuffixLink = -1
	st.backlog++
	for st.backlog > 0 {
		if st.activeLength == 0 {
			st.activeEdge = st.position
		}
		if _, ok := st.Edges(st.activeNode)[st.actE()]; !ok {
			leaf := st.newNode(st.position, oo)
			st.Edges(st.activeNode)[st.actE()] = leaf
			st.addSuffixLink(st.activeNode) // rule 2
		} else {
			next := st.Edges(st.activeNode)[st.actE()]
			if st.walkDown(next) {
				continue // observation 2
			}
			if st.text[st.nodes[next].start+st.activeLength] == c { // observation 1
				st.activeLength++
				st.addSuffixLink(st.activeNode) // observation 3
				break
			}
			split := st.newNode(st.nodes[next].start, st.nodes[next].start+st.activeLength)
			st.Edges(st.activeNode)[st.actE()] = split
			leaf := st.newNode(st.position, oo)
			st.Edges(split)[c] = leaf
			st.nodes[next].start += st.activeLength
			st.Edges(split)[st.text[st.nodes[next].start]] = next
			st.addSuffixLink(split) // rule 2
		}
		st.backlog--
		if st.activeNode == st.root && st.activeLength > 0 { // rule 1
			st.activeLength--
			st.activeEdge = st.position - st.backlog + 1
		} else {
			if st.nodes[st.activeNode].suffixLink > 0 { // rule 3
				st.activeNode = st.nodes[st.activeNode].suffixLink
			} else {
				st.activeNode = st.root
			}
		}
	}
}

func (st *Index) Print() {
	fmt.Println("digraph {")
	fmt.Println("\trankdir = LR;")
	fmt.Println("\tedge [arrowsize=0.4,fontsize=10];")
	fmt.Println("\tnode1 [label=\"\",style=filled,fillcolor=lightgrey,shape=circle,width=.1,height=.1];")

	fmt.Println("//------leaves------")
	st.printLeaves(st.root)

	fmt.Println("//------internal nodes------")
	st.printInternalNodes(st.root)

	fmt.Println("//------edges------")
	st.printEdges(st.root)

	fmt.Println("//------suffix links------")
	st.printSuffixLinks(st.root)

	fmt.Println("}")
}

func (st *Index) printLeaves(x int) {
	if st.IsLeaf(x) {
		fmt.Printf("\tnode%d [label=\"%d\",shape=point];\n", x, st.SuffixStart(x))
	} else {
		for _, child := range st.Edges(x) {
			st.printLeaves(child)
		}
	}
}

func (st *Index) printInternalNodes(x int) {
	if x != st.root && !st.IsLeaf(x) {
		fmt.Printf("\tnode%d [label=\"\",style=filled,fillcolor=lightgrey,shape=circle,width=.07,height=.07];\n", x)
	}
	for _, child := range st.Edges(x) {
		st.printInternalNodes(child)
	}
}

func (st *Index) printEdges(x int) {
	for _, child := range st.Edges(x) {
		fmt.Printf("\tnode%d -> node%d [label=\"%s\",weight=3];\n", x, child, st.EdgeLabel(child))
	}
	for _, child := range st.Edges(x) {
		st.printEdges(child)
	}
}

func (st *Index) printSuffixLinks(x int) {
	if st.nodes[x].suffixLink > 0 {
		fmt.Printf("\tnode%d -> node%d [label=\"\",weight=1,style=dotted];\n", x, st.nodes[x].suffixLink)
	}
	for _, child := range st.Edges(x) {
		st.printSuffixLinks(child)
	}
}
