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

// Index służy jako indeks łańcucha znaków.
// Pole `text` to indeksowany łańcuch.
// Pole `nodes` to tablica sufiksów łańcucha `text`.
// Pole `root` to korzeń drzewa sufiksów `nodes`.
// Pozostałe pola przydają się podczas budowania tablicy sufiksów
type Index struct {
	text                                          []byte
	nodes                                         []node
	root, position, currentNode, needSuffixLink   int
	backlog, activeNode, activeLength, activeEdge int
}

// New tworzy indeks łańcucha ‘text‘
func New(text []byte) *Index {
	ix := &Index{
		text:     make([]byte, len(text)+1), // Last byte is a sentinel
		nodes:    []node{},
		position: -1,
	}
	copy(ix.text, text)
	ix.root = ix.newNode(-1, -1)
	ix.activeNode = ix.root
	for _, c := range text {
		ix.addChar(c)
	}
	ix.addSuffixStartToLeaves(ix.root, len(text))
	return ix
}

func (ix *Index) addSuffixStartToLeaves(v int, suffixStart int) {
	if ix.IsLeaf(v) {
		ix.nodes[v].suffixStart = suffixStart
	}
	for _, w := range ix.Edges(v) {
		ix.addSuffixStartToLeaves(w, suffixStart-ix.edgeLength(&ix.nodes[w]))
	}
}

// NumNodes zwraca liczbę węzłów drzewa sufiksów `ix.nodes`
func (ix *Index) NumNodes() int {
	return len(ix.nodes)
}

// Root zwraca korzeń drzewa sufiksów `ix.nodes`
func (ix *Index) Root() int {
	return ix.root
}

// Edges zwraca mapę. Elementy tej mapy opisują te krawędzie, które
// wychodzą z wierzchołka `node` drzewa sufiksów `ix.nodes`. Każdy
// klucz tej mapy jest pierwszym znakiem etykiety pewnej krawędzi
// drzewa sufiksów `ix.nodes`, która wychodzi z wierzchołka `node`. Ta
// wartość mapy, która odpowiada temu kluczowi, jest tym wierzchołkiem
// drzewa sufiksów `ix.nodes`, do którego prowadzi ta krawędź
func (ix *Index) Edges(node int) map[byte]int {
	return ix.nodes[node].edges
}

// EdgeLabel zwraca etykietę tej krawędzi, która prowadzi do
// wierzchołka `node` drzewa sufiksów `ix.nodes`
func (ix *Index) EdgeLabel(node int) string {
	return string(ix.text[ix.nodes[node].start:min(ix.nodes[node].end, ix.position+1)])
}

// IsLeaf zwraca true, jeśli węzeł `node` jest liściem drzewa sufiksów
// `ix.nodes`
func (ix *Index) IsLeaf(node int) bool {
	return (len(ix.Edges(node)) == 0)
}

// SuffixStart zwraca indeks tej pozycji, od której zaczyna się ten
// sufiks łańcucha ‘ix.text‘, który jest równy połączonym etykietom
// krawędzi na ścieżce od korzenia do liścia `node` drzewa sufiksów
// `ix.nodes`
func (ix *Index) SuffixStart(node int) int {
	if !ix.IsLeaf(node) {
		log.Fatalf("SuffixStart(%d); not a leaf", node)
	}
	return ix.nodes[node].suffixStart
}

// LookupAll zwraca indeksy wszystkich tych pozycji, od których
// zaczynają się wystąpienia wzorca ‘pat‘ w łańcuchu ‘ix.text‘
func (ix *Index) LookupAll(pat []byte) []int {
	// Based on Listing 6.14 from the book "Algorytmika
	// praktyczna. Nie tylko dla mistrzów" by Piotr Stańczyk,
	// Wydawnictwo Naukowe PWN, Warszawa 2009
	v := ix.root
	patLen := 0
	for i := 0; i < len(pat); {
		var ok bool
		if v, ok = ix.Edges(v)[pat[i]]; !ok {
			return []int{}
		}
		patLen += ix.edgeLength(&ix.nodes[v])
		for x := ix.nodes[v].start; x < ix.nodes[v].end && i < len(pat); {
			if pat[i] != ix.text[x] {
				return []int{}
			}
			i++
			x++
		}
	}
	return ix.dfs([]int{}, v, len(ix.text)-1-patLen)
}

func (ix *Index) dfs(r []int, v, pos int) []int {
	if ix.IsLeaf(v) {
		r = append(r, pos)
	}
	for _, w := range ix.Edges(v) {
		r = ix.dfs(r, w, pos-ix.edgeLength(&ix.nodes[w]))
	}
	return r
}

func (ix *Index) newNode(start, end int) int {
	n := node{start: start, end: end, edges: make(map[byte]int)}
	ix.nodes = append(ix.nodes, n)
	return len(ix.nodes) - 1
}

func (ix *Index) edgeLength(n *node) int {
	return min(n.end, ix.position+1) - n.start
}

func (ix *Index) addSuffixLink(node int) {
	if ix.needSuffixLink > 0 {
		ix.nodes[ix.needSuffixLink].suffixLink = node
	}
	ix.needSuffixLink = node
}

func (ix *Index) walkDown(next int) bool {
	if ix.activeLength >= ix.edgeLength(&ix.nodes[next]) {
		ix.activeEdge += ix.edgeLength(&ix.nodes[next])
		ix.activeLength -= ix.edgeLength(&ix.nodes[next])
		ix.activeNode = next
		return true
	}
	return false
}

func (ix *Index) actE() byte {
	return ix.text[ix.activeEdge]
}

func (ix *Index) addChar(c byte) {
	ix.position++
	ix.needSuffixLink = -1
	ix.backlog++
	for ix.backlog > 0 {
		if ix.activeLength == 0 {
			ix.activeEdge = ix.position
		}
		if _, ok := ix.Edges(ix.activeNode)[ix.actE()]; !ok {
			leaf := ix.newNode(ix.position, oo)
			ix.Edges(ix.activeNode)[ix.actE()] = leaf
			ix.addSuffixLink(ix.activeNode) // rule 2
		} else {
			next := ix.Edges(ix.activeNode)[ix.actE()]
			if ix.walkDown(next) {
				continue // observation 2
			}
			if ix.text[ix.nodes[next].start+ix.activeLength] == c { // observation 1
				ix.activeLength++
				ix.addSuffixLink(ix.activeNode) // observation 3
				break
			}
			split := ix.newNode(ix.nodes[next].start, ix.nodes[next].start+ix.activeLength)
			ix.Edges(ix.activeNode)[ix.actE()] = split
			leaf := ix.newNode(ix.position, oo)
			ix.Edges(split)[c] = leaf
			ix.nodes[next].start += ix.activeLength
			ix.Edges(split)[ix.text[ix.nodes[next].start]] = next
			ix.addSuffixLink(split) // rule 2
		}
		ix.backlog--
		if ix.activeNode == ix.root && ix.activeLength > 0 { // rule 1
			ix.activeLength--
			ix.activeEdge = ix.position - ix.backlog + 1
		} else {
			if ix.nodes[ix.activeNode].suffixLink > 0 { // rule 3
				ix.activeNode = ix.nodes[ix.activeNode].suffixLink
			} else {
				ix.activeNode = ix.root
			}
		}
	}
}

func (ix *Index) Print() {
	fmt.Println("digraph {")
	fmt.Println("\trankdir = LR;")
	fmt.Println("\tedge [arrowsize=0.4,fontsize=10];")
	fmt.Println("\tnode1 [label=\"\",style=filled,fillcolor=lightgrey,shape=circle,width=.1,height=.1];")

	fmt.Println("//------leaves------")
	ix.printLeaves(ix.root)

	fmt.Println("//------internal nodes------")
	ix.printInternalNodes(ix.root)

	fmt.Println("//------edges------")
	ix.printEdges(ix.root)

	fmt.Println("//------suffix links------")
	ix.printSuffixLinks(ix.root)

	fmt.Println("}")
}

func (ix *Index) printLeaves(x int) {
	if ix.IsLeaf(x) {
		fmt.Printf("\tnode%d [label=\"%d\",shape=point];\n", x, ix.SuffixStart(x))
	} else {
		for _, child := range ix.Edges(x) {
			ix.printLeaves(child)
		}
	}
}

func (ix *Index) printInternalNodes(x int) {
	if x != ix.root && !ix.IsLeaf(x) {
		fmt.Printf("\tnode%d [label=\"\",style=filled,fillcolor=lightgrey,shape=circle,width=.07,height=.07];\n", x)
	}
	for _, child := range ix.Edges(x) {
		ix.printInternalNodes(child)
	}
}

func (ix *Index) printEdges(x int) {
	for _, child := range ix.Edges(x) {
		fmt.Printf("\tnode%d -> node%d [label=\"%s\",weight=3];\n", x, child, ix.EdgeLabel(child))
	}
	for _, child := range ix.Edges(x) {
		ix.printEdges(child)
	}
}

func (ix *Index) printSuffixLinks(x int) {
	if ix.nodes[x].suffixLink > 0 {
		fmt.Printf("\tnode%d -> node%d [label=\"\",weight=1,style=dotted];\n", x, ix.nodes[x].suffixLink)
	}
	for _, child := range ix.Edges(x) {
		ix.printSuffixLinks(child)
	}
}
