package core


type Graph struct {
    Nodes map[interface{}]*Node
    
    NamedNodes map[string]*Node
    
    eager map[string]*Node
    
    sources map[interface{}]map[interface{}]Node
    
    targets map[interface{}]map[interface{}]Node
}



func (g *Graph) GetByName(name string) *Node {
    return g.NamedNodes[name]
}

func (g *Graph) Get(key interface{}) *Node {
    return g.Nodes[key]
}

func (g *Graph) Add(node *Node, eager bool) {
    
    if eager {
        g.eager[node.NameKey] = node
    }
    if g.Nodes[node.Key] != nil {
        log.Debugf("Redefining:\n\t %s with \n\t %s", g.Nodes[node.Key], node)
        delete(g.Nodes, node.Key)
        g.Nodes[node.Key] = node
    } else {
        g.Nodes[node.Key] = node
    }
    if g.NamedNodes[node.NameKey] != nil {
        log.Debugf("Redefining:\n\t %s with \n\t %s", g.Nodes[node.NameKey], node)
        delete(g.NamedNodes, node.NameKey)
        g.NamedNodes[node.NameKey] = node
    } else {
        
        g.NamedNodes[node.NameKey] = node
    }
    
}

func (g *Graph) Connect(src, target interface{}) {

}

type Node struct {
    Key     interface{}
    NameKey string
    Value   *Instantiator
}

type Edge struct {
    Source Node
    Target Node
}


func NewGraph() *Graph {
    graph := new(Graph)
    graph.Nodes = make(map[interface{}]*Node, 0)
    graph.NamedNodes = make(map[string]*Node, 0)
    graph.eager = make(map[string]*Node, 0)
    return graph
}

