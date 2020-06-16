package resolver

import (
	"log"
	"testing"
)

func TestFindFirstIp(t *testing.T) {
	ipContext := "192.168.10.1\n"
	ip := FindFirstIp(ipContext)

	log.Println(ip)
	if ip != "192.168.10.1" {
		t.Fatalf("无法找出正确的正则表达式")
	}

	ipContext = "1oiasjdo122o.1sada.21.asd"

	ip = FindFirstIp(ipContext)

	if ip != "" {
		t.Fatalf("应该返回空字符串")
	}
}
