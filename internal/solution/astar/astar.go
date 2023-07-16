package astar

import (
	"container/heap"
	"fmt"
	"gonum.org/v1/gonum/graph"
	_ "gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"math"
	"time"
	"trains2/internal/solution"
	"trains2/internal/trains2"
)

var _ solution.Solution = &Solution{}

type Solution struct {
	stations  map[trains2.Station]graph.Node
	rStations map[graph.Node]trains2.Station
	trains    *[]trains2.Train
	packages  *[]trains2.Package
	roadGraph *simple.WeightedUndirectedGraph
	//deliveryGraph *simple.WeightedDirectedGraph
	steps    []trains2.Move
	distance map[string]int
	time     int
}

func (s *Solution) summaryDistance() int {
	sum := 0
	for _, d := range s.distance {
		sum += d
	}

	return sum
}

func summaryDistanceBy(distance map[string]int) int {
	sum := 0
	for _, d := range distance {
		sum += d
	}

	return sum
}

func (s *Solution) isSolved() bool {
	return s.summaryDistance() == 0
}

func (s *Solution) calcDistance(move trains2.Move) map[string]int {
	lc := s.distance
	for _, p := range move.PackagesPickedUp {
		lc[p.PackageName] = int(path.DijkstraFrom(s.stations[move.From], s.roadGraph).
			WeightTo(s.stations[move.To].ID()))
	}

	return lc
}

func NewAStarSolution(stations *[]trains2.Station, edges *[]trains2.Edge, packages *[]trains2.Package, trains *[]trains2.Train) *Solution {
	sol := Solution{
		stations:  map[trains2.Station]graph.Node{},
		rStations: map[graph.Node]trains2.Station{},
		trains:    trains,
		packages:  packages,
		roadGraph: simple.NewWeightedUndirectedGraph(0, math.Inf(1)),
		//deliveryGraph: simple.NewWeightedDirectedGraph(0, math.Inf(1)),
		steps:    nil,
		distance: map[string]int{},
	}
	for i, s := range *stations {

		sol.stations[s] = simple.Node(int64(i))

		sol.rStations[sol.stations[s]] = s
	}
	for _, e := range *edges {
		we := simple.WeightedEdge{
			F: sol.stations[e.Node1],
			T: sol.stations[e.Node2],
			W: float64(e.JourneyTimeInMinutes),
		}
		sol.roadGraph.SetWeightedEdge(we)
	}
	for _, p := range *packages {
		sol.distance[p.PackageName] = int(path.DijkstraFrom(sol.stations[*p.StartingNode], sol.roadGraph).
			WeightTo(sol.stations[*p.DestinationNode].ID()))
		//we := simple.WeightedEdge{
		//	F: sol.stations[string(p.StartingNode)],
		//	T: sol.stations[string(p.DestinationNode)],
		//	W: float64(p.WeightInKg),
		//}
		//sol.deliveryGraph.SetWeightedEdge(we)
	}

	return &sol
}

func (s *Solution) Calculate() error {
	//r := dynamic.NewDStarLite(s.roadGraph)
	pq := make(PQMovies, 0)
	//heap.Init(&pq)
	prevMovies := make(map[string]int, 0)
	for !s.isSolved() {
		time.Sleep(2 * time.Second)
		nm := s.nextMovies()
		for _, m := range nm {
			fmt.Println(m.String())
			i := &Item{
				value: m,
			}
			heap.Push(&pq, i)

			for ok := true; ok; _, ok = prevMovies[m.String()] {
				fmt.Println(m)
				m = heap.Pop(&pq).(*Item).value
				s.distance = s.calcDistance(m)
			}
		}
	}

	for pq.Len() > 0 {
		m := heap.Pop(&pq).(*Item).value
		s.steps = append(s.steps, m)
	}

	return nil
}

func (s *Solution) nextMovies() []trains2.Move {
	res := make([]trains2.Move, 0)

	for _, t := range *s.trains {
		for _, ns := range graph.NodesOf(s.roadGraph.From(s.stations[t.StartingNode].ID())) {
			nm := trains2.Move{
				When:               s.time,
				Train:              t,
				From:               t.StartingNode,
				To:                 s.rStations[ns],
				PackagesPickedUp:   nil,
				PackagesDroppedOff: nil,
			}
			//todo add capacity optimization
			for _, p := range *s.packages {
				if *p.StartingNode != t.StartingNode || p.StartingNode == p.DestinationNode {
					continue
				}

				//pick up only when a passing direction
				nd := path.DijkstraFrom(s.stations[t.StartingNode], s.roadGraph).
					WeightTo(s.stations[*p.DestinationNode].ID())
				if int(nd) < s.distance[p.PackageName] {
					continue
				}

				nm.PackagesPickedUp = append(nm.PackagesPickedUp, p)
				p.PickUp(&t)
			}
			for _, p := range t.CurrentPackages {
				//Local optimization - drop off all of packages
				p.DropOff(p.DestinationNode)
			}

			res = append(res, nm)
		}
	}

	return res
}

func (s *Solution) GetSteps() *[]trains2.Move {
	return &s.steps
}
