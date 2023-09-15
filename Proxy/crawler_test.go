package proxy

import (
	"fmt"
	"ippool/config"
	"ippool/http"
	"testing"
)

func TestCrawler(t *testing.T) {
	body, _ := http.FetchData()
	bodys := http.MatchIPs(string(body))

	newIppool := config.NewMap(len(bodys))

	for _, v := range bodys {
		//newIppool.Hash.Store(v, nil)
		newIppool.AddIP(v)
	}
	fmt.Println(len(bodys))

}
