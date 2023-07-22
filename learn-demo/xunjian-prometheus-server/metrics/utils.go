package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"xunjian-prometheus-server/tool"
)

type MetricResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type MetricData struct {
	ResultType string         `json:"resultType"`
	Result     []MetricResult `json:"result"`
}

type Metrics struct {
	Status string     `json:"status"`
	Data   MetricData `json:"data"`
}

func vectorMetricResult(promql_url string) string {
	// 发送请求
	resp, err := http.Get(promql_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 查询结果反序列化
	var metric Metrics
	err = json.Unmarshal(body, &metric)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(metric)

	// 判断是否有结果
	if len(metric.Data.Result) == 0 {
		return ""
	}
	re, ok := metric.Data.Result[0].Value[1].(string)
	if ok {
		return re
	}
	return ""
}

func getNodeCpuPercent(NodeName string) float64 {
	//0*node_uname_info+on(instance)group_left()(1-sum(increase(node_cpu_seconds_total{mode="idle"}[5m]))by(instance)/sum(increase(node_cpu_seconds_total[5m]))by(instance))*100
	// http://192.168.0.10:30003/api/v1/query?query=0*node_uname_info+on(instance)group_left()(1-sum(increase(node_cpu_seconds_total{mode="idle"}[5m]))by(instance)/sum(increase(node_cpu_seconds_total[5m]))by(instance))*100&time=1662723512.282
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] +
		"/api/v1/query?query=0*node_uname_info%7Bnodename%3D%22" + NodeName +
		"%22%7D%2Bon%28instance%29group_left%28%29%281-sum%28increase%28node_cpu_seconds_total%7Bmode%3D%22idle%22%7D%5B5m%5D%29%29by%28instance%29%2Fsum%28increase%28node_cpu_seconds_total%5B5m%5D%29%29by%28instance%29%29*100&time=" +
		fmt.Sprintf("%d", time.Now().Unix())

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		f64, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			log.Fatal(err)
		}
		return f64
	}
	return 0.0
}

func getNodeMemoryPercent(NodeName string) float64 {
	// 0*node_uname_info+on(instance)group_left()(1-node_memory_MemAvailable_bytes/node_memory_MemTotal_bytes)*100
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] + "/api/v1/query?" +
		"query=0*node_uname_info%7Bnodename%3D%22" + NodeName + "%22%7D%2Bon%28instance%29group_left%28%29%281-node_memory_MemAvailable_bytes%2Fnode_memory_MemTotal_bytes%29*100&time=" +
		fmt.Sprintf("%d", time.Now().Unix())

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		f64, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			log.Fatal(err)
		}
		return f64
	}
	return 0.0
}

func getNodeCpuTotal(NodeName string) int {
	// 0*node_uname_info{nodename="k8s-master-1"}+on(instance)group_left()count(node_cpu_seconds_total{mode='system'})by(instance)
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] + "/api/v1/query?" +
		"query=0*node_uname_info%7Bnodename%3D%22" + NodeName + "%22%7D%2Bon%28instance%29group_left%28%29count%28node_cpu_seconds_total%7Bmode%3D%27system%27%7D%29by%28instance%29&time=" +
		fmt.Sprintf("%d", time.Now().Unix())

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		i64, err := strconv.ParseInt(metric, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		return int(i64)
	}
	return 0
}

func getNodeLoad15(NodeName string) float64 {
	// 0*node_uname_info+on(instance)group_left()node_load15
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] + "/api/v1/query?" +
		"query=0*node_uname_info%7Bnodename%3D%22" + NodeName + "%22%7D%2Bon%28instance%29group_left%28%29node_load15&time=" +
		fmt.Sprintf("%d", time.Now().Unix())

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		f64, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			log.Fatal(err)
		}
		return f64
	}
	return 0.0
}

func getNodeDiskPercent(NodeName string) float64 {
	//0 * node_uname_info + on(instance) group_left() max((node_filesystem_size_bytes{fstype=~"ext.?|xfs"}-node_filesystem_free_bytes{fstype=~"ext.?|xfs"}) *100/(node_filesystem_avail_bytes {fstype=~"ext.?|xfs"}+(node_filesystem_size_bytes{fstype=~"ext.?|xfs"}-node_filesystem_free_bytes{fstype=~"ext.?|xfs"})))by(instance)
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] + "/api/v1/query?" +
		"query=0+*+node_uname_info%7Bnodename%3D%22" + NodeName + "%22%7D+%2B+on%28instance%29+group_left%28%29+max%28%28node_filesystem_size_bytes%7Bfstype%3D%7E%22ext.%3F%7Cxfs%22%7D-node_filesystem_free_bytes%7Bfstype%3D%7E%22ext.%3F%7Cxfs%22%7D%29+*100%2F%28node_filesystem_avail_bytes+%7Bfstype%3D%7E%22ext.%3F%7Cxfs%22%7D%2B%28node_filesystem_size_bytes%7Bfstype%3D%7E%22ext.%3F%7Cxfs%22%7D-node_filesystem_free_bytes%7Bfstype%3D%7E%22ext.%3F%7Cxfs%22%7D%29%29%29by%28instance%29&time=" +
		fmt.Sprintf("%d", time.Now().Unix())

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		f64, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			log.Fatal(err)
		}
		return f64
	}
	return 0.0
}

func getPodCpu(NameSpace, PodName string) float64 {
	// sum(irate(container_cpu_usage_seconds_total{container!="",container!="POD",pod="aaa",namespace="bbb"}[2m]))by(pod)
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] + "/api/v1/query?" +
		"query=sum%28irate%28container_cpu_usage_seconds_total%7Bcontainer%21%3D%22%22%2Ccontainer%21%3D%22POD%22%2Cpod%3D%22" + PodName + "%22%2Cnamespace%3D%22" + NameSpace + "%22%7D%5B2m%5D%29%29by%28pod%29&time=" +
		fmt.Sprintf("%d", time.Now().Unix())
	metric := vectorMetricResult(promql_url)
	if metric != "" {
		f64, err := strconv.ParseFloat(metric, 64)
		if err != nil {
			log.Fatal(err)
		}
		f64, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", f64), 64)
		return f64
	}
	return 0.0
}

func getPodMemory(NameSpace, PodName string) int64 {
	// promql: sum(container_memory_working_set_bytes{pod=~"monitor-grafana-99c94c4d4-ffm77",namespace=~"monitor",container !~ "POD",container !=""}) by(pod)
	// url: http://192.168.0.10:30003/api/v1/query?query=sum(container_memory_working_set_bytes{pod=~"monitor-grafana-99c94c4d4-ffm77",namespace=~"monitor",container+!~+"POD",container+!=""})+by(pod)&time=1662704823.883
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] +
		"/api/v1/query?query=sum(container_memory_working_set_bytes{pod=\"" + PodName + "\",namespace=\"" + NameSpace + "\",container+!=\"POD\",container+!=\"\"})+by(pod)&time=" + fmt.Sprintf("%d", time.Now().Unix())

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		i64, err := strconv.ParseInt(metric, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		return i64
	}
	return 0
}

func getPodRestartCount(NameSpace, PodName string) int64 {
	// sum(kube_pod_container_status_restarts_total{pod="monitor-grafana-6b47cf5798-qhs2b",namespace="monitor"})by(pod)
	promql_url := "http://" + tool.Conf.Prometheus["ip"] + ":" + tool.Conf.Prometheus["port"] + "/api/v1/query?" +
		"query=sum%28kube_pod_container_status_restarts_total%7Bpod%3D%22" + PodName + "%22%2Cnamespace%3D%22" + NameSpace + "%22%7D%29+by%28pod%29&time=1663078871.189"

	metric := vectorMetricResult(promql_url)
	if metric != "" {
		i64, err := strconv.ParseInt(metric, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		return i64
	}
	return 0
}
