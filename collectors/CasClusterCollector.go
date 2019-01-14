package collectors
//
//import (
//	"github.com/prometheus/client_golang/prometheus"
//	"virtual-exporter/config"
//	"fmt"
//	"encoding/xml"
//	"strconv"
//	"log"
//)
//
//type CasClusterCollector struct {
//	Target string
//}
//
//func (c CasClusterCollector) Describe(ch chan<- *prometheus.Desc) {
//	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
//}
//
//var (
//	cas_cluster_disk_free = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cluster_disk_free"), "cas cluster disk free",
//		[]string{"cas_ip", "cas_host_pool_id", "cas_cluster_id",}, nil)
//
//	cas_cluster_cpu_usage = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cluster_cpu_usage"), "cas cluster cpu usage",
//		[]string{"cas_ip", "cas_host_pool_id", "cas_cluster_id",}, nil)
//
//	cas_cluster_memory_usage = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cluster_memory_usage"), "cas cluster memory usage",
//		[]string{"cas_ip", "cas_host_pool_id", "cas_cluster_id",}, nil)
//
//	cas_cluster_disk_usage = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cluster_disk_usage"), "cas cluster disk usage",
//		[]string{"cas_ip", "cas_host_pool_id", "cas_cluster_id",}, nil)
//)
//
//func (c CasClusterCollector) Collect(ch chan<- prometheus.Metric) {
//
//	monitor_info := config.GetMonitorInfo(c.Target)
//	hostPoolId := monitor_info.Params_maps["hostPoolId"]
//	clusterId := monitor_info.Params_maps["clusterId"]
//	ip := monitor_info.Params_maps["ip"]
//	port := monitor_info.Params_maps["port"]
//	username := monitor_info.Params_maps["username"]
//	password := monitor_info.Params_maps["password"]
//	errs, ret := GetAnotherData(fmt.Sprintf("/cas/casrs/cluster/summary/%s", clusterId), ip, username, password, port)
//	casClusterDetil := &config.CasClusterDetails{}
//	if errs != nil {
//		log.Printf("get cluster detail error:%s", errs.Error())
//		return
//	}
//	errs = xml.Unmarshal(ret, &casClusterDetil)
//
//	labelvalues := []string{ip, hostPoolId, clusterId}
//	for _, clusterInfo := range (*casClusterDetil).CasClusterInfo {
//		if clusterInfo.Key == "本地总可用存储" {
//			valStr := config.Substr(clusterInfo.Value, 0, len(clusterInfo.Value)-2)
//			valUnit := config.Substr(clusterInfo.Value, len(clusterInfo.Value)-2, 2)
//			value, err := config.ValToGB(valStr, valUnit)
//			if err != nil {
//				log.Printf("parse cluster disk free:%s", err)
//			} else {
//				ch <- prometheus.MustNewConstMetric(cas_cluster_disk_free, prometheus.GaugeValue, value, labelvalues...) //MB
//			}
//		}
//		if clusterInfo.Key == "cpuRate" {
//			val, err := strconv.ParseFloat(clusterInfo.Value, 64)
//			if err != nil {
//				log.Printf("parse cluster cpu rate:%s", err)
//			} else {
//				ch <- prometheus.MustNewConstMetric(cas_cluster_cpu_usage, prometheus.GaugeValue, val, labelvalues...)
//			}
//		}
//		if clusterInfo.Key == "memRate" {
//			val, err := strconv.ParseFloat(clusterInfo.Value, 64)
//			if err != nil {
//				log.Printf("parse cluster mem rate:%s", err)
//			} else {
//				ch <- prometheus.MustNewConstMetric(cas_cluster_memory_usage, prometheus.GaugeValue, val, labelvalues...)
//			}
//		}
//
//		if clusterInfo.Key == "diskRate" {
//			val, err := strconv.ParseFloat(clusterInfo.Value, 64)
//			if err != nil {
//				log.Printf("parse cluster disk rate:%s", err)
//			} else {
//				ch <- prometheus.MustNewConstMetric(cas_cluster_disk_usage, prometheus.GaugeValue, val, labelvalues...)
//			}
//		}
//
//	}
//}
