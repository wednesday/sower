package util

import (
	"strings"
	"sync"
)

type Node struct {
	node
	sep string
	*sync.RWMutex
}
type node struct {
	node map[string]*node
}

func NewNode(sep string) *Node {
	return &Node{node{node: map[string]*node{}}, sep, &sync.RWMutex{}}
}
func NewNodeFromRules(sep string, rules ...string) *Node {
	n := NewNode(sep)
	for i := range rules {
		n.Add(rules[i])
	}
	return n
}

func (n *Node) String() string {
	n.RLock()
	defer n.RUnlock()
	return n.string("", "     ")
}
func (n *node) string(prefix, indent string) (out string) {
	for key, val := range n.node {
		out += prefix + key + "\n" + val.string(prefix+indent, indent)
	}
	return
}

func (n *Node) trim(item string) string {
	return strings.TrimSuffix(item, n.sep)
}

func (n *Node) Add(item string) {
	n.Lock()
	defer n.Unlock()
	n.add(strings.Split(n.trim(item), n.sep))
}
func (n *node) add(secs []string) {
	length := len(secs)
	switch length {
	case 0:
		return
	case 1:
		n.node[secs[length-1]] = &node{node: map[string]*node{}}
	default:
		subNode, ok := n.node[secs[length-1]]
		if !ok {
			subNode = &node{node: map[string]*node{}}
			n.node[secs[length-1]] = subNode
		}
		subNode.add(secs[:length-1])
	}
}

func (n *Node) Match(item string) bool {
	return n.matchSecs(strings.Split(n.trim(item), n.sep))
}

func (n *node) matchSecs(secs []string) bool {
	length := len(secs)
	if length == 0 {
		switch len(n.node) {
		case 0:
			return true
		case 1:
			_, ok := n.node["*"]
			return ok
		default:
			return false
		}
	}

	if n, ok := n.node[secs[length-1]]; ok {
		return n.matchSecs(secs[:length-1])
	}

	_, ok := n.node["*"]
	return ok
}
