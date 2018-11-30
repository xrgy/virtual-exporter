package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"virtual-exporter/config"
	"fmt"
	"encoding/xml"
	"strconv"
	"log"
	"strings"
)

type CvkCollector struct {
	Target string
}

func (c CvkCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

var (
	cas_cvk_host_cpu_usage = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_cpu_usage"), "cas cvk host cpu usage",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_model", "cas_cvk_host_vendor",}, nil)

	cas_cvk_host_memory_usage = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_memory_usage"), "cas cvk host memory usage",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name","cas_cvk_host_model", "cas_cvk_host_vendor",}, nil)
	cas_cvk_host_memory_size = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_memory_size"), "cas cvk host memory size",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_model", "cas_cvk_host_vendor",}, nil)

	cas_cvk_host_disk_size = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_memory_size"), "cas cas cvk host disk size",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_model", "cas_cvk_host_vendor",}, nil)

	cas_cvk_host_total_disk_usage = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_total_disk_usage"), "cas cvk host total disk usage",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_model", "cas_cvk_host_vendor", "cas_cvk_host_disk_device"}, nil)

	cas_cvk_host_pNIC_status = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_pNIC_status"), "cas cas cvk host pNIC status",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_pnic_name", "cas_cvk_host_pnic_macaddr"}, nil)

	cas_cvk_host_network_receiveRate = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_network_receiveRate"), "cas cas cvk host network receiveRate",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_pnic_name", "cas_cvk_host_pnic_macaddr"}, nil)

	cas_cvk_host_network_sendRate = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_network_sendRate"), "cas cvk host network sendRate",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_pnic_name", "cas_cvk_host_pnic_macaddr"}, nil)

	cas_cvk_vm_vswitch_status = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_vm_vswitch_status"), "cas cvk vm vswitch status",
		[]string{"cas_cvk_host_id", "cas_cvk_host_vswitch_id", "cas_cvk_host_vswitch_name"}, nil)

	cas_cvk_vm_port_receive_pkts = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_vm_port_receive_pkts"), "cas cvk vm port receive pkts",
		[]string{"cas_cvk_host_id", "cas_cvk_host_vswitch_id", "cas_cvk_host_vswitch_name", "cas_cvk_host_vswitch_portname"}, nil)

	cas_cvk_vm_port_receive_bytes = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_vm_port_receive_bytes"), "cas cvk vm port receive bytes",
		[]string{"cas_cvk_host_id", "cas_cvk_host_vswitch_id", "cas_cvk_host_vswitch_name", "cas_cvk_host_vswitch_portname"}, nil)

	cas_cvk_vm_port_send_pkts = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_vm_port_send_pkts"), "cas cvk vm port send pkts",
		[]string{"cas_cvk_host_id", "cas_cvk_host_vswitch_id", "cas_cvk_host_vswitch_name", "cas_cvk_host_vswitch_portname"}, nil)

	cas_cvk_vm_port_send_bytes = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_vm_port_send_bytes"), "cas cvk vm port send bytes",
		[]string{"cas_cvk_host_id", "cas_cvk_host_vswitch_id", "cas_cvk_host_vswitch_name", "cas_cvk_host_vswitch_portname"}, nil)

	cas_cvk_host_storage_status = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_storage_status"), "cas cvk host storage status",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_storage_name", "cas_cvk_host_storage_path","cas_cvk_host_storage_type"}, nil)

	cas_cvk_host_storage_freeSize = prometheus.NewDesc(prometheus.BuildFQName("cas", "", "cvk_host_storage_freeSize "), "cas cvk host storage freeSize",
		[]string{"cas_cvk_host_id", "cas_cvk_host_status", "cas_cvk_host_ip", "cas_cvk_host_name", "cas_cvk_host_storage_name", "cas_cvk_host_storage_path","cas_cvk_host_storage_type"}, nil)


)

func (c CvkCollector) Collect(ch chan<- prometheus.Metric) {
	monitor_info := config.GetMonitorInfo(c.Target)
	hostPoolId := monitor_info.Params_maps["hostPoolId"]
	hostId := monitor_info.Params_maps["hostId"]
	//clusterId := monitor_info.Params_maps["clusterId"]
	ip := monitor_info.Params_maps["ip"]
	port := monitor_info.Params_maps["port"]
	username := monitor_info.Params_maps["username"]
	password := monitor_info.Params_maps["password"]
	//获取该主机池的主机列表
	errs, ret := GetAnotherData(fmt.Sprintf("/cas/casrs/hostpool/host/%s?offset=0&limit=2000", hostPoolId), ip, username, password, port)
	if errs != nil {
		log.Printf("get host in hostpool error:%s", errs.Error())
		return
	}
	a_hostInfo := &config.HostInfoList{}
	errs = xml.Unmarshal(ret, &a_hostInfo)
	//hostInfo := &config.HostInfo{}
	//for _, host := range (*a_hostInfo).HostInfo {
	//	if host.Id == hostId {
	//		hostInfo = &host
	//		break
	//	}
	//}
	errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/host/id/%s", hostId), ip, username, password, port)
	hostInfoDetail := &config.HostInfoDetail{}
	errs = xml.Unmarshal(ret, &hostInfoDetail)
	if errs != nil {
		log.Printf("parse host info detail error:%s", errs.Error())
		return
	}
	hostStatus := hostInfoDetail.Status
	hostIp := hostInfoDetail.Ip
	hostName := hostInfoDetail.Name
	hostModel := hostInfoDetail.Model
	hostVendor := hostInfoDetail.Vendor
	labelvalues := []string{hostId, hostStatus, hostIp, hostName, hostModel, hostVendor}
	errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/host/id/%s/monitor", hostId), ip, username, password, port)
	hostMonitor := &config.HostMonitor{}
	if errs != nil {
		log.Printf("error get host %s info :%s", hostId, errs.Error())
		return
	}
	errs = xml.Unmarshal(ret, &hostMonitor)
	cpuRate, err := strconv.ParseFloat(hostMonitor.CpuRate, 64)
	if err != nil {
		log.Printf("error parse host %s cpurate info :%s", hostId, err.Error())
	} else {
		ch <- prometheus.MustNewConstMetric(cas_cvk_host_cpu_usage, prometheus.GaugeValue, cpuRate, labelvalues...)
	}
	memoryRate, err := strconv.ParseFloat(hostMonitor.MemRate, 64)
	if err != nil {
		log.Printf("error parse host %s memrate info :%s", hostId, err.Error())
	} else {
		ch <- prometheus.MustNewConstMetric(cas_cvk_host_memory_usage, prometheus.GaugeValue, memoryRate, labelvalues...)
	}
	ch <- prometheus.MustNewConstMetric(cas_cvk_host_memory_size, prometheus.GaugeValue, float64(hostInfoDetail.MemorySize), labelvalues...) //KB
	ch <- prometheus.MustNewConstMetric(cas_cvk_host_disk_size, prometheus.GaugeValue, float64(hostInfoDetail.DiskSize), labelvalues...)     //MB
	diskUsage, err := strconv.ParseFloat(config.Substr(hostMonitor.Disk.Usage, 0, 2), 64)
	if err != nil {
		log.Printf("error parsedisk usage :%s", err.Error())
	} else {
		labelvaluesdisk := []string{hostId, hostStatus, hostIp, hostName, hostModel, hostVendor, hostMonitor.Disk.Device}
		ch <- prometheus.MustNewConstMetric(cas_cvk_host_total_disk_usage, prometheus.GaugeValue, diskUsage, labelvaluesdisk...) //%
	}
	errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/host/allPhy?id=%s", hostId), ip, username, password, port)
	pnicList := &config.PNICList{}
	if errs != nil {
		log.Printf("error get host pnic %s info :%s", hostId, errs.Error())
		return
	}
	errs = xml.Unmarshal(ret, &pnicList)
	for _, pnic := range (*pnicList).PNIC {
		pnicStatus := pnic.Status
		pnicName := pnic.Name
		pnicMac := pnic.MacAddr
		pniclabelvalues := []string{hostId, hostStatus, hostIp, hostName, pnicName, pnicMac}
		ch <- prometheus.MustNewConstMetric(cas_cvk_host_pNIC_status, prometheus.GaugeValue, float64(pnicStatus), pniclabelvalues...)
		if pnicStatus == 1 {
			errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/host/pnic/traffic?mac=%s", pnicMac), ip, username, password, port)
			pnicTraffics := &config.PNICTraffics{}
			if errs != nil {
				log.Printf("error get host pnic traffic info :%s", errs.Error())
				return
			}
			errs = xml.Unmarshal(ret, &pnicTraffics)
			for _, trendRate := range (*pnicTraffics).TrendRate {
				if strings.Contains(trendRate.Name, "接收流量") {
					timeFlag := int64(0)
					var rRateFlag string
					for _, rtrend := range trendRate.Rates {
						if rtrend.Time > timeFlag {
							timeFlag = rtrend.Time
							rRateFlag = rtrend.Rate
						}
					}
					trendReceive, err := strconv.ParseFloat(rRateFlag, 64)
					if err != nil {
						log.Printf("error parse host %s receive rate info :%s", hostId, err.Error())
					} else {
						ch <- prometheus.MustNewConstMetric(cas_cvk_host_network_receiveRate, prometheus.GaugeValue, trendReceive, pniclabelvalues...)
					}
				}
				if strings.Contains(trendRate.Name, "发送流量") {
					timeFlag := int64(0)
					var sRateFlag string
					for _, strend := range trendRate.Rates {
						if strend.Time > timeFlag {
							timeFlag = strend.Time
							sRateFlag = strend.Rate
						}
					}
					trendSend, err := strconv.ParseFloat(sRateFlag, 64)
					if err != nil {
						log.Printf("error parse host %s send rate info :%s", hostId, err.Error())
					} else {
						ch <- prometheus.MustNewConstMetric(cas_cvk_host_network_sendRate, prometheus.GaugeValue, trendSend, pniclabelvalues...)
					}
				}
			}
		}

	}

	errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/host/id/%s/vswitch", hostId), ip, username, password, port)
	vswitchInfo := &config.VSwitchInfo{}
	if errs != nil {
		log.Printf("error get host %s info :%s", hostId, errs.Error())
		return
	}
	errs = xml.Unmarshal(ret, &vswitchInfo)
	for _, switchi := range (*vswitchInfo).Vswitch {
		vswitchId := switchi.Id
		vswitchName := switchi.Name
		vswitchStatus := switchi.Status
		vsstatuslabels := []string{hostId, vswitchId, vswitchName}
		ch <- prometheus.MustNewConstMetric(cas_cvk_vm_vswitch_status, prometheus.GaugeValue, float64(vswitchStatus), vsstatuslabels...) //0不活动1活动
		errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/host/vs/vport?id=%s&vsId=%s", hostId, vswitchId), ip, username, password, port)
		vportInfo := &config.VportTrafficInfos{}
		if errs != nil {
			log.Printf("error get host %s info :%s", hostId, errs.Error())
			return
		}
		errs = xml.Unmarshal(ret, &vportInfo)
		for _, vport := range (*vportInfo).VportTrafficInfo {
			vportName := vport.VPortName
			labelvaluesvport := []string{hostId, vswitchId, vswitchName, vportName}
			ch <- prometheus.MustNewConstMetric(cas_cvk_vm_port_receive_pkts, prometheus.GaugeValue, float64(vport.ReceivePkts), labelvaluesvport...)
			ch <- prometheus.MustNewConstMetric(cas_cvk_vm_port_receive_bytes, prometheus.GaugeValue, float64(vport.ReceiveBytes), labelvaluesvport...)
			ch <- prometheus.MustNewConstMetric(cas_cvk_vm_port_send_pkts, prometheus.GaugeValue, float64(vport.SendPkts), labelvaluesvport...)
			ch <- prometheus.MustNewConstMetric(cas_cvk_vm_port_send_bytes, prometheus.GaugeValue, float64(vport.SendBytes), labelvaluesvport...)
		}
	}


	errs, ret = GetAnotherData(fmt.Sprintf("/cas/casrs/storage/queryStoragePoolList?id=%s&limit=50&offset=0", hostId), ip, username, password, port)
	hostStoragePools := &config.HostStoragePools{}
	if errs != nil {
		log.Printf("error get host %s info :%s", hostId, errs.Error())
		return
	}
	errs = xml.Unmarshal(ret, &hostStoragePools)
	for _, hostStorage := range (*hostStoragePools).HostStoragePools{
		storageName := hostStorage.Name
		storageType := hostStorage.Type
		storagePath := hostStorage.Path
		storagelabelvalues :=[]string{hostId, hostStatus, hostIp, hostName,storageName,storagePath,storageType}
		ch <- prometheus.MustNewConstMetric(cas_cvk_host_storage_status, prometheus.GaugeValue, float64(hostStorage.Status), storagelabelvalues...) //0不活动 1活动
		remainSize, err := strconv.ParseFloat(hostStorage.RemainSize, 64)
		if err != nil {
			log.Printf("error parse host %s memrate info :%s", hostId, err.Error())
		} else {
			ch <- prometheus.MustNewConstMetric(cas_cvk_host_storage_freeSize, prometheus.GaugeValue, remainSize, storagelabelvalues...)
		}
	}
}
