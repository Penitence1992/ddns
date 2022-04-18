package main

import (
	"fmt"
	dns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net/http"
	_ "org.penitence/ddns/pkg/env"
	"org.penitence/ddns/pkg/resolver"
	"os"
	"strings"
	"time"
)

var gitCommit = ""
var buildStamp = ""

func main() {

	fmt.Printf("Git Commit : %s\n", gitCommit)
	fmt.Printf("Build Stamp : %s\n", buildStamp)

	accessKey := getEnvAndFatalWithEmpty("accessKey")
	accessSecret := getEnvAndFatalWithEmpty("accessSecret")
	baseDomain := getEnvAndFatalWithEmpty("baseDomain")
	domainRRs := getEnvAndFatalWithEmpty("domainRR")
	testUrl := getEnvOrDefault("testUrl", " https://cip.cc")

	log.Println("基础域名为: " + baseDomain)
	log.Printf("三级域名: %s \n", strings.Split(domainRRs, ","))
	log.Println("获取公网ip的测试url: " + testUrl)

	publicIp := findPublicIP(testUrl)

	client, _ := dns.NewClientWithAccessKey("cn-hangzhou", accessKey, accessSecret)

	for _, domainRR := range strings.Split(domainRRs, ",") {
		response, err := describeRR(client, baseDomain, domainRR)
		if err != nil {
			log.Printf("查询dns记录失败:%v", err)
		}
		invokeDDns(response, client, baseDomain, domainRR, publicIp)
	}
}

func describeRR(client *dns.Client, baseDomain, domainRR string) (*dns.DescribeDomainRecordsResponse, error) {
	recordRequest := dns.CreateDescribeDomainRecordsRequest()
	recordRequest.Scheme = "https"
	recordRequest.DomainName = baseDomain
	recordRequest.TypeKeyWord = "A"
	recordRequest.RRKeyWord = domainRR
	return client.DescribeDomainRecords(recordRequest)
}

func invokeDDns(response *dns.DescribeDomainRecordsResponse, client *dns.Client, baseDomain, domainRR, publicIp string) {
	if len(response.DomainRecords.Record) == 0 {
		log.Println("未发现dns解析记录, 添加一个解析")
		recordAddRequest := dns.CreateAddDomainRecordRequest()
		recordAddRequest.Scheme = "https"
		recordAddRequest.DomainName = baseDomain
		recordAddRequest.RR = domainRR
		recordAddRequest.Type = "A"
		recordAddRequest.Value = publicIp
		invokeAliSdk(client.AddDomainRecord(recordAddRequest))
	} else {
		for _, record := range response.DomainRecords.Record {
			log.Printf("发现dns解析记录, RR:%s", record.RR)
			log.Printf("记录内容 : %v\n", record)
			log.Println(record.Value)
			if record.Value == publicIp {
				log.Println("ip没有变动, 跳过更新")
			} else {
				log.Printf("更新域名%s.%s的ip为%s", domainRR, baseDomain, publicIp)
				recordChangeRequest := dns.CreateUpdateDomainRecordRequest()
				recordChangeRequest.Scheme = "https"
				recordChangeRequest.RecordId = record.RecordId
				recordChangeRequest.RR = record.RR
				recordChangeRequest.Type = "A"
				recordChangeRequest.Value = publicIp
				invokeAliSdk(client.UpdateDomainRecord(recordChangeRequest))
			}
		}
	}
}

func findPublicIP(testUrl string) (publicIp string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Get(testUrl)
	if err != nil {
		log.Fatalf("获取互联网ip失败 : %v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("读取请求内容失败 : %v", err)
	}
	publicIp = resolver.FindFirstIp(string(body))
	if publicIp == "" {
		log.Fatalln("响应中未发现ip地址")
	}
	log.Printf("访问%s获取到的互联网ip为:%s\n", testUrl, publicIp)
	return
}

func getEnvOrDefault(env, def string) (v string) {
	v = os.Getenv(env)
	if v == "" {
		return def
	}
	return v
}

func getEnvAndFatalWithEmpty(envname string) (v string) {
	v = os.Getenv(envname)
	if v == "" {
		log.Fatalf("无法获取变量:%s的内容", envname)
	}
	return
}

func invokeAliSdk(response interface{}, err error) {
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("response is %v", response)
}
