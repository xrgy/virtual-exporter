package config

import (
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"time"
	"database/sql"
)

var db *sql.DB

type ConnectInfoData struct {
	IP          string
	Port string
	Username string
	Password string
	HostpoolId string
	//ClusterId string
	HostId string
	VmId string
}
type ConnectInfo struct {
	ip     string
	port string
	username string
	password string
	hostpoolId string
	//clusterId string
	hostId string
	vmId string
}

func GetDBHandle() *sql.DB {
	var err error
	DBUsername := os.Getenv("DB_USERNAME")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBEndpoint := os.Getenv("DB_ENDPOINT")
	DBDatabase := os.Getenv("DB_DATABASE")
	dsn := DBUsername + ":" + DBPassword + "@(" + DBEndpoint + ")/" + DBDatabase
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("get DB handle error: %v", err)
	}
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(28000 * time.Second)
	err = db.Ping()
	if err != nil {
		log.Printf("connecting DB error: %v ", err)
	}
	return db
}
func GetCasMonitorInfo(id string) ConnectInfoData {
	info := queryCasConnectInfo(id)

	con_info_data:=ConnectInfoData{
		info.ip,
		info.port,
		info.username,
		info.password,
		"",
		"",
		"",
	}
	return con_info_data
}
func queryCasConnectInfo(id string) ConnectInfo {
	rows, err := db.Query("select ip,username,password,port from tbl_cas_monitor_record where uuid=?", id)
	if err != nil {
		log.Printf("query error")
	}
	info := ConnectInfo{}
	for rows.Next() {
		err = rows.Scan(&info.ip, &info.username,&info.password,&info.port)
	}
	defer rows.Close()
	return info
}

func GetCvkMonitorInfo(id string) ConnectInfoData {
	info := queryCvkConnectInfo(id)

	con_info_data:=ConnectInfoData{
		info.ip,
	  info.port,
	  info.username,
	  info.password,
	  info.hostpoolId,
	  info.hostId,
	  "",
	}
	return con_info_data
}
func queryCvkConnectInfo(id string) ConnectInfo {
	rows, err := db.Query("select ip,hostId,hostpoolId,casUuid from tbl_host_monitor_record where uuid=?", id)
	if err != nil {
		log.Printf("query error")
	}
	info := ConnectInfo{}
	casUuid:=""
	for rows.Next() {
		err = rows.Scan(&info.ip, &info.hostId,&info.hostpoolId,&casUuid)
	}

	rows2, err := db.Query("select username,password,port from tbl_cas_monitor_record where uuid=?", casUuid)
	if err != nil {
		log.Printf("query error")
	}
	for rows2.Next() {
		err = rows2.Scan(&info.username,&info.password,&info.port)
	}

	defer rows.Close()
	defer rows2.Close()
	return info
}

func GetVmMonitorInfo(id string) ConnectInfoData {
	info := queryVmConnectInfo(id)

	con_info_data:=ConnectInfoData{
		info.ip,
		info.port,
		info.username,
		info.password,
		"",
		info.hostId,
		info.vmId,
	}
	return con_info_data
}
func queryVmConnectInfo(id string) ConnectInfo {
	rows, err := db.Query("select ip,vmId,cvkUuid from tbl_vm_monitor_record where uuid=?", id)
	if err != nil {
		log.Printf("query error")
	}
	info := ConnectInfo{}
	cvkUuid:=""
	for rows.Next() {
		err = rows.Scan(&info.ip, &info.vmId,&cvkUuid)
	}

	rows2, err := db.Query("select hostId,casUuid from tbl_host_monitor_record where uuid=?", cvkUuid)
	if err != nil {
		log.Printf("query error")
	}
	casUuid:=""
	for rows2.Next() {
		err = rows2.Scan(&info.hostId,&casUuid)
	}

	rows3, err := db.Query("select username,password,port from tbl_cas_monitor_record where uuid=?", casUuid)
	if err != nil {
		log.Printf("query error")
	}
	for rows3.Next() {
		err = rows3.Scan(&info.username,&info.password,&info.port)
	}

	defer rows.Close()
	defer rows2.Close()
	defer rows3.Close()
	return info
}


func CloseDBHandle()  {
	db.Close()
}