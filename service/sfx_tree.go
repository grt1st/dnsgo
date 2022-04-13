package service

type suffixTreeNode struct {
	key      string
	value    string
	children map[string]*suffixTreeNode
}

func NewSuffixTreeRoot() *suffixTreeNode {
	return NewSuffixTree("", "")
}

func NewSuffixTree(key string, value string) *suffixTreeNode {
	root := &suffixTreeNode{
		key:      key,
		value:    value,
		children: map[string]*suffixTreeNode{},
	}
	return root
}

func (node *suffixTreeNode) ensureSubTree(key string) {
	if _, ok := node.children[key]; !ok {
		node.children[key] = NewSuffixTree(key, "")
	}
}

func (node *suffixTreeNode) insert(key string, value string) {
	if c, ok := node.children[key]; ok {
		c.value = value
	} else {
		node.children[key] = NewSuffixTree(key, value)
	}
}

func (node *suffixTreeNode) sinsert(keys []string, value string) {
	if len(keys) == 0 {
		return
	}

	key := keys[len(keys)-1]
	if len(keys) > 1 {
		node.ensureSubTree(key)
		node.children[key].sinsert(keys[:len(keys)-1], value)
		return
	}

	node.insert(key, value)
}

func (node *suffixTreeNode) search(keys []string) (string, bool) {

	if len(keys) == 0 {
		return "", false
	}

	key := keys[len(keys)-1]
	if n, ok := node.children[key]; ok {
		if nextValue, found := n.search(keys[:len(keys)-1]); found {
			return nextValue, found
		}
		return n.value, (n.value != "")
	}

	return "", false
}

func (node *suffixTreeNode) delete(keys []string) {
	if len(keys) == 0 {
		return
	}

	key := keys[len(keys)-1]
	if len(keys) > 1 {
		node.ensureSubTree(key)
		node.children[key].delete(keys[:len(keys)-1])
		return
	}

	node.value = ""
}

func (node *suffixTreeNode) searchWidcard(keys []string) (string, bool) {

	if len(keys) == 0 {
		return "", false
	}

	key := keys[len(keys)-1]
	n, ok := node.children[key]
	if ok == false {
		n, ok = node.children["*"]
	}
	if ok {
		if nextValue, found := n.searchWidcard(keys[:len(keys)-1]); found {
			return nextValue, found
		}
		return n.value, n.value != ""
	}

	return "", false
}