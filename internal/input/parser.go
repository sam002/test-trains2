package input

import (
	"fmt"
	"strconv"
	"strings"
	"trains2/internal/trains2"
)

func ParseInput() (*[]trains2.Station, *[]trains2.Edge, *[]trains2.Package, *[]trains2.Train) {
	numStation := 0
	fmt.Println("Enter num of station:")
	if _, err := fmt.Scanln(&numStation); err != nil {
		panic(err)
	}

	fmt.Println("Stations names (each on new line):")
	stations := make([]trains2.Station, numStation)
	ms := make(map[string]trains2.Station, numStation)
	for i := 0; i < numStation; i++ {
		s := ""
		if _, err := fmt.Scanln(&s); err != nil {
			panic(err)
		}
		stations[i] = trains2.Station(s)
		ms[s] = trains2.Station(s)
	}
	fmt.Printf("Stations (%d): %v\n", numStation, stations)

	numEdges := 0
	fmt.Println("Enter num of edges:")
	if _, err := fmt.Scanln(&numEdges); err != nil {
		panic(err)
	}
	//  2 // number of edges
	edges := make([]trains2.Edge, numEdges)

	fmt.Println("Stations names (each on new line):")
	for i := range edges {
		s := ""
		if _, err := fmt.Scanln(&s); err != nil {
			panic(err)
		}
		r := strings.Split(s, ",")
		edges[i].Name = r[0]
		if sn, ok := ms[r[1]]; !ok {
			panic("Not parse edge: " + s)
		} else {
			edges[i].Node1 = sn
		}
		if sn, ok := ms[r[2]]; !ok {
			panic("Not parse edge: " + s)
		} else {
			edges[i].Node2 = sn
		}
		if l, err := strconv.Atoi(r[3]); err != nil {
			panic(err)
		} else {
			edges[i].JourneyTimeInMinutes = l
		}
	}
	//  E1,A,B,30 // route from A to B that takes 30 minutes
	//  E2,B,C,10 // route from B to C that takes 10 minutes
	fmt.Printf("Edges (%d): %v\n", numEdges, edges)

	numPacks := 0
	//  1 // number of deliveries to be performed
	if _, err := fmt.Scanln(&numPacks); err != nil {
		panic(err)
	}
	packages := make([]trains2.Package, numPacks)
	for i := range packages {
		//  K1,5,A,C // package K1 with weight 5 located currently at station A that must be delivered to station C
		s := ""
		if _, err := fmt.Scanln(&s); err != nil {
			panic(err)
		}
		r := strings.Split(s, ",")
		packages[i].PackageName = r[0]
		if l, err := strconv.Atoi(r[1]); err != nil {
			panic(err)
		} else {
			packages[i].WeightInKg = l
		}
		if sn, ok := ms[r[2]]; !ok {
			panic("Not parse package: " + s)
		} else {
			packages[i].StartingNode = &sn
		}
		if sn, ok := ms[r[3]]; !ok {
			panic("Not parse package: " + s)
		} else {
			packages[i].DestinationNode = &sn
		}
	}
	fmt.Printf("Package (%d): %v\n", numPacks, packages)

	numTrains := 0
	//  1 // number of trains
	if _, err := fmt.Scanln(&numTrains); err != nil {
		panic(err)
	}
	trains := make([]trains2.Train, numTrains)
	for i := range packages {
		//  Q1,6,B // train Q1 with capacity 6 located at station B
		s := ""
		if _, err := fmt.Scanln(&s); err != nil {
			panic(err)
		}
		r := strings.Split(s, ",")
		trains[i].TrainName = r[0]
		if l, err := strconv.Atoi(r[1]); err != nil {
			panic(err)
		} else {
			trains[i].CapacityInKg = l
		}
		if sn, ok := ms[r[2]]; !ok {
			panic("Not parse train: " + s)
		} else {
			trains[i].StartingNode = sn
		}
	}
	fmt.Printf("Train (%d): %v\n", numTrains, trains)
	return &stations, &edges, &packages, &trains
}
