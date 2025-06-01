package trie

type Trie struct {
    root *Node
}

type Node struct {
    child map[rune]*Node
    isEnd bool
}

func NewPathTrie() *Trie {
    return &Trie{root: &Node{
        child: make(map[rune]*Node),
    }}
}

func (t *Trie) Insert(content string) {
    node := t.root
    for _, i := range content {
        if _, ok := node.child[i]; !ok {
            node.child[i] = &Node{child: make(map[rune]*Node)}
        }
        node = node.child[i]
    }
    node.isEnd = true
}

func (t *Trie) Match(content string) bool {
    node := t.root
    for _, i := range content {
        wildcardNode, ok := node.child['*']
        if ok && wildcardNode.isEnd {
            return true
        }
        child, ok := node.child[i]
        if !ok {
            return false
        }
        node = child
    }
    return node.isEnd
}
