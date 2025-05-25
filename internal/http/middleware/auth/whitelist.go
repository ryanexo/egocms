package auth

type whitelist struct {
    root *charNode
}

type charNode struct {
    child map[rune]*charNode
    isEnd bool
}

func newWhitelist() *whitelist {
    return &whitelist{root: &charNode{
        child: make(map[rune]*charNode),
    }}
}

func (t *whitelist) insert(content string) {
    node := t.root
    for _, i := range content {
        if _, ok := node.child[i]; !ok {
            node.child[i] = &charNode{child: make(map[rune]*charNode)}
        }
        node = node.child[i]
    }
    node.isEnd = true
}

func (t *whitelist) match(content string) bool {
    node := t.root
    for _, i := range content {
        child, ok := node.child[i]
        
        if !ok {
            wildcardNode, ok := node.child['*']
            if ok {
                return wildcardNode.isEnd
            }
            return false
        }
        
        node = child
    }
    return node.isEnd
}
