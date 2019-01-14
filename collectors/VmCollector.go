package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
	"virtual-exporter/config"
	"encoding/xml"
	"strconv"
	"log"
)

type VmCollector struct {
	Target string
}
func (c VmCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}
var (
	/*cas_vm_cpu = prometheus.NewDesc("cas_vm_cpu","cas virtual machine cpu",
		[]string{"cas_vm_id","cas_cvk_host_id",},nil)
	cas_vm_memory = prometheus.NewDesc("cas_vm_memory","cas virtual machine memory",
		[]string{"cas_vm_id","cas_cvk_host_id",},nil)
	casVmMemoryFree = prometheus.NewDesc("cas_vm_memory_free","cas virtual machine memory free",
		[]string{"cas_vm_id","cas_cvk_host_id",},nil)
	casVmMemoryUsed = prometheus.NewDesc("cas_vm_memory_used","cas virtual machine memory used",
		[]string{"cas_vm_id","cas_cvk_host_id",},nil)*/
	cas_vm_monitorstatus = prometheus.NewDesc(prometheus.BuildFQName("cas","","vm_monitorstatus"),
		"cas vm monitorstatus",nil,nil)
	cas_vm_cpu_usage = prometheus.NewDesc("cas_vm_cpu_usage","cas virtual machine cpu usage",
		[]string{"cas_vm_id","cas_cvk_host_id",},nil)
	cas_vm_memory_usage = prometheus.NewDesc("cas_vm_memory_usage","cas virtual machine memory usage",
		[]string{"cas_vm_id","cas_cvk_host_id",},nil)
)



func (c VmCollector)Collect(ch chan<-prometheus.Metric ) {
	monitor_info := config.GetVmMonitorInfo(c.Target)
	//hostPoolId := monitor_info.Params_maps["hostPoolId"]
	hostId := monitor_info.HostId
	//clusterId := monitor_info.Params_maps["clusterId"]
	vmId := monitor_info.VmId
	ip := monitor_info.IP
	port := monitor_info.Port
	username := monitor_info.Username
	password := monitor_info.Password
	if ip =="" || username == "" ||password == "" || port == ""  || vmId == "" || hostId == ""{
		ch <- prometheus.MustNewConstMetric(cas_vm_monitorstatus,prometheus.GaugeValue,float64(0))
		return
	}

	errs,ret:=GetAnotherData(fmt.Sprintf("/cas/casrs/vm/id/%s/monitor",vmId),ip,username,password,port)
	info := &config.VMMonitor{}
	if errs!=nil {
		log.Printf("error get vm %s info :%s",vmId,errs.Error())
		ch <- prometheus.MustNewConstMetric(cas_vm_monitorstatus,prometheus.GaugeValue,float64(0))
		return
	}
	ch <- prometheus.MustNewConstMetric(cas_vm_monitorstatus,prometheus.GaugeValue,float64(1))
	errs=xml.Unmarshal(ret,&info)
	metricValue,err:=strconv.ParseFloat(info.CpuRate,64)
	if err!=nil {
		log.Printf("error parse vm %s cpurate info :%s",vmId,err.Error())
	}
	ch<-prometheus.MustNewConstMetric(cas_vm_cpu_usage,prometheus.GaugeValue,metricValue,vmId,hostId)
	memortRate,err:=strconv.ParseFloat(info.MemRate,64)
	ch<-prometheus.MustNewConstMetric(cas_vm_memory_usage,prometheus.GaugeValue,memortRate,vmId,hostId)
	if err!=nil {
		log.Printf("error parse vm %s memrate info :%s",vmId,err.Error())
	}

}
/*
func GetVirtualMachineInfo(uuid,hostId,vmId,cas_ip,cas_port,cas_password,cas_username string){
	if uuid!=nil {
		errs,ret := getdata("/cas/casrs/vm/vmList",cas_ip,cas_port,cas_password,cas_username)
		if errs!=nil{
			log.Printf("error request:%s",errs)
			return
		}
		if ret==nil {
			log.Printf("/cas/casrs/vm/vmList response is null")
			return
		}
		cvm_list:=&config.CVM_VMList{}
		errs=nil
		errs=xml.Unmarshal(ret,&cvm_list)
		if errs!=nil {
			log.Printf("Unmarshal error: %s" ,errs)
			return
		}
		count:=0
		for _,vm := range (*cvm_list).VM{
			if vm.uu {
				
			}
		}
	}
}

func getVmdata(cas_ip,cas_port,cas_password,cas_username string,ch chan<-prometheus.Metric)()  {

}
func getdata(string,string,string,string, string)(error,[]byte)  {

}
*/
