// Package nacos @Author lanpang
// @Date 2024/8/2 下午2:05:00
// @Desc
package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var cidrs []*net.IPNet
var tablerow []string

func NewNacos(c Config, writefile string) *Nacos {
	return &Nacos{
		Config: c,
		//Webport:   ":" + webport,
		Writefile: writefile,
	}
}

func (d *Nacos) WithAuth() {
	log.Println("更新 token")
	_url := fmt.Sprintf("%s/nacos/v1/auth/login", d.Config.Server)
	formData := map[string]string{
		"username": d.Config.Username,
		"password": d.Config.Password,
	}
	res := d.post(_url, formData)
	if len(gjson.GetBytes(res, "accessToken").String()) != 0 {
		log.Println("Authentication successful...")
		d.Token = gjson.GetBytes(res, "accessToken").String()
	} else {
		log.Println("Authentication failed!")
	}
}
func ContainerdIPCheck(ip string) bool {
	for i := range cidrs {
		if cidrs[i].Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}
func (d *Nacos) GetService(url string, namespaceId string, group string) []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/service/list?pageNo=1&pageSize=500&namespaceId=%s&groupName=%s", url, namespaceId, group)
	res := d.get(_url)
	return res
}

func (d *Nacos) GetInstance(url string, servicename string, namespaceId string, group string) []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/instance/list?serviceName=%s&namespaceId=%s&groupName=%s", url, servicename, namespaceId, group)
	//fmt.Println(_url)
	res := d.get(_url)
	return res
}

func InString(filed string, array []string) bool {
	for _, str := range array {
		if filed == str {
			return true
		}
	}
	return false
}

func (d *Nacos) tableAppend(table *tablewriter.Table, data []string) {
	datastr := strings.Join(data, "-")
	if !InString(datastr, tablerow) {
		tablerow = append(tablerow, datastr)
		table.Append(data)
	}
}
func (d *Nacos) TableRender() {
	tablerow = []string{}
	nacosServer := d.Clusterdata
	tabletitle := []string{"命名空间", "服务名称", "实例", "健康状态", "主机名", "权重", "容器", "组"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	for _, v := range nacosServer.HealthInstance {
		tabledata := []string{v.NamespaceName, v.ServiceName, v.IpAddr, v.Health, v.Hostname, v.Weight, v.Container, v.GroupName}
		d.tableAppend(table, tabledata)
	}
	fmt.Printf("健康实例:(%d 个)\n", table.NumLines())
	table.Render()
	if len(nacosServer.UnHealthInstance) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(tabletitle)
		for _, v := range nacosServer.UnHealthInstance {
			tabledata := []string{v.NamespaceName, v.ServiceName, v.IpAddr, v.Health, v.Hostname, v.Weight, v.Container, v.GroupName}
			d.tableAppend(table, tabledata)
		}
		fmt.Printf("异常实例:(%d 个)\n", table.NumLines())
		table.Render()
	}
}

func (d *Nacos) GetNacosInstance() {
	log.Println("更新服务数据")
	var ser Service
	var cluster ClusterStatus
	_url := d.Config.Server
	namespace := d.Config.Namespace
	group := "DEFAULT_GROUP"
	res := d.GetService(_url, namespace, group)
	err := json.Unmarshal(res, &ser)
	if err != nil {
		fmt.Println(err)
	}
	for _, se := range ser.Doms {
		res := d.GetInstance(_url, se, namespace, group)
		var in Instance
		err := json.Unmarshal(res, &in)
		if err != nil {
			fmt.Printf("json序列化错误:%s\n", err)
		}
		for _, host := range in.Hosts {
			instance := ServerInstance{
				NamespaceName: namespace,
				ServiceName:   se,
				IpAddr:        fmt.Sprintf("%s:%d", host.Ip, host.Port),
				Health:        strconv.FormatBool(host.Healthy),
				Hostname:      host.Ip,
				Weight:        fmt.Sprintf("%.1f", host.Weight),
				Container:     strconv.FormatBool(ContainerdIPCheck(host.Ip)),
				Ip:            host.Ip,
				Port:          strconv.Itoa(host.Port),
				GroupName:     in.GroupName,
			}
			if host.Healthy {
				cluster.HealthInstance = append(cluster.HealthInstance, instance)
			} else {
				cluster.UnHealthInstance = append(cluster.UnHealthInstance, instance)
			}
			//fmt.Println(instance)
		}
		d.Clusterdata = cluster
	}
	//fmt.Println(d.Clusterdata.HealthInstance)
}
func (d *Nacos) GetJson(resultType string) (result interface{}, err error) {
	//mutex.Lock()
	//defer mutex.Unlock()
	defer func() {
		if funcErr := recover(); funcErr != nil {
			result = []byte("[]")
			err = errors.New("error")
		}
	}()
	if d.Web {
		d.GetNacosInstance()
	}
	var nacos Nacosfile
	nacosServer := d.Clusterdata
	if len(nacosServer.HealthInstance) != 0 {
		for _, na := range nacosServer.HealthInstance {
			var ta Nacostarget
			ta.Labels = make(map[string]string)
			ta.Targets = append(ta.Targets, na.IpAddr)
			ta.Labels["namespace"] = na.NamespaceName
			ta.Labels["service"] = na.ServiceName
			ta.Labels["hostname"] = na.Hostname
			ta.Labels["weight"] = na.Weight
			ta.Labels["ip"] = na.Ip
			ta.Labels["port"] = na.Port
			ta.Labels["group"] = na.GroupName
			ta.Labels["containerd"] = na.Container
			nacos.Data = append(nacos.Data, ta)
		}

	}

	if resultType == "json" {
		result = nacos.Data
		return result, err
	}
	data, err := json.MarshalIndent(&nacos.Data, "", "  ")
	if err != nil {
		fmt.Println("json序列化失败!")
		result = []byte("[]")
		return result, err
	}
	result = data
	//result = []byte("[]")
	return result, err
}

func (d *Nacos) WriteFile() {
	var basedir string
	var filename string
	basedir = path.Dir(d.Writefile)
	filename = path.Base(d.Writefile)
	if err := os.MkdirAll(basedir, os.ModePerm); err != nil {
		os.Exit(1)
	}
	file, err := os.OpenFile(basedir+"/.nacos_tmp.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("创建文件失败 %s", err)
		os.Exit(2)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	jsondata, err := d.GetJson("byte")
	data := make([]byte, 0)
	var check bool
	if data, check = jsondata.([]byte); !check {
		log.Println("转换失败")
	}
	if _, err := file.Write(data); err != nil {
		log.Println("写入失败", err)
		os.Exit(1)
	}
	err = file.Close()
	if err != nil {
		return
	}
	if err := os.Rename(basedir+"/.nacos_tmp.json", basedir+"/"+filename); err != nil {
		log.Println("写入失败:", basedir+"/"+filename)
	} else {
		log.Println("写入成功:", basedir+"/"+filename)
	}
}
