package trains2

import (
	"fmt"
)

type Station string

type Edge struct {
	Name                 string
	Node1                Station
	Node2                Station
	JourneyTimeInMinutes int
}

func (e Edge) String() string {
	return fmt.Sprintf("{%s, %s<->%s, %d(min)}", e.Name, e.Node1, e.Node2, e.JourneyTimeInMinutes)
}

type Train struct {
	TrainName    string
	CapacityInKg int
	StartingNode Station
}

//func (t Train) String() string {
//	return fmt.Sprintf("{%s, %s<->%s, %d(min)}", e.Name, e.Node1.Name, e.Node2.Name, e.JourneyTimeInMinutes)
//}

type Package struct {
	PackageName     string
	WeightInKg      int
	StartingNode    Station
	DestinationNode Station
}

func (p *Package) String() string {
	return fmt.Sprintf("%s %d %s %s", p.PackageName, p.WeightInKg, p.StartingNode, p.DestinationNode)
}

type Move struct {
	When               int
	Train              Train
	From               Station
	To                 Station
	PackagesPickedUp   []Package
	PackagesDroppedOff []Package
}

func (m *Move) String() string {
	return fmt.Sprintf("W=%d, T=%s, N1=%s, P1=%v, N2=%s, P2=%v",
		m.When, m.Train.TrainName, m.From, m.PackagesPickedUp, m.To, m.PackagesDroppedOff)
}
