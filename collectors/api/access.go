package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"errors"

	"fmt"
)
type AccessReq struct {
	monitorInfo map[string]string `json:"monitorInfo"`
}

type AccessResp struct {
	Result map[string]string `json:"result"`
}

func getResponse(w http.ResponseWriter,err error)  {
	resultMap := make(map[string]string,2)
	if err!=nil {
		resultMap["accessible"] = "false"
		resultMap["message"] = err.Error()
	}else {
		resultMap["accessible"] = "true"
		resultMap["message"] = ""
	}
	accessResp := AccessResp{Result:resultMap}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(accessResp)
}

func CasAccess(w http.ResponseWriter,r *http.Request)  {
	var info AccessReq
    body,err := ioutil.ReadAll(r.Body)
	if err!=nil {
		log.Printf("get body data error:"+err.Error())
	}
	err = json.Unmarshal(body,&info)
	ip := info.monitorInfo["ip"]
	username := info.monitorInfo["username"]
	password := info.monitorInfo["password"]
	port := info.monitorInfo["port"]
	if ip =="" || username == "" ||password == "" || port == "" {
		log.Printf("illegal request bosy:%s",string(body))
		http.Error(w,"illegal request bosy",http.StatusBadRequest)
		return
	}
	command := "ping -i 0.3 -w 5 "+ip+" -c 3 | tail -n 2"
	cmd := exec.Command("/bin/sh","-c",command)
	ret,_ := cmd.Output()
	s := string(ret)
	if strings.Contains(s,"100% packet loss") {
		getResponse(w,errors.New("ip doesn't exist."))
		return
	}

	err,_=GetAnotherData("/cas/casrs/hostpool/all",ip,username,password,port)
	getResponse(w,err)
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