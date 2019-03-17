package main

import (
	"./pathedge"
	"./permutations"

	"flag"
	"fmt"
	"time"
)

func main() {

	displayResCount := flag.Int("display", 3, "the number of displayed results")

	flag.Parse()

	for p := range permutations.PairsGenerator(permutations.KeysFromGraph(pathedge.RelatedGraph)) {
		var results []pathedge.PathEdge
		worklist := make(chan pathedge.PathEdge, 10)
		startpointId, endpointId := p[0].(int), p[1].(int)

		startpoint := pathedge.NewPathEdge(pathedge.OutCome{
			time.Duration(0), 0},
			[]int{startpointId}, make(map[int]pathedge.TrainLeg))

		startpoint.Path(worklist, &pathedge.RelatedGraph, endpointId)

		for p := range worklist {
			results = append(results, p)
		}
		fmt.Printf("############   StartId: %v  ####  FinishId: %v   ###############\n", startpointId, endpointId)
		pathedge.ShowResult(results, *displayResCount)
	}

}
