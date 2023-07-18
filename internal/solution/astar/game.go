package astar

import (
	"fmt"
	"golang.org/x/exp/maps"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"math"
	"trains2/internal/trains2"
)

type GameMap struct {
	stations  map[trains2.Station]graph.Node
	rStations map[graph.Node]trains2.Station
	roadGraph *simple.WeightedUndirectedGraph
}

func NewGameMap(stations *[]trains2.Station, edges *[]trains2.Edge) *GameMap {
	gm := GameMap{
		stations:  make(map[trains2.Station]graph.Node),
		rStations: make(map[graph.Node]trains2.Station),
		roadGraph: simple.NewWeightedUndirectedGraph(0, math.Inf(1)),
	}
	for i, s := range *stations {
		gm.stations[s] = simple.Node(int64(i))
		gm.rStations[gm.stations[s]] = s
	}

	for _, e := range *edges {
		we := simple.WeightedEdge{
			F: gm.stations[e.Node1],
			T: gm.stations[e.Node2],
			W: float64(e.JourneyTimeInMinutes),
		}
		gm.roadGraph.SetWeightedEdge(we)
	}

	return &gm
}

func (gm GameMap) timeBetween(s1 *trains2.Station, s2 *trains2.Station) int {
	return int(path.DijkstraFrom(gm.stations[*s1], gm.roadGraph).
		WeightTo(gm.stations[*s2].ID()))
}

func (gm GameMap) NeighborsTo(station trains2.Station) []trains2.Station {
	res := make([]trains2.Station, 0)
	for _, n := range graph.NodesOf(gm.roadGraph.From(gm.stations[station].ID())) {
		res = append(res, gm.rStations[n])
	}
	return res
}

type GameState struct {
	gameMap  GameMap
	trains   []trains2.Train
	packages []trains2.Package
	time     int
	distance map[string]int
}

func NewGameStateAuto(
	gameMap GameMap,
	trains []trains2.Train,
	packages []trains2.Package,
	time int,
) *GameState {
	gs := GameState{
		gameMap:  gameMap,
		trains:   make([]trains2.Train, len(trains)),
		packages: make([]trains2.Package, len(packages)),
		time:     time,
		distance: map[string]int{},
	}
	copy(gs.trains, trains)
	copy(gs.packages, packages)
	for _, p := range packages {
		gs.distance[p.PackageName] = int(path.DijkstraFrom(gs.gameMap.stations[p.StartingNode], gs.gameMap.roadGraph).
			WeightTo(gs.gameMap.stations[p.DestinationNode].ID()))
	}
	return &gs
}

func NewGameState(
	gameMap GameMap,
	trains []trains2.Train,
	packages []trains2.Package,
	time int,
	distance map[string]int,
) *GameState {
	gs := &GameState{
		gameMap:  gameMap,
		trains:   make([]trains2.Train, len(trains)),
		packages: make([]trains2.Package, len(packages)),
		time:     time,
		distance: make(map[string]int),
	}

	copy(gs.trains, trains)
	copy(gs.packages, packages)
	maps.Copy(gs.distance, distance)

	return gs
}

func (g *GameState) summaryDistance() int {
	sum := 0
	for _, d := range g.distance {
		sum += d
	}

	return sum
}

func (g *GameState) IsSolved() bool {
	return g.summaryDistance() == 0
}

type GameMove struct {
	Move *trains2.Move
	Game GameState
}

func (g *GameMove) StateID() string {
	return fmt.Sprintf("%v%v%v", g.Game.packages, g.Game.trains, g.Move)
}

func (g *GameState) NextMoves() []GameMove {
	res := make([]GameMove, 0)

	for ti, t := range g.trains {
		for _, ns := range g.gameMap.NeighborsTo(t.StartingNode) {
			time := g.gameMap.timeBetween(&t.StartingNode, &ns)

			gm := GameMove{
				Move: &trains2.Move{
					When:               g.time,
					Train:              t,
					From:               t.StartingNode,
					To:                 ns,
					PackagesPickedUp:   nil,
					PackagesDroppedOff: nil,
				},
				Game: *NewGameState(g.gameMap, g.trains, g.packages, g.time+time, g.distance),
			}

			gm.Game.trains[ti].StartingNode = ns

			//todo add capacity optimization
			for pi, p := range g.packages {
				if p.StartingNode != t.StartingNode || p.StartingNode == p.DestinationNode {
					continue
				}
				nd := path.DijkstraFrom(g.gameMap.stations[ns], g.gameMap.roadGraph).
					WeightTo(g.gameMap.stations[p.DestinationNode].ID())
				//local optimization, pick up only when a passing direction
				if int(nd) >= g.distance[p.PackageName] {
					continue
				}
				gm.Game.distance[p.PackageName] = int(nd)
				gm.Game.packages[pi].StartingNode = gm.Move.To

				gm.Move.PackagesPickedUp = append(gm.Move.PackagesPickedUp, p)
				gm.Move.PackagesDroppedOff = gm.Move.PackagesPickedUp
			}

			res = append(res, gm)
		}
	}

	return res
}
