package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strings"
	"fmt"
	"github.com/gorilla/mux"
	"virtual-exporter/config"
	"virtual-exporter/collectors"
	"virtual-exporter/collectors/api"
)
var listenAddress = kingpin.Flag("web.listen-address","Address to listen on for web " +
	"interface and telemetry.").Default(":9107").String()



func runCollector(collector prometheus.Collector,w http.ResponseWriter,r *http.Request)  {
	registry:= prometheus.NewRegistry()
	registry.MustRegister(collector)
	h:=promhttp.HandlerFor(registry,promhttp.HandlerOpts{})
	h.ServeHTTP(w,r)
}
func init() {
	config.GetDBHandle()
}
func main() {
	kingpin.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/v1/cas/resourceList",handleOther)
	r.HandleFunc("/api/v1/mysql/access",api.CasAccess).Methods(http.MethodPost)
	r.HandleFunc("/v1/cas/clusterList",api.GetAllCluster)//没啥用

	r.HandleFunc("/cas",handler)
	r.HandleFunc("/cvk",handler)
	r.HandleFunc("/vm",handler)
	http.ListenAndServe(*listenAddress,r)

}
func handleOther(w http.ResponseWriter, r *http.Request) {
	collectors.GetCasInfo(w,r)
}
func handler(w http.ResponseWriter,r *http.Request)  {
	var collectorType prometheus.Collector
	target:= r.URL.Query().Get("target")
	if target=="" {
		http.Error(w,"'target' parameter must be specified",400)
		return
	}
	switch strings.Split(fmt.Sprintf("%s",r.URL),"?")[0] {
	case "/cas":
		collectorType = collectors.CasCollector{target}
		break
	case "/cvk":
		collectorType = collectors.CvkCollector{target}
		break
	case "/vm":
		collectorType = collectors.VmCollector{target}
		break
	default:
		break
	}

	runCollector(collectorType,w,r)
}