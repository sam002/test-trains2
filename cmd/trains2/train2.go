package main

import (
	"fmt"
	"trains2/internal/input"
	"trains2/internal/solution/astar"
)

func main() {
	fmt.Println("START TRAINS2")
	stations, edges, packages, trains := input.ParseInput()
	solution := astar.NewAStarSolution(stations, edges, packages, trains)
	err := solution.Calculate()
	if err != nil {
		panic(err)
	}
	for _, s := range solution.GetSteps() {
		fmt.Println(s.String())
	}

	fmt.Println("END TRAINS2")
}
