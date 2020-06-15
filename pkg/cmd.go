package main

import (
	dns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net/http"
	_ "org.penitence/ddns/pkg/env"
	"os"
)

func main() {

	accessKey := getEnvAndFatalWithEmpty("accessKey")
	accessSecret := getEnvAndFatalWithEmpty("accessSecret")
	baseDomain := getEnvAndFatalWithEmpty("baseDomain")
	domainRR := getEnvAndFatalWithEmpty("domainRR")

	publicIp := findPublicIP()

	client, _ := dns.NewClientWithAccessKey("cn-hangzhou", accessKey, accessSecret)

	recordRequest := dns.CreateDescribeDomainRecordsRequest()
	recordRequest.Scheme = "https"
	recordRequest.DomainName = baseDomain
	recordRequest.TypeKeyWord = "A"
	recordRequest.RRKeyWord = domainRR

	response, err := client.DescribeDomainRecords(recordRequest)

	if err != nil {
		log.Fatalln(err)
	}

	if len(response.DomainRecords.Record) == 0 {
		log.Println("未发现dns解析记录, 添加一个解析")

		recordAddRequest := dns.CreateAddDomainRecordRequest()
		recordAddRequest.Scheme = "https"
		recordAddRequest.DomainName = baseDomain
		recordAddRequest.RR = domainRR
		recordAddRequest.Type = "A"
		recordAddRequest.Value = "192.168.1.1"
		invokeAliSdk(client.AddDomainRecord(recordAddRequest))
	} else {
		log.Println("发现dns解析记录")
		for _, record := range response.DomainRecords.Record {
			log.Printf("记录内容 : %v\n", record)
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
func findPublicIP() (publicIp string) {
	res, err := http.Get("https://ip.cip.cc")
	if err != nil {
		log.Fatalf("获取互联网ip失败 : %v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("读取请求内容失败 : %v", err)
	}
	publicIp = string(body)
	log.Printf("访问https://ifconfig.me/ip获取到的互联网ip为:%s\n", publicIp)
	return
}

func getEnvAndFatalWithEmpty (envname string) (v string) {
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
