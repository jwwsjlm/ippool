package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"ippool/config"
	"ippool/utils"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var newIppool *config.IpPool

func main() {

	newIppool = config.NewMap(50, "https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt")
	go newIppool.LoopGetIP()
	fmt.Println("ip导入成功")
	server := &http.Server{
		Addr: ":9981",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodConnect {
				handleTunneling(w, r) // 处理 CONNECT 方法的请求
			} else {
				handleHTTP(w, r) // 处理普通的 HTTP 请求
			}
		}),
	}
	log.Fatal(server.ListenAndServe())
}

// 处理普通的 HTTP 请求

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	Autho := req.Header.Get("Proxy-Authorization") // 获取请求头中的代理授权信息
	Authos := strings.Split(Autho, " ")

	utils.ProcessProtocolHeader(req) // 删除协议头当中hop-by-hop协议头
	//fmt.Println(req.Header)

	//if len(Authos) == 2 {
	var decodedAuth []byte
	if len(Authos) == 2 {
		decodedAuth, _ = base64.StdEncoding.DecodeString(Authos[1])
	} else {
		//decodedAuth= interface{}
	}
	//	decodedAuth, _ := base64.StdEncoding.DecodeString(Authos[1]) // 解码代理授权信息
	// 打印解码后的代理授权信息
	p := "https://" + newIppool.WriteToMap(string(decodedAuth))
	proxyurl, _ := url.Parse(p)
	//returnhttp://httpbin.org/get?show_env=1
	//}
	// 打印解码后的代理授权信息
	fmt.Println(p)
	// 创建一个自定义的 Transport
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyurl), // 设置代理地址
		TLSClientConfig: tlsConfig,               //跳过ssl验证
	}

	// 创建一个自定义的 Client
	client := &http.Client{
		Transport: transport,
	}

	//resp, err := client.Do(req) // 发起 HTTP 请求
	fmt.Println(req.URL)
	resp, err := client.Do(&http.Request{
		Method: req.Method,
		URL:    req.URL,
		Header: req.Header,
		Body:   req.Body,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header) // 复制响应头部
	w.WriteHeader(resp.StatusCode)      // 设置响应状态码

	io.Copy(w, resp.Body) // 将响应体复制到客户端
}

// 复制头部字段
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// 处理 CONNECT 方法的请求
func handleTunneling(w http.ResponseWriter, r *http.Request) {
	//Autho := r.Header.Get("CrawlProxy-Authorization") // 获取请求头中的代理授权信息
	headers := r.Header
	Autho := r.Header.Get("Proxy-Authorization") // 获取请求头中的代理授权信息
	// 遍历协议头并打印
	for key, values := range headers {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	//fmt.Println(Autho)
	Authos := strings.Split(Autho, " ")
	//body, err := io.ReadAll(r.Body)
	//fmt.Println(len(Authos), len(Authos) != 0, string(body), len(headers))
	if len(Authos) == 2 {
		decodedAuth, _ := base64.StdEncoding.DecodeString(Authos[1]) // 解码代理授权信息
		// 打印解码后的代理授权信息
		p := newIppool.WriteToMap(string(decodedAuth))

		fmt.Println(string(decodedAuth), p)

		response := "Custom Response: " + p
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		write, err := w.Write([]byte(response))
		if err != nil {
			fmt.Println("错误", write)
			return
		}

		return
	}
	//fmt.Println(len(r.Header))
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second) // 建立与目标服务器的 TCP 连接
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)

	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack() // 获取客户端连接
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(dest_conn, client_conn) // 将客户端数据转发给目标服务器
	go transfer(client_conn, dest_conn) // 将目标服务器的响应转发给客户端
}

// 将数据从源连接复制到目标连接
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()

	io.Copy(destination, source)
}
