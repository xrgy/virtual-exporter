package config

type HostpoolList struct {
	HostPool []HostPool `xml:"hostPool"`
}
type HostPool struct {
	Id           string `xml:"id"`
	HostpoolId   string `xml:"hostpoolId"`
	HostpoolName string `xml:"hostpoolName"`
	Name         string `xml:"name"`
	Cluster []Cluster `xml:"cluster"`
	Host []Host `xml:"host"`
}
type ResourceData struct {
	HostPool []HostPool `xml:"hostPool"`
}

type ChildNode struct {
	ClusterList []ClusterInfo `xml:"clusterList"`
	NodeHostList []HostInfo `xml:"nodeHostList"`
}
type ClusterInfo struct {
	Id          string `xml:"id"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
}
type Cluster struct {
	ClusterId   string `xml:"clusterId"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Host []Host `xml:"host"`
}

type ClusterHosts struct {
	ClusterHostInfo []HostInfo `xml:"clusterHostInfo"`
}
type Host struct {
	Name   string `xml:"name"`
	Id     string `xml:"id"`
	Status string `xml:"status"`
	Ip     string `xml:"ip"`
	Vm  []VM `xml:"vm"`
}

type HostInfoList struct {
	HostInfo []HostInfo `xml:"host"`
}

type HostInfo struct {
	Id         string `xml:"id"`
	Ip         string `xml:"ip"`
	HostPoolId string `xml:"hostPoolId"`
	ClusterId  string `xml:"clusterId"`
	Name       string `xml:"name"`
	VmNum      string `xml:"vmNum"`
	VmRunCount string `xml:"vmRunCount"`
	VmShutoff  string `xml:"vmShutoff"`
	CpuRate    string `xml:"cpuRate"`
	MemRate    string `xml:"memRate"`
	Status     string `xml:"status"`
	Storage    string `xml:"storage"`
	Version    string `xml:"version"`
}

type HostInfoDetail struct {
	Id string `xml:"id"`
	Name string `xml:"name"`
	Ip string `xml:"ip"`
	Model string `xml:"model"`
	Vendor string `xml:"vendor"`
	CpuCount string `xml:"cpuCount"`
	CpuModel string `xml:"cpuModel"`
	CpuFrequence string `xml:"cpuFrequence"`
	DiskSize uint64 `xml:"diskSize"`
	MemorySize uint64 `xml:"memorySize"`
	Status string `xml:"status"`
	CpuSockets string `xml:"cpuSockets"`
	CpuCores string `xml:"cpuCores"`
}

type VMMonitor struct {
	Id           string `xml:"id"`
	CpuRate      string `xml:"cpuRate"`
	MemRate      string `xml:"memRate"`
	Status       string `xml:"status"`
	Uuid         string `xml:"uuid"`
	Flag         string `xml:"flag"`
	Type         string `xml:"type"`
	Deployed     string `xml:"deployed"`
	Hoststatus   string `xml:"hoststatus"`
	HostHaEnable string `xml:"hostHaEnable"`
	//VMDisk string 	`xml:"vmDisk"`
}
type VMList struct {
	VM []HostVM `xml:"vm"`
}
type HostVM struct {
	Id string `xml:"id"`
	Title string `xml:"title"`
	VMStatus string	`xml:"vmStatus"`
	OsDesc string `xml:"osDesc"`
} 
type VM struct {
	Id string `xml:"id"`
	Name string `xml:"name"`
	Status string 	`xml:"status"`
	Os string `xml:"os"`
	Ip string `xml:"ip"`
}
type VMNetworkList struct {
	VMNetwork []VMNetwork `xml:"vmNetwork"`	
}
type VMNetwork struct {
	IpAddr string	`xml:"ipAddr"`
}

type CVM_VMList struct {
    VM []VM `xml:"vm"`
}

type PNICList struct {
	PNIC []PNIC `xml:"pnic"`
}
type PNIC struct {
	Name string `xml:"name"`
	Description string `xml:"description"`
	Status int64 	`xml:"status"`
	MacAddr string `xml:"macAddr"`
	Speed string `xml:"speed"`
	Duplex string `xml:"duplex"`
	Carrier string `xml:"carrier"`
	Mtu string `xml:"Mtu"`
}

type CasClusterDetails struct {
	CasClusterInfo []CasClusterDetail `json:"keyValue"`
}
type CasClusterDetail struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type HostMonitor struct {
	CpuRate      string `xml:"cpuRate"`
	MemRate      string `xml:"memRate"`
	Disk HostDisk `xml:"disk"`
}
type HostDisk struct {
	Device string `xml:"device"`
	Usage string `xml:"usage"` //本地存储使用百分比
}

type PNICTraffics struct {
	TrendRate []TrendRate `xml:"trendRate"`
}
type TrendRate struct {
	Name string `xml:"name"`
	Rates []Rates `xml:"rates"`
}

type Rates struct {
	Time int64 `xml:"time"`
	Rate string `xml:"rate"`
}
type VSwitchInfo struct {
	Vswitch []VSwitch `xml:"vSwitch"`
}

type VSwitch struct {
	Id string `xml:"id"`
	HostId string `xml:"hostId"`
	Name string `xml:"name"`
	PortNum string `xml:"portNum"`
	Mode string `xml:"mode"`
	Pnic string `xml:"pnic"`
	Address string 	`xml:"address"`
	Netmask string 	`xml:"netmask"`
	Gateway string 	`xml:"gateway"`
	EnableLacp string `xml:"enableLacp"`
	BondMode string `xml:"bondMode"`
	IsManage string `xml:"isManage"`
	Status int64 `xml:"status"`
}
type VportTrafficInfos struct {
	VportTrafficInfo []VportTrafficInfo `xml:"trafficInfo"`
}
type VportTrafficInfo struct {
	VsName string `xml:"vsName"`
	VPortName string `xml:"vPortName"`
	VmName string `xml:"vmName"`
	Title string `xml:"title"`
	VmMac string `xml:"vmMac"`
	ReceivePkts int64 `xml:"rxPkts"`
	ReceiveBytes int64 `xml:"rxBytes"`
	ReceiveErrs int64 `xml:"rxErrs"`
	SendPkts int64 `xml:"txPkts"`
	SendBytes int64 `xml:"txBytes"`
	SendErrs int64 `xml:"txErrs"`
}

type HostStoragePools struct {
	HostStoragePools []HostStoragePool `xml:"storagePool"`
}
type HostStoragePool struct {
	Name string `xml:"name"`
	Path string `xml:"path"`
	RemainSize string `xml:"remainSize"`
	Status int64 `xml:"status"`
	Title string `xml:"title"`
	Type string `xml:"type"`
}

type HostStorageVolumes struct {
	HostStorageVolumes []HostStorageVolume `xml:"storageVolume"`
}

type HostStorageVolume struct {
	Allocation string `xml:"allocation"`
	BackingStore string `xml:"backingStore"`
	BaseFile string `xml:"baseFile"`
	Format string `xml:"format"`
	Name string `xml:"name"`
	Size string `xml:"size"`
}