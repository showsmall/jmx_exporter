package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sre/jmx_exporter/config"
	"time"
)

type JmxData struct {
	Value     Value  `json:"value"`
	Error     string `json:"error"`
	Status    int    `json:"status"`
	Timestamp int    `json:"timestamp"`
}
type Value struct {
	CollectionCount int `json:"CollectionCount"`
	ThreadCount     int `json:"ThreadCount"`
}

func GetJmxData(data *map[string]interface{}) (re *JmxData, err error) {
	// 初始化返回值结构体
	r_jmxdata := &JmxData{}

	// 初始化请求
	transport := http.Transport{}
	client := http.Client{
		Transport: &transport,
	}
	// 配置请求超时
	client.Timeout = 120 * time.Second
	postjson, err := json.Marshal(data)
	if err != nil {
		log.Println("post的map数据转json出错:", err)
		return nil, err
	}
	log.Println(string(postjson))
	// 构造请求
	req, err := http.NewRequest("POST", config.C.GetString("jolokia.url"), bytes.NewBuffer(postjson))
	req.SetBasicAuth(config.C.GetString("jolokia.username"), config.C.GetString("jolokia.password"))
	if err != nil {
		log.Println("构造请求出错:", err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println("发起请求出现错误:", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		jmxdata, _ := ioutil.ReadAll(res.Body)
		log.Println("请求响应状态码不是200，查看日志", string(jmxdata))
		_ = res.Body.Close()
		return nil, errors.New(fmt.Sprintf("请求响应状态码不是200，查看日志:%s", string(jmxdata)))
	}
	jmxdata, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(jmxdata, r_jmxdata)
	if err != nil {
		log.Println("接收到的jmxdata转struct失败:", err)
		return nil, err
	}
	//log.Println(string(jmxdata))
	if r_jmxdata.Status != 200 {
		log.Println("接收到的jmxdataStatus != 200 :", r_jmxdata.Error)
		return nil, errors.New(r_jmxdata.Error)
	}

	_ = res.Body.Close()
	return r_jmxdata, nil
}
