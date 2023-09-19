package config

import (
	"ippool/utils"
	"sync"
)

type IpPool struct {
	Hash   sync.Map
	IpChan chan string
	url    string
	ip     []string
	Count  int
}

func NewMap(i int, u string) *IpPool {

	ip := IpPool{
		Hash:   sync.Map{},
		IpChan: make(chan string, i),
		url:    u,
		Count:  i,
	}

	return &ip

}
func (i IpPool) LoopGetIP() {
	for {
		body, err := utils.FetchURL(i.url)
		if err != nil {
			panic(err)
		}
		ips := utils.MatchIPs(string(body))
		for _, v := range ips {
			//newIppool.Hash.Store(v, nil)
			//fmt.Println(len(ips), len(i.IpChan), i.Count, v)
			if len(i.IpChan) < i.Count {
				i.AddIP(v)
			} else {
				break
			}
		}
	}
}
func (i *IpPool) AddIP(ip string) {
	i.IpChan <- ip
}
func (i *IpPool) GetIP() string {
	p := <-i.IpChan
	return p
}

func (i *IpPool) WriteToMap(key string) string {
	if key == "" {
		return i.GetIP()
	}
	existingValue, loaded := i.Hash.LoadOrStore(key, i.GetIP()) // 尝试向 sync.Map 中写入数据

	if loaded {
		return existingValue.(string) // 如果已存在相同的键，则返回已存在的值
	}
	ret, ok := i.Hash.Load(key)
	if ok {
		return ret.(string)
	}
	return i.GetIP() // 如果不存在相同的键，则返回特定的一句话
}
