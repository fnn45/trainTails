package pathedge

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OutCome struct {
	Duration time.Duration
	Price    float32
}

type PathEdge struct {
	OutCome
	keys  []int
	nodes map[int]TrainLeg
}

type DurationSortablePathSet []PathEdge

type PriceSortablePathSet []PathEdge

func (p *PathEdge) Path(worklist chan PathEdge, g *Graph, endpoint int) {
	edges := (*g)[p.keys[len(p.keys)-1]]

	go func() {
		isway := false
		chans := make([]chan PathEdge, 10)
		for _, edge := range edges {
			if !elemInArray(edge.ArrivalStationId, p.keys) {
				isway = true
				newQ := p.Copy()
				newQ.nodes[len(p.keys)-1] = edge
				newQ.keys = append(p.keys, edge.ArrivalStationId)
				newQ.Price += edge.Price
				newQ.Duration += CalculateDuration(edge.DepartureTimeString, edge.ArrivalTimeString)
				if endpoint == edge.ArrivalStationId {
					worklist <- newQ
				} else {
					ch := make(chan PathEdge)
					chans = append(chans, ch)
					newQ.Path(ch, g, endpoint)
					for c := range ch {
						worklist <- c
					}
				}
			}
		}

		close(worklist)
		if !isway {
			for _, ch := range chans {
				if ch != nil {
					close(ch)
				}
			}

		}
	}()
}

func elemInArray(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func NewPathEdge(out OutCome, keys []int, nodes map[int]TrainLeg) PathEdge {
	return PathEdge{OutCome: out, keys: keys, nodes: nodes}
}

func (q *PathEdge) Copy() PathEdge {
	newNodes := make(map[int]TrainLeg, len(q.keys))
	newKeys := make([]int, len(q.keys))
	copy(newKeys, q.keys)
	for k, v := range q.nodes {
		newNodes[k] = v
	}
	return PathEdge{OutCome: OutCome{Duration: q.Duration, Price: q.Price}, keys: newKeys, nodes: newNodes}
}

func CalculateDuration(start, finish string) time.Duration {
	time1, _ := ParseTime(start)
	time2, _ := ParseTime(finish)
	diff := time1 - time2
	if diff < 0 {
		diff = diff + time.Duration(24*time.Hour)
	}
	return diff
}

func ParseTime(t string) (time.Duration, error) {
	var secs, mins, hours int
	var err error

	parts := strings.SplitN(t, ":", 3)

	hours, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	mins, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	secs, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	if secs > 59 || secs < 0 || mins > 59 || mins < 0 || hours > 23 || hours < 0 {
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	return time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute + time.Duration(secs)*time.Second, nil
}

func (d DurationSortablePathSet) Len() int {
	return len(d)
}

func (d DurationSortablePathSet) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DurationSortablePathSet) Less(i, j int) bool {
	return d[i].Duration < d[j].Duration
}

func (p PriceSortablePathSet) Len() int {
	return len(p)
}

func (p PriceSortablePathSet) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PriceSortablePathSet) Less(i, j int) bool {
	return p[i].Price < p[j].Price
}

func (d DurationSortablePathSet) String() string {
	var buffer bytes.Buffer
	splitLine := "-----------------------------------------------------\n"
	data := "TrainId: %v (DepId: %v, ArrId: %v)\n"
	buffer.WriteString("\n")
	buffer.WriteString("<<< --------- sorted by duration  --------->>>\n\n")
	for i := 0; i <= len(d)-1; i++ {
		buffer.WriteString(
			fmt.Sprintf("duration:  %v\n", d[i].Duration.String()))
		nodes := d[i].nodes
		for j := 0; j <= len(nodes)-1; j++ {
			node := nodes[j]
			buffer.WriteString(
				fmt.Sprintf(data, node.TrainId, node.DepartureStationId, node.ArrivalStationId))
		}
		buffer.WriteString(splitLine)
	}
	return buffer.String()
}

func (p PriceSortablePathSet) String() string {
	var buffer bytes.Buffer
	splitLine := "-----------------------------------------------------\n"
	data := "TrainId: %v (DepId: %v, ArrId: %v)\n"
	buffer.WriteString("\n")
	buffer.WriteString("<<< --------- sorted by price  --------->>>\n\n")
	for i := 0; i <= len(p)-1; i++ {
		buffer.WriteString(
			fmt.Sprintf("price:  %v\n", p[i].Price))
		nodes := p[i].nodes
		for j := 0; j <= len(nodes)-1; j++ {
			node := nodes[j]
			buffer.WriteString(
				fmt.Sprintf(data, node.TrainId, node.DepartureStationId, node.ArrivalStationId))
		}
		buffer.WriteString(splitLine)
	}
	return buffer.String()
}

func ShowResult(qs []PathEdge, rescount int) {
	durres := make(chan DurationSortablePathSet)
	prres := make(chan PriceSortablePathSet)
	defer close(durres)
	defer close(prres)
	go func() {
		arr := make([]PathEdge, rescount)
		copy(arr, qs)
		sort.Sort(DurationSortablePathSet(arr))
		durres <- arr
	}()
	go func() {
		arr := make([]PathEdge, rescount)
		copy(arr, qs)
		sort.Sort(PriceSortablePathSet(arr))
		prres <- arr
	}()
	fmt.Println(<-durres)
	fmt.Println(<-prres)
}
