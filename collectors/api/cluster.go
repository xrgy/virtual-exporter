package api

import (
	"net/http"
	"virtual-exporter/config"
	"encoding/xml"
	"fmt"
	"encoding/json"
	"log"
)

func GetAllCluster(w http.ResponseWriter, r *http.Request) {
	ip:=r.URL.Query().Get("ip")
	username:=r.URL.Query().Get("username")
	password:=r.URL.Query().Get("password")
	port:=r.URL.Query().Get("port")
	err,ret:=GetAnotherData("/cas/casrs/hostpool/all",ip,username,password,port)

	hostPools := &config.HostpoolList{}
	resourceData := &config.ResourceData{}
	err=xml.Unmarshal(ret,&hostPools)

	for i:=0;i<len((*hostPools).HostPool);i++ {
		hostPool := &config.HostPool{}
		hostPoolId := (*hostPools).HostPool[i].Id
		hostPool.HostpoolId = hostPoolId
		hostPool.HostpoolName = (*hostPools).HostPool[i].Name
		errs, ret := GetAnotherData(fmt.Sprintf("/cas/casrs/hostpool/%s/allChildNode", hostPoolId), ip, username, password, port)
		if err!=nil{
			log.Printf("get hostpool child error",errs)
			continue
		}
		childNodeInfo := &config.ChildNode{}
		errs = xml.Unmarshal(ret, childNodeInfo)
		for _,clusterInfo:=range (*childNodeInfo).ClusterList {
			cluster := &config.Cluster{}
			cluster.ClusterId = clusterInfo.Id
			cluster.Name = clusterInfo.Name
			cluster.Description = clusterInfo.Description
			hostPool.Cluster = append(hostPool.Cluster, *cluster)
		}
		resourceData.HostPool=append(resourceData.HostPool,*hostPool)
	}
	ret,errs:=json.Marshal(*resourceData)
	w.Header().Set("Content-type","application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(ret)
	if err!=nil{
		log.Printf("get hostpool child error",errs)
	}
	//return errs
}