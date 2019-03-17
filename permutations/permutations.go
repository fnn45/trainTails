package permutations

import "../pathedge"

func PairsGenerator(s []interface{}) chan []interface{} {

	ch := make(chan []interface{})
	go func() {
		defer close(ch)
		for i := 0; i < len(s); i++ {
			for j := i + 1; j < len(s); j++ {
				ch <- []interface{}{s[i], s[j]}
				ch <- []interface{}{s[j], s[i]}
			}
		}
	}()
	return ch
}

func KeysFromGraph(g pathedge.Graph) (res []interface{}) {

	for k, _ := range g {
		res = append(res, k)
	}
	return
}
