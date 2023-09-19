package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func ProcessProtocolHeader(req *http.Request) {
	hopHeaders := []string{
		"Connection",
		"CrawlProxy-Connection",
		"Keep-Alive",
		"CrawlProxy-Authenticate",
		"CrawlProxy-Authorization",
		"Te",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}

	for key, _ := range req.Header {

		if contains(hopHeaders, key) {
			req.Header.Del(key)
		}
	}
}
func contains(list []string, key string) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}
	return false
}
func FetchURL(url string) ([]byte, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "proxyspace.pro")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
func MatchIPs(text string) []string {
	// 将文本按行分割成切片
	lines := strings.Split(text, "\n")

	// 定义 IP 地址的正则表达式
	ipRegex := `\b(?:\d{1,3}\.){3}\d{1,3}:\d+\b`

	// 编译正则表达式
	regex, err := regexp.Compile(ipRegex)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return nil
	}

	// 存储匹配到的 IP 地址
	var matches []string

	// 遍历每一行文本进行匹配
	for _, line := range lines {
		// 查找匹配的 IP 地址
		ip := regex.FindString(line)
		if ip != "" {
			matches = append(matches, ip)
		}
	}

	return matches
}
