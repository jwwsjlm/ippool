package http

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func FetchData() ([]byte, error) {
	url := "https://proxyspace.pro/https.txt"

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
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
