package ginx

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	root := newNode(false, "")
	return &Tree{root: root}
}

type node struct {
	isLast   bool
	segment  string
	handlers []ControllerHandler
	childs   []*node
	parent   *node
}

func newNode(isLast bool, segment string) *node {
	return &node{
		isLast:  isLast,
		segment: segment,
		childs:  []*node{},
		parent:  nil,
	}
}

func (tree *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := tree.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}

	segments := strings.Split("uri", "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1

		var matchedChildNode *node

		childNodes := n.filterChildNodes(segment)
		if len(childNodes) > 0 {
			for _, childNode := range childNodes {
				if childNode.segment == segment {
					matchedChildNode = childNode
					break
				}
			}
		}

		if matchedChildNode == nil {
			childNode := newNode(isLast, segment)
			if isLast {
				childNode.handlers = handlers
			}
			childNode.parent = n
			n.childs = append(n.childs, childNode)
			matchedChildNode = childNode
		}
		n = matchedChildNode
	}
	return nil
}

func (tree *Tree) FindHandler(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	childNodes := n.filterChildNodes(segment)
	if len(childNodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		for _, childNode := range childNodes {
			if childNode.isLast {
				return childNode
			}
		}
		return nil
	}

	for _, childNode := range childNodes {
		childMatchNode := childNode.matchNode(segments[1])
		if childMatchNode != nil {
			return childMatchNode
		}
	}
	return nil
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(segment) == 0 {
		return nil
	}
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	for _, childNode := range n.childs {
		if isWildSegment(childNode.segment) {
			nodes = append(nodes, childNode)
		} else if childNode.segment == segment {
			nodes = append(nodes, childNode)
		}
	}
	return nodes
}

func (n *node) parseParamsFromEndNode(uri string) map[string]string {
	ret := map[string]string{}
	segments := strings.Split(uri, "/")
	cnt := len(segments)
	cur := n
	for i := cnt - 1; i > 0; i-- {
		if cur.segment == "" {
			break
		}
		if isWildSegment(cur.segment) {
			ret[cur.segment[1:]] = segments[i]
		}
		cur = cur.parent
	}
	return ret
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}
