package pathedge

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/xml"
)

type Graph map[int][]TrainLeg

var RelatedGraph Graph

type TrainLegs struct {
	XMLname   xml.Name   `xml:"TrainLegs"`
	TrainLegs []TrainLeg `xml:"TrainLeg"`
}

type TrainLeg struct {
	TrainId             int     `xml:"TrainId,attr"`
	DepartureStationId  int     `xml:"DepartureStationId,attr"`
	ArrivalStationId    int     `xml:"ArrivalStationId,attr"`
	Price               float32 `xml:"Price,attr"`
	ArrivalTimeString   string  `xml:"ArrivalTimeString,attr"`
	DepartureTimeString string  `xml:"DepartureTimeString,attr"`
}

func init() {

	xmlFile, err := os.Open("data.xml")
	if err != nil {
		log.Fatal("Error: %v ", err)
	}
	defer xmlFile.Close()
	byteVal, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	var legs TrainLegs
	if err := xml.Unmarshal(byteVal, &legs); err != nil {
		log.Fatal(err)
	}

	RelatedGraph = make(Graph)

	for _, l := range legs.TrainLegs {
		RelatedGraph[l.DepartureStationId] = append(RelatedGraph[l.DepartureStationId], l)
	}
}
