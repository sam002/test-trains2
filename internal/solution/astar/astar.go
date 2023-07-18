package astar

import (
	"container/heap"
	_ "gonum.org/v1/gonum/graph/iterator"
	"trains2/internal/solution"
	"trains2/internal/trains2"
)

var _ solution.Solution = &Solution{}

type Solution struct {
	//deliveryGraph *simple.WeightedDirectedGraph
	steps []GameMove
	time  int
}

//func (s *Solution) calcDistance(move trains2.Move) map[string]int {
//	lc := s.distance
//	for _, p := range move.PackagesPickedUp {
//		lc[p.PackageName] = int(path.DijkstraFrom(s.stations[move.From], s.roadGraph).
//			WeightTo(s.stations[move.To].ID()))
//	}
//
//	return lc
//}

func NewAStarSolution(stations *[]trains2.Station, edges *[]trains2.Edge, packages *[]trains2.Package, trains *[]trains2.Train) *Solution {
	sol := Solution{
		steps: nil,
	}
	gm := NewGameMap(stations, edges)
	gs := NewGameStateAuto(*gm, *trains, *packages, 0)

	sol.steps = append(sol.steps, GameMove{
		Move: nil,
		Game: *gs,
	})
	return &sol
}

type Node struct {
	prev *Node
	gm   GameMove
	time int
}

func (s *Solution) Calculate() error {
	possibleMoves := make(PQMovies, 0)
	heap.Init(&possibleMoves)
	visited := make(map[string]Node, 0)
	n := Node{
		prev: nil,
		gm:   s.steps[0],
		time: 0,
	}
	//fmt.Println("Add to visited:")
	//fmt.Println(n.gm.StateID())
	visited[n.gm.StateID()] = n

	for !n.gm.Game.IsSolved() {
		//time.Sleep(1 * time.Second)
		nm := n.gm.Game.NextMoves()
		for _, m := range nm {

			prev := n
			if n.prev == nil || n.prev.gm.StateID() != m.StateID() {
				//fmt.Println("add to PQ:")
				//fmt.Println(m.StateID())
				heap.Push(&possibleMoves, NewItem(Node{
					prev: &prev,
					gm:   m,
					time: m.Move.When,
				}))
			}
		}

		_, ok := visited[n.gm.StateID()]
		for ok {
			if possibleMoves.Len() == 0 {
				panic("SOLUTION IMPOSSIBLE")
			}
			n = heap.Pop(&possibleMoves).(*Item).GetNode()

			_, ok = visited[n.gm.StateID()]
			//fmt.Println("Pop from queue:")
			//fmt.Println(n.gm.StateID(), ok)
		}

		//fmt.Println("Add to visited:")
		//fmt.Println(n.gm.StateID())
		visited[n.gm.StateID()] = n
	}

	for n.prev != nil {
		s.steps = append(s.steps, n.gm)
		n = *n.prev
	}

	return nil
}

func (s *Solution) GetSteps() []trains2.Move {
	res := make([]trains2.Move, 0)
	for i := len(s.steps) - 1; i > 0; i-- {
		if s.steps[i].Move == nil {
			continue
		}
		res = append(res, *s.steps[i].Move)
	}

	return res
}
