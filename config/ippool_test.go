package config

import (
	"fmt"
	"testing"
)

func TestMyFunction(t *testing.T) {

	ippool := NewMap(50, "https://proxyspace.pro/https.txt")
	go ippool.LoopGetIP()
	for {

		//time.Sleep(time.Second * 1)
		fmt.Println(ippool.GetIP(), "开始循环")
	}
}
