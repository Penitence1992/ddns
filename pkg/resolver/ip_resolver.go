package resolver

import (
	"log"
	"regexp"
)

var ipMatcher *regexp.Regexp

func init() {

	matcher, err := regexp.Compile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	if err != nil {
		log.Fatalf("初始化ip正则表达式提取器失败 : %v", err)
	}

	ipMatcher = matcher
}

func FindFirstIp(context string) (ip string) {
	ip = ipMatcher.FindString(context)
	return
}
