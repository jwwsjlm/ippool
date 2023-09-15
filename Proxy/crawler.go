package proxy

import (
	"ippool/config"
	"ippool/http"
)

func Crawler(ipp *config.IpPool) {
	body, _ := http.FetchData()
	bodys := http.MatchIPs(string(body))
	for _, v := range bodys {
		ipp.AddIP(v)
	}
}
func WriteToMap(ip *config.IpPool, key string) string {
	existingValue, loaded := ip.Hash.LoadOrStore(key, ip.GetIP()) // 尝试向 sync.Map 中写入数据

	if loaded {
		return existingValue.(string) // 如果已存在相同的键，则返回已存在的值
	}
	ret, ok := ip.Hash.Load(key)
	if ok {
		return ret.(string)
	}
	return ip.GetIP() // 如果不存在相同的键，则返回特定的一句话
}
