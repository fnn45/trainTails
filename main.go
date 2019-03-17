package main

import (
	"./pathedge"
	"flag"
	"log"
	"time"
)

func main() {

	startpointId := flag.Int("start", 1902, "starting point")
	endpointId := flag.Int("end", 1902, "end point")
	displayResCount := flag.Int("display", 3, "the number of displayed results")

	flag.Parse()

	var results []pathedge.PathEdge
	worklist := make(chan pathedge.PathEdge, 10)

	if _, ok := pathedge.RelatedGraph[*startpointId]; !ok {
		log.Fatal("Enter the correct starting id")
	}
	if _, ok := pathedge.RelatedGraph[*endpointId]; !ok {
		log.Fatal("Enter the correct finish id")
	}

	startpoint := pathedge.NewPathEdge(pathedge.OutCome{
		time.Duration(0), 0},
		[]int{*startpointId}, make(map[int]pathedge.TrainLeg))

	startpoint.Path(worklist, &pathedge.RelatedGraph, *endpointId)

	for p := range worklist {
		results = append(results, p)
	}
	pathedge.ShowResult(results, *displayResCount)
}
