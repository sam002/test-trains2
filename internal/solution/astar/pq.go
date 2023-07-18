package astar

type Item struct {
	n        Node
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func NewItem(n Node) *Item {
	priority := n.gm.Game.summaryDistance() * n.gm.Game.time
	return &Item{n: n, priority: priority}
}

func (i *Item) GetNode() Node {
	return i.n
}

// A PQMovies implements heap.Interface and holds Items.
type PQMovies []*Item

func (pq PQMovies) Len() int { return len(pq) }

func (pq PQMovies) Less(i, j int) bool {
	// MinPQ
	return pq[i].priority < pq[j].priority
}

func (pq PQMovies) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PQMovies) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PQMovies) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
