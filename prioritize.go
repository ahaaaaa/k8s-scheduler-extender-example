package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"strconv"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error)
}

func (p Prioritize) Handler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	return p.Func(*args.Pod, args.Nodes.Items)
}



func StealTimePreemption(promethuesUrl string, nodes []v1.Node) map[string] int {
	var nodeResult = getNodeStealTimeMetrics(promethuesUrl)
	var nodeStealTimeMap = make(map[string] float64)
	var pairList PairList
	for _, result := range nodeResult.Data.Result{
		hostname := result.Metric.Hostname
		stealtime, err := strconv.ParseFloat(result.Value[1].(string), 64)

		if( err != nil){
			fmt.Println("Hostname %v, err", hostname, err)
			panic(err)
		}


		if(len(hostname) == 0){
			panic("hostname is null")
		}

		nodeStealTimeMap[hostname] = stealtime

		pairList = SortMapByValueASC(nodeStealTimeMap)
	}

	var nodeStealTimeScoreMap = make(map[string] int)
	i := 0
	for _,p := range pairList{
		nodeStealTimeScoreMap [p.Key] = i
		i++
		fmt.Println("key : %v, v:%v", p.Key, p.Value)
	}
	return nodeStealTimeScoreMap
}
