package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type NodeStealTimePercentage struct {
	Status string `json : "status"`
	Data NodeData `json : "data"`

}

type NodeData struct {
	ResultType string `json : "resultType"`
	Result [] NodeDataResult `json : "result"`
}

type NodeDataResult struct {
	Metric NodeDataResultMetric `json : "metric"`
	Value [] interface{} `json : value`
}
type NodeDataResultMetric struct{
	Name string `json:"__name__"`
	Hostname string `json : "hostname"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}


func getNodeStealTimeMetrics(prometheusUrl string) NodeStealTimePercentage{
	resp, err := httpClient.Get(prometheusUrl + "query=record_node_steal_time_percentage{hostname=~\".*-k8snode-ad.*|.*-k8snode-ingress-internal.*|.*-k8snode-ingress-external.*\"}")
	if err != nil {
		print("getNodeStealTimeMetrics error, ", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
	}

	var data NodeStealTimePercentage
	json.Unmarshal(body, &data)
	fmt.Println("Result : %v\n", data)
	return data
}




func main() {
	var nodeResult = getNodeStealTimeMetrics("http://prometheus.nj.iapp.jpushoa.com/api/v1/query?")
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

	for _,p := range pairList{
		fmt.Println("key : %v, v:%v", p.Key, p.Value)
	}

	fmt.Println(nodeStealTimeMap)
}
