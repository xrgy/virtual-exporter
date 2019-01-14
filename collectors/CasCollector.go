package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"virtual-exporter/config"
	"net/http"
	"sync"
	"encoding/xml"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type CasCollector struct {
	Target string
}
var (
	cas_monitorstatus = prometheus.NewDesc(prometheus.BuildFQName("cas","","monitorstatus"),
		"cas monitorstatus",nil,nil)
	)
func (c CasCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

func GetCasInfo(w http.ResponseWriter,r *http.Request) (err error) {
		ip:=r.URL.Query().Get("ip")
		username:=r.URL.Query().Get("username")
		password:=r.URL.Query().Get("password")
		port:=r.URL.Query().Get("port")
		err,ret:=GetAnotherData("/cas/casrs/hostpool/all",ip,username,password,port)
		locker:=sync.RWMutex{}

		hostPools := &config.HostpoolList{}
		resourceData := &config.ResourceData{}
		err=xml.Unmarshal(ret,&hostPools)

		for i:=0;i<len((*hostPools).HostPool);i++{
			hostPool:=&config.HostPool{}
			hostPoolId := (*hostPools).HostPool[i].Id
			hostPool.HostpoolId = hostPoolId
			hostPool.HostpoolName= (*hostPools).HostPool[i].Name
			errs,ret:=GetAnotherData(fmt.Sprintf("/cas/casrs/hostpool/%s/allChildNode",hostPoolId),ip,username,password,port)
			childNodeInfo := &config.ChildNode{}
			errs = xml.Unmarshal(ret,&childNodeInfo)
			wg:=sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				for _,clusterInfo:=range (*childNodeInfo).ClusterList  {
					cluster:=&config.Cluster{}
					cluster.ClusterId=clusterInfo.Id
					cluster.Name=clusterInfo.Name
					cluster.Description=clusterInfo.Description
					errs,ret=GetAnotherData(fmt.Sprintf("/cas/casrs/cluster/hosts?offset=0&limit=2000&clusterId=%s",clusterInfo.Id),ip,username,password,port)
					clusterHosts:=&config.ClusterHosts{}
					errs=xml.Unmarshal(ret,clusterHosts)
					for _,hostInfo:=range (*clusterHosts).ClusterHostInfo{
						host:=&config.Host{}
						host.Name=hostInfo.Name
						host.Id=hostInfo.Id
						host.Status = hostInfo.Status
						host.Ip = hostInfo.Ip
						errs,ret=GetAnotherData(fmt.Sprintf("/cas/casrs/vm/vmList?hostId=%s&sortField=id&sortDir=Asc",hostInfo.Id),ip,username,password,port)
						vmList:=&config.VMList{}
						errs=xml.Unmarshal(ret,vmList)
						for i:=0;i<len((*vmList).VM);i++{
							hostVm:=&config.VM{}
							vm:=(*vmList).VM[i]
							hostVm.Id = vm.Id
							hostVm.Name=vm.Title
							hostVm.Status = vm.VMStatus
							hostVm.Os = vm.OsDesc
							errs,ret=GetAnotherData(fmt.Sprintf("/cas/casrs/vm/network/%s",vm.Id),ip,username,password,port)
							vmNetwork:=&config.VMNetworkList{}
							errs=xml.Unmarshal(ret,vmNetwork)
							ip:=""
							for _,network:=range (*vmNetwork).VMNetwork{
								if(network.IpAddr!=""){
									ip=network.IpAddr
									break
								}
							}
							hostVm.Ip=ip
							host.Vm=append(host.Vm,*hostVm)
						}
						cluster.Host=append(cluster.Host,*host)
					}
					locker.Lock()
					hostPool.Cluster = append(hostPool.Cluster,*cluster)
					locker.Unlock()
				}
			}()
			go func() {
				defer wg.Done()
				for _,hostInfo:=range (*childNodeInfo).NodeHostList{
					host:=&config.Host{}
					host.Name = hostInfo.Name
					host.Id = hostInfo.Id
					host.Status = hostInfo.Status
					host.Ip = hostInfo.Ip
					errs,ret=GetAnotherData(fmt.Sprintf("/cas/casrs/vm/vmList?hostId=%s&sortField=id&sortDir=Asc",hostInfo.Id),ip,username,password,port)
					vmList:=&config.VMList{}
					errs=xml.Unmarshal(ret,vmList)
					for i:=0;i<len((*vmList).VM);i++{
						hostVm:=&config.VM{}
						vm:=(*vmList).VM[i]
						hostVm.Id = vm.Id
						hostVm.Name=vm.Title
						hostVm.Status = vm.VMStatus
						hostVm.Os = vm.OsDesc
						errs,ret=GetAnotherData(fmt.Sprintf("/cas/casrs/vm/network/%s",vm.Id),ip,username,password,port)
						vmNetwork:=&config.VMNetworkList{}
						errs=xml.Unmarshal(ret,vmNetwork)
						ip:=""
						for _,network:=range (*vmNetwork).VMNetwork{
							if(network.IpAddr!=""){
								ip=network.IpAddr
								break
							}
						}
						hostVm.Ip=ip
						host.Vm=append(host.Vm,*hostVm)
					}
					locker.Lock()
					hostPool.Host=append(hostPool.Host,*host)
					locker.Unlock()
				}
			}()
			wg.Wait()
			resourceData.HostPool=append(resourceData.HostPool,*hostPool)
		}
		ret,errs:=json.Marshal(*resourceData)
		w.Header().Set("Content-type","application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(ret)
		return errs
}

func GetAnotherData(url string,ip string, username string, password string, port string) (error,[]byte) {
  //wating test
	path:=fmt.Sprintf("%s:%s/%s"+ip,port,url)
	//client := &http.Client{}
	//resp,err:=client.Get(path)
	req, err := http.NewRequest("GET",path,nil)
	req.SetBasicAuth(username,password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("read response error:%s ",err)
		return err,nil
	}
	bytes,err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		log.Printf("read response body error:%s ",err)
		return err,nil
	}
	return nil,bytes
}

func (c CasCollector) Collect(ch chan<- prometheus.Metric) {
	monitor_info:= config.GetCasMonitorInfo(c.Target)
	ip:= monitor_info.IP
	port:= monitor_info.Port
	username:=monitor_info.Username
	password := monitor_info.Password
	if ip =="" || username == "" ||password == "" || port == ""  {
		ch <- prometheus.MustNewConstMetric(cas_monitorstatus,prometheus.GaugeValue,float64(0))
		return
	}
	command := "ping -i 0.3 -w 5 "+ip+" -c 3 | tail -n 2"
	cmd := exec.Command("/bin/sh","-c",command)
	ret,_ := cmd.Output()
	s := string(ret)
	if strings.Contains(s,"100% packet loss") {
		ch <- prometheus.MustNewConstMetric(cas_monitorstatus,prometheus.GaugeValue,float64(0))
		return
	}
	err,_:=GetAnotherData("/cas/casrs/hostpool/all",ip,username,password,port)
	if err!=nil {
		ch <- prometheus.MustNewConstMetric(cas_monitorstatus,prometheus.GaugeValue,float64(0))
		return
	}
	ch <- prometheus.MustNewConstMetric(cas_monitorstatus,prometheus.GaugeValue,float64(1))

}